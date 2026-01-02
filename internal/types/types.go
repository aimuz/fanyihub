// Package types provides shared type definitions for the application.
package types

// Provider represents an LLM provider configuration.
type Provider struct {
	Name            string  `json:"name"`
	Type            string  `json:"type"` // "openai", "openai-compatible", "gemini", "claude"
	BaseURL         string  `json:"base_url,omitempty"`
	APIKey          string  `json:"api_key"`
	Model           string  `json:"model"`
	SystemPrompt    string  `json:"system_prompt,omitempty"`
	MaxTokens       int     `json:"max_tokens,omitempty"`
	Temperature     float64 `json:"temperature,omitempty"`
	Active          bool    `json:"active"`
	DisableThinking bool    `json:"disable_thinking,omitempty"` // For Gemini: set thinkingBudget to 0
}

// DefaultMaxTokens is the default max tokens if not specified.
const DefaultMaxTokens = 1000

// DefaultTemperature is the default temperature if not specified.
const DefaultTemperature = 0.3

// TranslateRequest represents a translation request from the frontend.
type TranslateRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
}

// DetectResult represents the result of language detection.
type DetectResult struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	DefaultTarget string `json:"defaultTarget"`
}

// Usage represents token usage statistics from LLM API calls.
type Usage struct {
	PromptTokens     int  `json:"promptTokens"`
	CompletionTokens int  `json:"completionTokens"`
	TotalTokens      int  `json:"totalTokens"`
	CacheHit         bool `json:"cacheHit"`
}

// TranslateResult represents the result of a translation request.
type TranslateResult struct {
	Text  string `json:"text"`
	Usage Usage  `json:"usage"`
}
