// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import (
	"fmt"
)

// ChatCompletionRequest is the OpenAI-compatible Nebula Agent request body
// of `create_chat_completion` (POST /chat/completions).
// Docs: https://docs.hiddensystems.ai/nebula/reference/post-chat-completions/
type ChatCompletionRequest struct {
	// Messages is the complete conversation history; at least one message is
	// required and the Agent is stateless, so resend the full history each call.
	Messages []ChatMessage
	// Model is accepted for OpenAI compatibility and ignored by Nebula.
	Model *string
	// ThinkingMode selects the reasoning mode: auto or deep.
	ThinkingMode *ThinkingMode
	// Stream requests server-sent events. Callers normally leave this nil and
	// pick Agent.CreateChatCompletion or Agent.StreamChatCompletion instead.
	Stream *bool
}

// ChatMessage is one OpenAI-style conversation message.
type ChatMessage struct {
	Role    ChatRole `json:"role"`
	Content string   `json:"content"`
}

func (r *ChatCompletionRequest) validate() error {
	if r == nil {
		return &ValidationError{Op: "create_chat_completion", Field: "request", Msg: "request is required"}
	}
	if len(r.Messages) == 0 {
		return &ValidationError{Op: "create_chat_completion", Field: "messages", Msg: "at least one message is required"}
	}
	for i, m := range r.Messages {
		if m.Role == "" {
			return &ValidationError{Op: "create_chat_completion", Field: "messages", Msg: fmt.Sprintf("message %d: role is required", i)}
		}
		if m.Content == "" {
			return &ValidationError{Op: "create_chat_completion", Field: "messages", Msg: fmt.Sprintf("message %d: content must be non-empty", i)}
		}
	}
	return nil
}

type chatCompletionBody struct {
	Messages     []ChatMessage `json:"messages"`
	Model        *string       `json:"model,omitempty"`
	ThinkingMode *ThinkingMode `json:"thinking_mode,omitempty"`
	Stream       bool          `json:"stream"`
}

func (r *ChatCompletionRequest) body(stream bool) chatCompletionBody {
	return chatCompletionBody{
		Messages:     r.Messages,
		Model:        r.Model,
		ThinkingMode: r.ThinkingMode,
		Stream:       stream,
	}
}
