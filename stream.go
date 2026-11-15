package nebula

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
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
//	    for _, choice := range chunk.Choices {
//	        if d := choice.Delta; d != nil && d.Content != nil { fmt.Print(*d.Content) }
//	    }
//	}
//	if err := stream.Err(); err != nil { return err }
type ChatStream struct {
	events chan *ChatCompletionResponse
	done   chan struct{}
	once   sync.Once
	body   io.Closer
	cancel context.CancelFunc
	ctx    context.Context
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
	u, err := t.buildURL(ctx, path, nil)
	if err != nil {
		return nil, err
	}
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, &ValidationError{Op: op, Field: "body", Msg: "cannot serialize request body"}
	}
	ctx, cancel := context.WithTimeout(ctx, t.cfg.timeout)
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
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		err := decodeAPIError(resp, op)
		cancel()
		return nil, err
	}
	if mediaType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type")); err != nil || mediaType != "text/event-stream" {
		resp.Body.Close()
		cancel()
		return nil, &StreamError{Op: op, Err: errors.New("expected Content-Type text/event-stream")}
	}
	s := &ChatStream{
		events: make(chan *ChatCompletionResponse, 16),
		done:   make(chan struct{}),
		body:   resp.Body,
		cancel: cancel,
		ctx:    ctx,
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
		var frame struct {
			Error   json.RawMessage `json:"error"`
			Message string          `json:"message"`
		}
		if json.Unmarshal([]byte(payload), &frame) == nil && len(frame.Error) > 0 && string(frame.Error) != "false" && string(frame.Error) != "null" {
			message := frame.Message
			if message == "" {
				message = "API error event"
			}
			s.setErr(&StreamError{Op: op, Err: errors.New(message)})
			return false
		}
		var chunk ChatCompletionResponse
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil || len(chunk.Choices) == 0 && chunk.Usage == nil {
			s.setErr(&StreamError{Op: op, Err: errors.New("malformed SSE completion data")})
			return false
		}
		chunk.Metadata = s.meta
		select {
		case s.events <- &chunk:
			return true
		case <-s.done:
			return false
		case <-s.ctx.Done():
			s.setErr(&StreamError{Op: op, Err: s.ctx.Err()})
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
			if dataBuf.Len()+len(line) > 4*1024*1024 {
				s.setErr(&StreamError{Op: op, Err: errors.New("SSE event exceeds 4 MiB")})
				return
			}
			dataBuf.WriteString(strings.TrimPrefix(line, "data:"))
			dataBuf.WriteByte('\n')
			continue
		}
		// ignore event:, id:, retry: fields — data carries the payload
	}
	select {
	case <-s.done: // Explicit Close is a normal way to stop early.
		return
	default:
	}
	if err := s.ctx.Err(); err != nil {
		s.setErr(&StreamError{Op: op, Err: err})
	} else if err := scanner.Err(); err != nil {
		s.setErr(&StreamError{Op: op, Err: err})
	} else {
		s.setErr(&StreamError{Op: op, Err: io.ErrUnexpectedEOF})
	}
}
