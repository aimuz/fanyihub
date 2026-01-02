package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/aimuz/fanyihub/clipboard"
	"github.com/aimuz/fanyihub/config"
	"github.com/aimuz/fanyihub/hotkey"
	"github.com/aimuz/fanyihub/internal/types"
	"github.com/aimuz/fanyihub/langdetect"
	"github.com/aimuz/fanyihub/llm"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

// App is the main application struct bound to Wails.
type App struct {
	ctx    context.Context
	cfg    *config.Config
	hotkey *hotkey.HotkeyManager
}

func NewApp() *App {
	return &App{}
}

// ─────────────────────────────────────────────────────────────────────────────
// Lifecycle
// ─────────────────────────────────────────────────────────────────────────────

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		cfg = &config.Config{}
	}
	a.cfg = cfg

	a.setupHotkey()
}

func (a *App) shutdown(_ context.Context) {
	if a.hotkey != nil {
		a.hotkey.Stop()
	}
}

func (a *App) setupHotkey() {
	a.hotkey = hotkey.NewHotkeyManager(func() {
		a.ToggleWindowVisibility()
	})

	a.hotkey.SetStatusCallback(func(granted bool) {
		runtime.EventsEmit(a.ctx, "accessibility-permission", granted)
		if granted {
			slog.Info("accessibility permission granted")
		} else {
			slog.Warn("accessibility permission denied")
		}
	})

	if err := a.hotkey.Start(); err != nil {
		slog.Error("start hotkey", "error", err)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Window & Clipboard
// ─────────────────────────────────────────────────────────────────────────────

func (a *App) ToggleWindowVisibility() {
	runtime.WindowShow(a.ctx)

	text, err := clipboard.GetText(a.ctx)
	if err != nil {
		slog.Error("get clipboard", "error", err)
		return
	}
	if text != "" {
		runtime.EventsEmit(a.ctx, "set-clipboard-text", text)
	}
}

func (a *App) GetAccessibilityPermission() bool {
	return hotkey.IsAccessibilityEnabled(false)
}

// ─────────────────────────────────────────────────────────────────────────────
// Provider Management (Delegated to Config)
// ─────────────────────────────────────────────────────────────────────────────

func (a *App) GetProviders() []types.Provider {
	return a.cfg.Providers
}

func (a *App) AddProvider(p types.Provider) error {
	return a.cfg.AddProvider(p)
}

func (a *App) UpdateProvider(name string, p types.Provider) error {
	return a.cfg.UpdateProvider(name, p)
}

func (a *App) RemoveProvider(name string) error {
	return a.cfg.RemoveProvider(name)
}

func (a *App) SetProviderActive(name string) error {
	return a.cfg.SetProviderActive(name)
}

func (a *App) GetActiveProvider() *types.Provider {
	return a.cfg.GetActiveProvider()
}

// ─────────────────────────────────────────────────────────────────────────────
// Language Settings
// ─────────────────────────────────────────────────────────────────────────────

func (a *App) GetDefaultLanguages() map[string]string {
	return a.cfg.DefaultLanguages
}

func (a *App) SetDefaultLanguage(src, dst string) error {
	if a.cfg.DefaultLanguages == nil {
		a.cfg.DefaultLanguages = make(map[string]string)
	}
	a.cfg.DefaultLanguages[src] = dst
	return a.cfg.Save()
}

func (a *App) DetectLanguage(text string) types.DetectResult {
	code, name := langdetect.Detect(text)

	target := "en"
	if code != "auto" && a.cfg.DefaultLanguages != nil {
		if t, ok := a.cfg.DefaultLanguages[code]; ok {
			target = t
		}
	}

	return types.DetectResult{
		Code:          code,
		Name:          name,
		DefaultTarget: target,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Translation
// ─────────────────────────────────────────────────────────────────────────────

func (a *App) TranslateWithLLM(req types.TranslateRequest) (string, error) {
	provider := a.GetActiveProvider()
	if provider == nil {
		return "", fmt.Errorf("no active provider")
	}

	client := llm.NewClient(provider)

	messages := []llm.Message{
		{Role: "system", Content: provider.SystemPrompt},
		{Role: "user", Content: fmt.Sprintf(
			"please translate the following text from %s to %s:\n\n%s",
			req.SourceLang, req.TargetLang, req.Text,
		)},
	}

	return client.Complete(messages)
}

// ─────────────────────────────────────────────────────────────────────────────
// Main Entry
// ─────────────────────────────────────────────────────────────────────────────

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "FanyiHub",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind:       []any{app},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHidden(),
			About: &mac.AboutInfo{
				Title:   "FanyiHub",
				Message: "©2025 FanyiHub. All rights reserved.",
			},
		},
	})
	if err != nil {
		slog.Error("run app", "error", err)
	}
}
