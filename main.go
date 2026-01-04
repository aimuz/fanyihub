package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.aimuz.me/transy/cache"
	"go.aimuz.me/transy/clipboard"
	"go.aimuz.me/transy/config"
	"go.aimuz.me/transy/hotkey"
	"go.aimuz.me/transy/internal/types"
	"go.aimuz.me/transy/langdetect"
	"go.aimuz.me/transy/llm"
	"go.aimuz.me/transy/ocr"
	"go.aimuz.me/transy/screenshot"
)

//go:embed all:frontend/dist
var assets embed.FS

// App is the main application struct bound to Wails.
type App struct {
	ctx    context.Context
	cfg    *config.Config
	hotkey *hotkey.HotkeyManager
	cache  *cache.Cache
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

	// Initialize cache
	a.setupCache()

	a.setupHotkey()
}

func (a *App) shutdown(_ context.Context) {
	if a.hotkey != nil {
		a.hotkey.Stop()
	}
	if a.cache != nil {
		if err := a.cache.Close(); err != nil {
			slog.Error("close cache", "error", err)
		}
	}
}

func (a *App) setupCache() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		slog.Error("get config dir for cache", "error", err)
		return
	}

	cachePath := filepath.Join(configDir, "transy", "cache")
	c, err := cache.New(cachePath)
	if err != nil {
		slog.Error("init cache", "error", err)
		return
	}
	a.cache = c
	slog.Info("cache initialized", "path", cachePath)
}

func (a *App) setupHotkey() {
	a.hotkey = hotkey.NewHotkeyManager(
		func() {
			a.ToggleWindowVisibility()
		},
		func() {
			// Run in goroutine to not block the hotkey listener
			go func() {
				if _, err := a.TakeScreenshotAndOCR(); err != nil {
					slog.Error("ocr screenshot", "error", err)
				}
			}()
		},
	)

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

// TakeScreenshotAndOCR captures a screenshot and performs OCR.
// Returns the recognized text.
func (a *App) TakeScreenshotAndOCR() (string, error) {
	// Hide window to allow capturing screen behind it
	runtime.WindowHide(a.ctx)

	// Give a little time for window to hide
	time.Sleep(100 * time.Millisecond)

	imagePath, err := screenshot.CaptureInteractive()
	if err != nil {
		// If cancelled or failed, show window again if not active
		runtime.WindowShow(a.ctx)
		return "", fmt.Errorf("capture screenshot: %w", err)
	}
	defer os.Remove(imagePath) // Clean up temp file

	text, err := ocr.RecognizeText(imagePath)
	if err != nil {
		runtime.WindowShow(a.ctx)
		return "", fmt.Errorf("recognize text: %w", err)
	}

	// Show window and populate text
	runtime.WindowShow(a.ctx)

	if text != "" {
		runtime.EventsEmit(a.ctx, "set-clipboard-text", text)
	}

	return text, nil
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

func (a *App) TranslateWithLLM(req types.TranslateRequest) (types.TranslateResult, error) {
	provider := a.GetActiveProvider()
	if provider == nil {
		return types.TranslateResult{}, fmt.Errorf("no active provider configured")
	}

	cacheKey := a.translationCacheKey(provider, req)

	// Check cache first.
	if result, ok := a.getCachedTranslation(cacheKey); ok {
		return result, nil
	}

	// Call LLM API.
	text, usage, err := a.callLLM(provider, req)
	if err != nil {
		return types.TranslateResult{}, fmt.Errorf("translate %q: %w", truncate(req.Text, 32), err)
	}

	// Store result in cache (best effort).
	a.cacheTranslation(cacheKey, text, usage)

	return types.TranslateResult{Text: text, Usage: usage}, nil
}

// translationCacheKey generates a cache key for the translation request.
func (a *App) translationCacheKey(p *types.Provider, req types.TranslateRequest) string {
	return cache.GenerateKey(p.Name, p.Model, req.SourceLang, req.TargetLang, req.Text)
}

// getCachedTranslation retrieves a cached translation if available.
func (a *App) getCachedTranslation(key string) (types.TranslateResult, bool) {
	if a.cache == nil {
		return types.TranslateResult{}, false
	}

	entry, found := a.cache.Get(key)
	if !found {
		return types.TranslateResult{}, false
	}

	return types.TranslateResult{
		Text: entry.Text,
		Usage: types.Usage{
			PromptTokens:     entry.Usage.PromptTokens,
			CompletionTokens: entry.Usage.CompletionTokens,
			TotalTokens:      entry.Usage.TotalTokens,
			CacheHit:         true,
		},
	}, true
}

// cacheTranslation stores a translation result in the cache.
func (a *App) cacheTranslation(key, text string, usage types.Usage) {
	if a.cache == nil {
		return
	}

	entry := &cache.Entry{
		Text: text,
		Usage: cache.Usage{
			PromptTokens:     usage.PromptTokens,
			CompletionTokens: usage.CompletionTokens,
			TotalTokens:      usage.TotalTokens,
		},
		CreatedAt: time.Now(),
	}

	if err := a.cache.Set(key, entry, cache.DefaultTTL); err != nil {
		slog.Warn("cache translation", "error", err)
	}
}

// callLLM invokes the LLM API to perform translation.
func (a *App) callLLM(p *types.Provider, req types.TranslateRequest) (string, types.Usage, error) {
	client := llm.NewClient(p)

	messages := []llm.Message{
		{Role: "system", Content: p.SystemPrompt},
		{Role: "user", Content: fmt.Sprintf(
			"please translate the following text from %s to %s:\n\n%s",
			req.SourceLang, req.TargetLang, req.Text,
		)},
	}

	return client.Complete(messages)
}

// truncate shortens a string for logging purposes.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// ─────────────────────────────────────────────────────────────────────────────
// Main Entry
// ─────────────────────────────────────────────────────────────────────────────

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:  "Transy",
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
				Title:   "Transy",
				Message: "©2025 Transy. All rights reserved.",
			},
		},
	})
	if err != nil {
		slog.Error("run app", "error", err)
	}
}
