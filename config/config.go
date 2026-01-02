// Package config handles application configuration.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/aimuz/fanyihub/internal/types"
)

const (
	appName        = "fanyihub"
	configFileName = "config.json"
)

// Config represents the application configuration.
type Config struct {
	Providers        []types.Provider  `json:"providers"`
	DefaultLanguages map[string]string `json:"default_languages"`
}

// Load loads configuration from the config file.
// Returns default config if file doesn't exist.
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, fmt.Errorf("get config path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// Ensure default languages exist
	if cfg.DefaultLanguages == nil {
		cfg.DefaultLanguages = defaultLanguages()
	}

	return &cfg, nil
}

// Save persists the configuration to disk.
func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return fmt.Errorf("get config path: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	return nil
}

// AddProvider adds a new provider.
func (c *Config) AddProvider(p types.Provider) error {
	if err := validateProvider(p); err != nil {
		return err
	}
	applyDefaults(&p)

	// First provider or explicitly active: deactivate others
	if len(c.Providers) == 0 || p.Active {
		for i := range c.Providers {
			c.Providers[i].Active = false
		}
		p.Active = true
	}

	c.Providers = append(c.Providers, p)
	return c.Save()
}

// UpdateProvider updates an existing provider.
func (c *Config) UpdateProvider(name string, p types.Provider) error {
	if err := validateProvider(p); err != nil {
		return err
	}
	applyDefaults(&p)

	idx := slices.IndexFunc(c.Providers, func(x types.Provider) bool {
		return x.Name == name
	})
	if idx == -1 {
		return fmt.Errorf("provider not found: %s", name)
	}

	wasActive := c.Providers[idx].Active
	if p.Active && !wasActive {
		for i := range c.Providers {
			c.Providers[i].Active = false
		}
	} else {
		p.Active = wasActive
	}

	c.Providers[idx] = p
	return c.Save()
}

// RemoveProvider removes a provider.
func (c *Config) RemoveProvider(name string) error {
	idx := slices.IndexFunc(c.Providers, func(p types.Provider) bool {
		return p.Name == name
	})
	if idx == -1 {
		return fmt.Errorf("provider not found: %s", name)
	}

	wasActive := c.Providers[idx].Active
	c.Providers = slices.Delete(c.Providers, idx, idx+1)

	if wasActive && len(c.Providers) > 0 {
		c.Providers[0].Active = true
	}

	return c.Save()
}

// SetProviderActive checks if provider exists and sets it active.
func (c *Config) SetProviderActive(name string) error {
	found := false
	for i := range c.Providers {
		if c.Providers[i].Name == name {
			c.Providers[i].Active = true
			found = true
		} else {
			c.Providers[i].Active = false
		}
	}
	if !found {
		return fmt.Errorf("provider not found: %s", name)
	}
	return c.Save()
}

// GetActiveProvider returns the currently active provider.
func (c *Config) GetActiveProvider() *types.Provider {
	for i := range c.Providers {
		if c.Providers[i].Active {
			p := c.Providers[i]
			return &p
		}
	}
	// Auto-activate first if none active
	if len(c.Providers) > 0 {
		c.Providers[0].Active = true
		_ = c.Save()
		p := c.Providers[0]
		return &p
	}
	return nil
}

// Helper functions

func validateProvider(p types.Provider) error {
	if p.Name == "" {
		return fmt.Errorf("provider name required")
	}
	if p.APIKey == "" {
		return fmt.Errorf("api key required")
	}
	if p.Model == "" {
		return fmt.Errorf("model required")
	}
	if p.Type == "openai-compatible" && p.BaseURL == "" {
		return fmt.Errorf("base url required for openai-compatible")
	}
	return nil
}

func applyDefaults(p *types.Provider) {
	if p.MaxTokens == 0 {
		p.MaxTokens = types.DefaultMaxTokens
	}
	if p.Temperature == 0 {
		p.Temperature = types.DefaultTemperature
	}
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}
	return filepath.Join(dir, appName, configFileName), nil
}

func defaultConfig() *Config {
	return &Config{
		Providers:        []types.Provider{},
		DefaultLanguages: defaultLanguages(),
	}
}

func defaultLanguages() map[string]string {
	return map[string]string{
		"zh": "en",
		"en": "zh",
	}
}
