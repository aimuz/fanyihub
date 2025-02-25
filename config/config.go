package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aimuz/fanyihub/llm"
)

const (
	// AppName is the name of the application
	AppName = "fanyihub"
	// ConfigFileName is the name of the config file
	ConfigFileName = "config.json"
)

// Config represents the application configuration
type Config struct {
	Providers []llm.Provider `json:"providers"`
	// 默认翻译语言对
	DefaultLanguages map[string]string `json:"default_languages"`
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("get config path: %w", err)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果配置文件不存在，返回默认配置
			return &Config{
				Providers: []llm.Provider{},
				DefaultLanguages: map[string]string{
					"zh": "en", // 中文默认翻译为英语
					"en": "zh", // 英语默认翻译为中文
				},
			}, nil
		}
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// 确保默认语言对存在
	if config.DefaultLanguages == nil {
		config.DefaultLanguages = map[string]string{
			"zh": "en", // 中文默认翻译为英语
			"en": "zh", // 英语默认翻译为中文
		}
	}

	return &config, nil
}

// Save saves the configuration to the config file
func (c *Config) Save() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("get config path: %w", err)
	}

	// 确保配置目录存在
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("get user config dir: %w", err)
	}

	return filepath.Join(configDir, AppName, ConfigFileName), nil
}
