package main

import (
	"context"
	"embed"
	"fmt"

	"github.com/aimuz/fanyihub/pkg/config"
	"github.com/aimuz/fanyihub/pkg/llm"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

// App struct
type App struct {
	ctx    context.Context
	config *config.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		cfg = &config.Config{} // 使用空配置
	}
	a.config = cfg
}

// GetProviders returns the list of providers
func (a *App) GetProviders() []llm.Provider {
	return a.config.Providers
}

// AddProvider adds a new provider
func (a *App) AddProvider(provider llm.Provider) error {
	// 验证提供商配置
	if provider.Name == "" {
		return fmt.Errorf("provider name is required")
	}
	if provider.APIKey == "" {
		return fmt.Errorf("api key is required")
	}
	if provider.Model == "" {
		return fmt.Errorf("model is required")
	}
	if provider.Type == "openai-compatible" && provider.BaseURL == "" {
		return fmt.Errorf("base url is required for openai-compatible provider")
	}

	// 添加默认值
	if provider.MaxTokens == 0 {
		provider.MaxTokens = 1000
	}
	if provider.Temperature == 0 {
		provider.Temperature = 0.3
	}

	// 如果是第一个提供商或者设置为激活状态，停用其他提供商
	if len(a.config.Providers) == 0 || provider.Active {
		for i := range a.config.Providers {
			a.config.Providers[i].Active = false
		}
		provider.Active = true
	}

	a.config.Providers = append(a.config.Providers, provider)
	return a.config.Save()
}

// SetProviderActive sets a provider as active
func (a *App) SetProviderActive(name string) error {
	var found bool
	for i := range a.config.Providers {
		if a.config.Providers[i].Name == name {
			a.config.Providers[i].Active = true
			found = true
		} else {
			a.config.Providers[i].Active = false
		}
	}
	if !found {
		return fmt.Errorf("provider not found")
	}
	return a.config.Save()
}

// GetActiveProvider returns the active provider
func (a *App) GetActiveProvider() *llm.Provider {
	for i := range a.config.Providers {
		if a.config.Providers[i].Active {
			provider := a.config.Providers[i] // 创建一个副本
			return &provider
		}
	}
	// 如果没有激活的提供商，但有提供商，激活第一个
	if len(a.config.Providers) > 0 {
		a.config.Providers[0].Active = true
		if err := a.config.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
		}
		provider := a.config.Providers[0] // 创建一个副本
		return &provider
	}
	return nil
}

// RemoveProvider removes a provider by name
func (a *App) RemoveProvider(name string) error {
	var wasActive bool
	var index = -1

	// 找到要删除的提供商
	for i, p := range a.config.Providers {
		if p.Name == name {
			wasActive = p.Active
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("provider not found")
	}

	// 删除提供商
	a.config.Providers = append(a.config.Providers[:index], a.config.Providers[index+1:]...)

	// 如果删除的是激活的提供商，且还有其他提供商，激活第一个
	if wasActive && len(a.config.Providers) > 0 {
		a.config.Providers[0].Active = true
	}

	return a.config.Save()
}

type TranslateRequest struct {
	Text      string `json:"text"`
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
}

// TranslateWithLLM translates text using the specified provider
func (a *App) TranslateWithLLM(req TranslateRequest) (string, error) {
	provider := a.GetActiveProvider()
	if provider == nil {
		return "", fmt.Errorf("no active provider")
	}

	// 创建客户端
	client := llm.NewClient(provider)

	// 准备消息
	messages := []llm.ChatMessage{
		{Role: "system", Content: provider.SystemPrompt},
		{Role: "user", Content: fmt.Sprintf("请从%s 到 %s 翻译以下文本：\n%s", req.SourceLang, req.TargetLang, req.Text)},
	}

	// 发送请求
	return client.ChatCompletion(messages)
}

//go:embed assets/*
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "FanyiHub",
		Width:     1024,
		Height:    768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
		},
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   "FanyiHub",
				Message: " 2024 FanyiHub. All rights reserved.",
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
