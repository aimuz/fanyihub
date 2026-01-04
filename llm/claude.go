package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.aimuz.me/transy/internal/types"
)

// https://api.anthropic.com/v1/messages
const defaultClaudeBaseURL = "https://api.anthropic.com/v1/messages"

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeRequest struct {
	Model     string          `json:"model"`
	Messages  []claudeMessage `json:"messages"`
	System    string          `json:"system,omitempty"`
	MaxTokens int             `json:"max_tokens,omitempty"`
}

type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Usage *struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage,omitempty"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *Client) completeClaude(messages []Message) (string, types.Usage, error) {
	var claudeMsgs []claudeMessage
	var systemPrompt string

	for _, msg := range messages {
		if msg.Role == "system" {
			systemPrompt += msg.Content
			continue
		}
		claudeMsgs = append(claudeMsgs, claudeMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	reqBody := claudeRequest{
		Model:     c.provider.Model,
		Messages:  claudeMsgs,
		System:    systemPrompt,
		MaxTokens: c.provider.MaxTokens,
	}

	if reqBody.MaxTokens == 0 {
		reqBody.MaxTokens = 1024 // Claude requires max_tokens
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", types.Usage{}, fmt.Errorf("marshal request: %w", err)
	}

	baseURL := defaultClaudeBaseURL
	if c.provider.BaseURL != "" {
		baseURL = c.provider.BaseURL
	}

	req, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", types.Usage{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("x-api-key", c.provider.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", types.Usage{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", types.Usage{}, fmt.Errorf("read response: %w", err)
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", types.Usage{}, fmt.Errorf("unmarshal response: %w", err)
	}

	if claudeResp.Error != nil {
		return "", types.Usage{}, fmt.Errorf("api error: %s - %s", claudeResp.Error.Type, claudeResp.Error.Message)
	}

	if len(claudeResp.Content) == 0 {
		return "", types.Usage{}, fmt.Errorf("no content returned")
	}

	var usage types.Usage
	if claudeResp.Usage != nil {
		usage = types.Usage{
			PromptTokens:     claudeResp.Usage.InputTokens,
			CompletionTokens: claudeResp.Usage.OutputTokens,
			TotalTokens:      claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens,
		}
	}

	return claudeResp.Content[0].Text, usage, nil
}
