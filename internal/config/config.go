package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	AppDirName = ".cc-notify"
	FileName   = "config.json"
)

type Config struct {
	Schema   any                      `json:"schema,omitempty"`
	Prefix   string                   `json:"prefix,omitempty"`
	LogLevel string                   `json:"log_level,omitempty"`
	LogDir   string                   `json:"log_dir,omitempty"`
	Channels map[string]ChannelConfig `json:"channels,omitempty"`
	Events   map[string]EventConfig   `json:"events,omitempty"`
}

type ChannelConfig struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Webhook string `json:"webhook,omitempty"`
	Secret  string `json:"secret,omitempty"`
	// DingTalk
	Markdown *bool `json:"markdown,omitempty"`
	// Feishu
	Card *bool `json:"card,omitempty"`
}

type EventConfig struct {
	Channels []string `json:"channels,omitempty"`
}

func DefaultPath() string {
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, AppDirName, FileName)
	}
	return FileName
}

func LegacyPath() string {
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".cc-hooks-notify", FileName)
	}
	return ""
}

func Load(path string) (Config, string, error) {
	candidates := []string{}
	if path != "" {
		candidates = append(candidates, path)
	} else if envPath := strings.TrimSpace(os.Getenv("CC_NOTIFY_CONFIG_PATH")); envPath != "" {
		candidates = append(candidates, envPath)
	} else {
		candidates = append(candidates, DefaultPath())
		if legacy := LegacyPath(); legacy != "" {
			candidates = append(candidates, legacy)
		}
	}

	var cfg Config
	loadedPath := ""
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if _, err := os.Stat(candidate); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return Config{}, "", err
		}
		data, err := os.ReadFile(candidate)
		if err != nil {
			return Config{}, "", err
		}
		if err := json.Unmarshal(data, &cfg); err != nil {
			return Config{}, "", fmt.Errorf("parse config %s: %w", candidate, err)
		}
		loadedPath = candidate
		break
	}

	applyDefaults(&cfg)
	applyEnvOverrides(&cfg)
	return cfg, loadedPath, nil
}

func Save(path string, cfg Config) error {
	if path == "" {
		path = DefaultPath()
	}
	applyDefaults(&cfg)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

func Redacted(cfg Config) Config {
	out := cfg
	out.Channels = map[string]ChannelConfig{}
	for name, ch := range cfg.Channels {
		if ch.Webhook != "" {
			ch.Webhook = redact(ch.Webhook)
		}
		if ch.Secret != "" {
			ch.Secret = redact(ch.Secret)
		}
		out.Channels[name] = ch
	}
	return out
}

func applyDefaults(cfg *Config) {
	if cfg.LogLevel == "" {
		cfg.LogLevel = "INFO"
	}
	if cfg.Channels == nil {
		cfg.Channels = map[string]ChannelConfig{}
	}
	if cfg.Events == nil {
		cfg.Events = map[string]EventConfig{}
	}
}

func applyEnvOverrides(cfg *Config) {
	setChannel := func(name, webhookKey, secretKey string) {
		webhook := firstEnv(webhookKey, "CC_HOOKS_NOTIFY_"+strings.TrimPrefix(webhookKey, "CC_NOTIFY_"))
		secret := firstEnv(secretKey, "CC_HOOKS_NOTIFY_"+strings.TrimPrefix(secretKey, "CC_NOTIFY_"))
		if webhook == "" && secret == "" {
			return
		}
		ch := cfg.Channels[name]
		if webhook != "" {
			ch.Webhook = webhook
			enabled := true
			ch.Enabled = &enabled
		}
		if secret != "" {
			ch.Secret = secret
		}
		cfg.Channels[name] = ch
	}
	setChannel("dingtalk", "CC_NOTIFY_DINGTALK_WEBHOOK", "CC_NOTIFY_DINGTALK_SECRET")
	setChannel("feishu", "CC_NOTIFY_FEISHU_WEBHOOK", "CC_NOTIFY_FEISHU_SECRET")
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func redact(value string) string {
	if len(value) <= 12 {
		return "***"
	}
	return value[:6] + "..." + value[len(value)-4:]
}
