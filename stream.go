package nebula

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

// ChatStream is a closeable server-sent-events stream of chat completion
// chunks. Read events from Events() until the channel closes, then check
// Err(). Always call Close (or cancel the context) to release the body.
//
//	stream, err := client.Agent.StreamChatCompletion(ctx, req)
//	if err != nil { return err }
//	defer stream.Close()
//	for chunk := range stream.Events() {
//	    if d := chunk.Choices[0].Delta; d != nil { fmt.Print(d.Content) }
//	}
//	if err := stream.Err(); err != nil { return err }
type ChatStream struct {
	events chan *ChatCompletionResponse
	done   chan struct{}
	once   sync.Once
	body   io.Closer
	cancel context.CancelFunc
	err    error
	meta   ResponseMetadata
	mu     sync.Mutex
}

// Events yields decoded chat completion chunks; the channel closes on [DONE],
// on stream error, or after Close/context cancellation.
func (s *ChatStream) Events() <-chan *ChatCompletionResponse { return s.events }

// Metadata returns the HTTP metadata of the stream's response headers.
func (s *ChatStream) Metadata() ResponseMetadata { return s.meta }

// Err returns the terminal error, if any, after Events() closed.
func (s *ChatStream) Err() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}

func (s *ChatStream) setErr(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.err == nil {
		s.err = err
	}
}

// Close releases the stream body and stops the reader goroutine.
func (s *ChatStream) Close() error {
	s.once.Do(func() {
		if s.cancel != nil {
			s.cancel()
		}
		close(s.done)
		if s.body != nil {
			s.body.Close()
		}
	})
	return nil
}

// streamSSE issues the POST with stream=true and spawns the SSE reader.
func streamSSE(ctx context.Context, c *Client, path string, body any, op string) (*ChatStream, error) {
	if c.cfg.x402Mode {
		return nil, fmt.Errorf("nebula: %s: %w", op, ErrNotX402Eligible)
	}
	t := c.tr
	u, err := t.buildURL(path, nil)
	if err != nil {
		return nil, err
	}
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, &ValidationError{Op: op, Field: "body", Msg: "cannot serialize request body"}
	}
	ctx, cancel := context.WithCancel(ctx)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(buf))
	if err != nil {
		cancel()
		return nil, err
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", t.cfg.userAgent)
	t.applyAuth(req)
	locale, localeMode := localeFor(ctx, t.cfg)
	if locale != "" && localeMode == LocaleModeHeader {
		req.Header.Set("X-Nebula-Locale", string(locale))
	}
	if t.cfg.logger != nil {
		t.cfg.logger.Printf("nebula: POST %s (stream)", redactURL(u))
	}
	resp, err := t.cfg.httpClient.Do(req)
	if err != nil {
		cancel()
		return nil, &TransportError{Op: op, Err: err}
	}
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		cancel()
		return nil, decodeAPIError(resp, op)
	}
	s := &ChatStream{
		events: make(chan *ChatCompletionResponse, 16),
		done:   make(chan struct{}),
		body:   resp.Body,
		cancel: cancel,
		meta:   parseMetadata(resp.StatusCode, resp.Header),
	}
	go s.read(resp.Body, op)
	return s, nil
}

// read parses the SSE stream: event/data pairs, [DONE], API error frames.
func (s *ChatStream) read(body io.Reader, op string) {
	defer close(s.events)
	defer s.Close()
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	var dataBuf bytes.Buffer
	flush := func() bool {
		if dataBuf.Len() == 0 {
			return true
		}
		payload := strings.TrimSpace(dataBuf.String())
		dataBuf.Reset()
		if payload == "[DONE]" {
			return false
		}
		var chunk ChatCompletionResponse
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			// an API error frame inside the stream
			var apiErr struct {
				Error   *bool   `json:"error"`
				Message *string `json:"message"`
			}
			if json.Unmarshal([]byte(payload), &apiErr) == nil && apiErr.Message != nil {
				s.setErr(&StreamError{Op: op, Err: errors.New(*apiErr.Message)})
				return false
			}
			s.setErr(&StreamError{Op: op, Err: fmt.Errorf("malformed SSE data: %.200s", payload)})
			return false
		}
		chunk.Metadata = s.meta
		select {
		case s.events <- &chunk:
			return true
		case <-s.done:
			return false
		}
	}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if !flush() {
				return
			}
			continue
		}
		if strings.HasPrefix(line, ":") {
			continue // comment/heartbeat
		}
		if strings.HasPrefix(line, "data:") {
			dataBuf.WriteString(strings.TrimPrefix(line, "data:"))
			dataBuf.WriteByte('\n')
			continue
		}
		// ignore event:, id:, retry: fields — data carries the payload
	}
	if err := scanner.Err(); err != nil {
		if !errors.Is(err, context.Canceled) && !errors.Is(err, io.EOF) {
			select {
			case <-s.done:
			default:
				s.setErr(&StreamError{Op: op, Err: err})
			}
		}
		return
	}
	// EOF without [DONE]: flush any buffered data then end.
	dataBuf.Reset()
}
