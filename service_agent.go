// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

import "context"

// CreateChatCompletion - Nebula Agent chat completions. POST /chat/completions (non-streaming).
// Docs: https://docs.hiddensystems.ai/nebula/reference/post-chat-completions/
func (s *AgentService) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return postJSON[ChatCompletionResponse](ctx, s.client, "/chat/completions", req.body(false), "create_chat_completion")
}

// StreamChatCompletion - Nebula Agent chat completions. POST /chat/completions with stream=true (SSE).
// Docs: https://docs.hiddensystems.ai/nebula/reference/post-chat-completions/
func (s *AgentService) StreamChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatStream, error) {
	if err := req.validate(); err != nil {
		return nil, err
	}
	return streamSSE(ctx, s.client, "/chat/completions", req.body(true), "create_chat_completion")
}
