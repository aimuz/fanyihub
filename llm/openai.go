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

type openaiRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

type openaiResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (c *Client) completeOpenAI(messages []Message) (string, types.Usage, error) {
	url := defaultBaseURL
	if c.provider.Type == "openai-compatible" && c.provider.BaseURL != "" {
		url = c.provider.BaseURL
	}

	reqBody := openaiRequest{
		Model:       c.provider.Model,
		Messages:    messages,
		MaxTokens:   c.provider.MaxTokens,
		Temperature: c.provider.Temperature,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", types.Usage{}, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", types.Usage{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.provider.APIKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", types.Usage{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", types.Usage{}, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", types.Usage{}, fmt.Errorf("api error: %d - %s", resp.StatusCode, string(body))
	}

	var chatResp openaiResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", types.Usage{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", types.Usage{}, fmt.Errorf("no choices")
	}

	usage := types.Usage{
		PromptTokens:     chatResp.Usage.PromptTokens,
		CompletionTokens: chatResp.Usage.CompletionTokens,
		TotalTokens:      chatResp.Usage.TotalTokens,
	}

	return chatResp.Choices[0].Message.Content, usage, nil
}
