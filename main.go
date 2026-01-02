package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"slices"

	"github.com/aimuz/fanyihub/clipboard"
	"github.com/aimuz/fanyihub/config"
	"github.com/aimuz/fanyihub/hotkey"
	"github.com/aimuz/fanyihub/langdetect"
	"github.com/aimuz/fanyihub/llm"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	config        *config.Config
	hotkeyManager *hotkey.HotkeyManager
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

	// 初始化全局快捷键
	a.setupGlobalHotkeys()
}

// setupGlobalHotkeys 设置全局快捷键
func (a *App) setupGlobalHotkeys() {
	// 创建窗口切换回调函数
	toggleCb := func() {
		// 显示或隐藏主窗口
		a.ToggleWindowVisibility()
		slog.Info("触发全局快捷键：切换窗口")
	}

	// 初始化热键管理器
	a.hotkeyManager = hotkey.NewHotkeyManager(toggleCb)
	err := a.hotkeyManager.Start()
	if err != nil {
		slog.Error("启动全局快捷键失败", "error", err.Error())
	} else {
		slog.Info("全局快捷键已启动")
	}
}

// ToggleWindowVisibility 切换窗口的可见性（显示/隐藏）
func (a *App) ToggleWindowVisibility() {
	runtime.WindowShow(a.ctx)
	// 获取剪贴板内容并发送到前端
	clipboardText, err := clipboard.GetText(a.ctx)
	if err != nil {
		slog.Error("获取剪贴板内容失败", "error", err.Error())
	}
	if clipboardText != "" {
		slog.Info("剪贴板内容", "text", clipboardText)
		runtime.EventsEmit(a.ctx, "set-clipboard-text", clipboardText)
	}
	slog.Info("显示窗口并置于前台")
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	// 停止全局快捷键监听
	if a.hotkeyManager != nil {
		a.hotkeyManager.Stop()
		slog.Info("全局快捷键已停止")
	}
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
	// 找到要删除的提供商
	index := slices.IndexFunc(a.config.Providers, func(p llm.Provider) bool {
		return p.Name == name
	})
	if index == -1 {
		return fmt.Errorf("provider not found")
	}

	wasActive := a.config.Providers[index].Active

	// 删除提供商
	a.config.Providers = slices.Delete(a.config.Providers, index, index+1)
	// 如果删除的是激活的提供商，且还有其他提供商，激活第一个
	if wasActive && len(a.config.Providers) > 0 {
		a.config.Providers[0].Active = true
	}

	return a.config.Save()
}

// UpdateProvider updates an existing provider
func (a *App) UpdateProvider(name string, updatedProvider llm.Provider) error {
	// 验证提供商配置
	if updatedProvider.Name == "" {
		return fmt.Errorf("provider name is required")
	}
	if updatedProvider.APIKey == "" {
		return fmt.Errorf("api key is required")
	}
	if updatedProvider.Model == "" {
		return fmt.Errorf("model is required")
	}
	if updatedProvider.Type == "openai-compatible" && updatedProvider.BaseURL == "" {
		return fmt.Errorf("base url is required for openai-compatible provider")
	}

	// 添加默认值
	if updatedProvider.MaxTokens == 0 {
		updatedProvider.MaxTokens = 1000
	}
	if updatedProvider.Temperature == 0 {
		updatedProvider.Temperature = 0.3
	}

	// 找到要更新的提供商
	index := slices.IndexFunc(a.config.Providers, func(p llm.Provider) bool {
		return p.Name == name
	})
	if index == -1 {
		return fmt.Errorf("provider not found")
	}

	// 保存原来的激活状态
	wasActive := a.config.Providers[index].Active

	// 如果新的提供商设置为激活状态，则停用其他提供商
	if updatedProvider.Active && !wasActive {
		for i := range a.config.Providers {
			a.config.Providers[i].Active = false
		}
	} else {
		// 保持原有的激活状态
		updatedProvider.Active = wasActive
	}

	// 更新提供商
	a.config.Providers[index] = updatedProvider

	return a.config.Save()
}

type TranslateRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
}

// GetDefaultLanguages 获取默认翻译语言对
func (a *App) GetDefaultLanguages() map[string]string {
	return a.config.DefaultLanguages
}

// SetDefaultLanguage 设置默认翻译语言对
func (a *App) SetDefaultLanguage(sourceLang, targetLang string) error {
	if a.config.DefaultLanguages == nil {
		a.config.DefaultLanguages = make(map[string]string)
	}
	a.config.DefaultLanguages[sourceLang] = targetLang
	return a.config.Save()
}

type DetectLanguageResponse struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	DefaultTarget string `json:"defaultTarget"`
}

// DetectLanguage 检测文本的语言
func (a *App) DetectLanguage(text string) DetectLanguageResponse {
	langCode, langName := langdetect.DetectLanguage(text)

	// 如果检测到了语言，并且有默认的目标语言，一并返回
	targetLang := "en" // 默认为英语
	if langCode != "auto" && a.config.DefaultLanguages != nil {
		if defaultTarget, exists := a.config.DefaultLanguages[langCode]; exists {
			targetLang = defaultTarget
		}
	}

	return DetectLanguageResponse{
		Code:          langCode,
		Name:          langName,
		DefaultTarget: targetLang,
	}
}

// TranslateWithLLM translates text using the specified provider
func (a *App) TranslateWithLLM(req TranslateRequest) (string, error) {
	provider := a.GetActiveProvider()
	if provider == nil {
		return "", fmt.Errorf("no active provider")
	}

	// 直接使用前端传入的确定语言，不再进行自动检测
	sourceLang := req.SourceLang
	targetLang := req.TargetLang

	slog.Info("TranslateWithLLM", "sourceLang", sourceLang, "targetLang", targetLang, "text", req.Text)

	// 创建客户端
	client := llm.NewClient(provider)

	// 准备消息
	messages := []llm.ChatMessage{
		{Role: "system", Content: provider.SystemPrompt},
		{Role: "user", Content: fmt.Sprintf("please translate the following text from %s to %s:\n\n%s", sourceLang, targetLang, req.Text)},
	}

	// 发送请求
	return client.ChatCompletion(messages)
}

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "FanyiHub",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
		},
		OnShutdown: func(ctx context.Context) {
			app.shutdown(ctx)
		},
		Bind: []any{
			app,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHidden(),
			About: &mac.AboutInfo{
				Title:   "FanyiHub",
				Message: "©2025 FanyiHub. All rights reserved.",
			},
		},
	})

	if err != nil {
		slog.Error("run app", "error", err.Error())
	}
}
