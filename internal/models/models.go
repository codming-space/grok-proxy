package models

import "time"

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIRequest represents the incoming OpenAI-compatible request
type OpenAIRequest struct {
	Model     string    `json:"model"`
	Stream    bool      `json:"stream"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
}

// GrokRequest represents the request sent to Grok API
type GrokRequest struct {
	Message   string `json:"message"`
	ModelName string `json:"modelName"`
}

// Model represents an AI model
type Model struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// ModelList represents a list of AI models
type ModelList struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// StreamResponseChunk represents a chunk in a streaming response
type StreamResponseChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Delta        map[string]string `json:"delta"`
		Index        int               `json:"index"`
		FinishReason interface{}       `json:"finish_reason"`
	} `json:"choices"`
}

// CompletionResponse represents a complete, non-streaming response
type CompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
	} `json:"choices"`
}

// NewModelList creates a new ModelList with predefined models
func NewModelList() ModelList {
	now := time.Now().Unix()
	return ModelList{
		Object: "list",
		Data: []Model{
			{
				ID:      "grok-latest",
				Object:  "model",
				Created: now,
				OwnedBy: "xai",
			},
			{
				ID:      "grok-3",
				Object:  "model",
				Created: now,
				OwnedBy: "xai",
			},
		},
	}
}
