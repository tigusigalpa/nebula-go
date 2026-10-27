package nebula

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func streamRequest() *ChatCompletionRequest {
	return &ChatCompletionRequest{
		Messages: []ChatMessage{{Role: ChatRoleUser, Content: "hello"}},
	}
}

func TestStreamCompletion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "text/event-stream" {
			t.Errorf("Accept=%q", r.Header.Get("Accept"))
		}
		var body struct {
			Stream bool `json:"stream"`
		}
		readJSON(r, &body)
		if !body.Stream {
			t.Error("stream flag must be true in body")
		}
		data, _ := os.ReadFile(filepath.Join("testdata", "chat_stream.sse"))
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write(data)
	}))
	defer srv.Close()

	c, err := New("k", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	stream, err := c.Agent.StreamChatCompletion(context.Background(), streamRequest())
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	var chunks []*ChatCompletionResponse
	for ch := range stream.Events() {
		chunks = append(chunks, ch)
	}
	if err := stream.Err(); err != nil {
		t.Fatalf("stream err: %v", err)
	}
	if len(chunks) != 3 {
		t.Fatalf("chunks=%d want 3", len(chunks))
	}
	if chunks[0].Choices[0].Delta == nil || *chunks[0].Choices[0].Delta.Role != "assistant" {
		t.Fatalf("chunk0 delta: %+v", chunks[0].Choices[0].Delta)
	}
	if chunks[1].Choices[0].Delta == nil || *chunks[1].Choices[0].Delta.Content != "NVIDIA sentiment strengthened" {
		t.Fatalf("chunk1 delta: %+v", chunks[1].Choices[0].Delta)
	}
	if chunks[2].Choices[0].FinishReason == nil || *chunks[2].Choices[0].FinishReason != "stop" {
		t.Fatalf("finish: %+v", chunks[2].Choices[0].FinishReason)
	}
}

func TestStreamHTTPErrorBeforeStream(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		_, _ = w.Write([]byte(`{"error":true,"message":"bad key"}`))
	}))
	defer srv.Close()
	c, _ := New("k", WithBaseURL(srv.URL))
	_, err := c.Agent.StreamChatCompletion(context.Background(), streamRequest())
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("want unauthorized, got %v", err)
	}
}

func TestStreamMalformedEvent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = w.Write([]byte("data: {this is not json}\n\ndata: [DONE]\n\n"))
	}))
	defer srv.Close()
	c, _ := New("k", WithBaseURL(srv.URL))
	stream, err := c.Agent.StreamChatCompletion(context.Background(), streamRequest())
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()
	for range stream.Events() {
	}
	if stream.Err() == nil {
		t.Fatal("expected malformed event error")
	}
	if !strings.Contains(stream.Err().Error(), "malformed") {
		t.Fatalf("err=%v", stream.Err())
	}
}

func TestStreamCancellation(t *testing.T) {
	started := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(started)
		w.Header().Set("Content-Type", "text/event-stream")
		fl, _ := w.(http.Flusher)
		_, _ = w.Write([]byte("data: {\"id\":\"c1\",\"object\":\"chat.completion.chunk\",\"created\":1,\"model\":\"nebula-agent\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"hi\"},\"finish_reason\":null}]}\n\n"))
		if fl != nil {
			fl.Flush()
		}
		<-r.Context().Done() // hold stream open until client cancels
	}))
	defer srv.Close()
	c, _ := New("k", WithBaseURL(srv.URL))
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := c.Agent.StreamChatCompletion(ctx, streamRequest())
	if err != nil {
		t.Fatal(err)
	}
	<-started
	first, ok := <-stream.Events()
	if !ok || first == nil {
		t.Fatal("expected first chunk")
	}
	cancel()
	stream.Close()
	// drain; channel must close promptly
	timeout := time.After(3 * time.Second)
	for {
		select {
		case _, ok := <-stream.Events():
			if !ok {
				return
			}
		case <-timeout:
			t.Fatal("stream did not close after cancellation")
		}
	}
}

func TestNonStreamingCompletion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Stream   bool `json:"stream"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
			Model        *string `json:"model"`
			ThinkingMode *string `json:"thinking_mode"`
		}
		readJSON(r, &body)
		if body.Stream {
			t.Error("non-streaming call must send stream=false")
		}
		if len(body.Messages) != 1 || body.Messages[0].Role != "user" {
			t.Errorf("messages: %+v", body.Messages)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"chatcmpl_1","object":"chat.completion","created":1786114800,"model":"nebula-agent","choices":[{"index":0,"message":{"role":"assistant","content":"done"},"finish_reason":"stop"}],"usage":{"credits_used":3}}`))
	}))
	defer srv.Close()
	c, _ := New("k", WithBaseURL(srv.URL))
	res, err := c.Agent.CreateChatCompletion(context.Background(), streamRequest())
	if err != nil {
		t.Fatal(err)
	}
	if res.Usage == nil || res.Usage.CreditsUsed != 3 {
		t.Fatalf("usage=%+v", res.Usage)
	}
	if res.Choices[0].Message == nil || res.Choices[0].Message.Content != "done" {
		t.Fatalf("choice=%+v", res.Choices[0])
	}
}

func readJSON(r *http.Request, v any) {
	defer r.Body.Close()
	_ = json.NewDecoder(r.Body).Decode(v)
}
