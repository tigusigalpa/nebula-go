// Code generated from resources/openapi/nebula-openapi.json (OpenAPI 3.0.3, info.version 1.0.2).
// Do not edit by hand: regenerate with tools/nebula-gen.

package nebula

// ChatCompletionResponse models an object of the `agent` service.
type ChatCompletionResponse struct {
	BaseResponse
	// Generated completion choices.
	Choices []ChatCompletionResponseChoice `json:"choices"`
	// Unix timestamp when the completion was created.
	Created int64 `json:"created"`
	// Unique identifier assigned to the chat completion.
	ID string `json:"id"`
	// Nebula agent model identifier.
	Model string `json:"model"`
	// OpenAI-compatible object type.
	Object string `json:"object"`
	// Credits consumed by this completion.
	Usage *ChatCompletionResponseUsage `json:"usage"`
}

// ChatCompletionResponseChoice models an object of the `agent` service.
type ChatCompletionResponseChoice struct {
	// Streaming message fragment.
	Delta *ChatCompletionResponseChoiceDelta `json:"delta"`
	// Reason generation stopped when available.
	FinishReason *string `json:"finish_reason"`
	// Zero-based index of the generated completion choice.
	Index int `json:"index"`
	// Completed assistant message for a non-streaming response.
	Message *ChatCompletionResponseChoiceMessage `json:"message"`
}

// ChatCompletionResponseChoiceDelta models an object of the `agent` service.
type ChatCompletionResponseChoiceDelta struct {
	// Partial assistant response text in this stream fragment.
	Content *string `json:"content"`
	// Conversation role when emitted in this stream fragment.
	Role *string `json:"role"`
}

// ChatCompletionResponseChoiceMessage models an object of the `agent` service.
type ChatCompletionResponseChoiceMessage struct {
	// Completed assistant response text.
	Content string `json:"content"`
	// Conversation role of the completed message.
	Role string `json:"role"`
}

// ChatCompletionResponseUsage models an object of the `agent` service.
type ChatCompletionResponseUsage struct {
	// Nebula credits consumed by the request.
	CreditsUsed int64 `json:"credits_used"`
}
