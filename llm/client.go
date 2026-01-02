// Package llm provides HTTP client for LLM API calls.
package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aimuz/fanyihub/internal/types"
)

const defaultBaseURL = "https://api.openai.com/v1/chat/completions"

// Message represents a chat message.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatRequest represents an OpenAI-compatible chat request.
type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// chatResponse represents an OpenAI-compatible chat response.
type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
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

// Complete sends a chat completion request and returns the response text.
func (c *Client) Complete(messages []Message) (string, error) {
	url := defaultBaseURL
	if c.provider.Type == "openai-compatible" && c.provider.BaseURL != "" {
		url = c.provider.BaseURL
	}

	req := chatRequest{
		Model:       c.provider.Model,
		Messages:    messages,
		MaxTokens:   c.provider.MaxTokens,
		Temperature: c.provider.Temperature,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.provider.APIKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("api error: status=%d body=%s", resp.StatusCode, respBody)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return chatResp.Choices[0].Message.Content, nil
}
