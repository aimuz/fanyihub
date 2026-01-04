// Package llm provides HTTP client for LLM API calls.
package llm

import (
	"net/http"

	"go.aimuz.me/transy/internal/types"
)

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Client is an HTTP client for LLM APIs.
type Client struct {
	provider *types.Provider
	http     *http.Client
}

// NewClient creates a new LLM client for the given provider.
func NewClient(p *types.Provider) *Client {
	return &Client{
		provider: p,
		http:     &http.Client{},
	}
}

// Complete sends a chat completion request and returns the response text and usage.
func (c *Client) Complete(messages []Message) (string, types.Usage, error) {
	switch c.provider.Type {
	case "gemini":
		return c.completeGemini(messages)
	case "claude":
		return c.completeClaude(messages)
	case "openai", "openai-compatible":
		return c.completeOpenAI(messages)
	default:
		// Default to OpenAI format for compatibility
		return c.completeOpenAI(messages)
	}
}
