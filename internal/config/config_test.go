package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoadAndEnvOverride(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	markdown := true
	err := Save(path, Config{
		Channels: map[string]ChannelConfig{
			"dingtalk": {Webhook: "old", Markdown: &markdown},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("CC_NOTIFY_DINGTALK_WEBHOOK", "new")
	cfg, loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded != path {
		t.Fatalf("loaded = %q", loaded)
	}
	if cfg.Channels["dingtalk"].Webhook != "new" {
		t.Fatalf("webhook = %q", cfg.Channels["dingtalk"].Webhook)
	}
}

func TestRedacted(t *testing.T) {
	cfg := Redacted(Config{
		Channels: map[string]ChannelConfig{
			"feishu": {Webhook: "https://example.com/abcdef123456", Secret: "secret-value-123456"},
		},
	})
	if cfg.Channels["feishu"].Webhook == "https://example.com/abcdef123456" {
		t.Fatal("webhook was not redacted")
	}
}

func TestDefaultPathUsesHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	got := DefaultPath()
	want := filepath.Join(home, AppDirName, FileName)
	if got != want {
		t.Fatalf("DefaultPath() = %q, want %q", got, want)
	}
}

func TestLoadMissingConfigIsEmpty(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	cfg, loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded != "" {
		t.Fatalf("loaded = %q", loaded)
	}
	if cfg.LogLevel == "" {
		t.Fatal("defaults not applied")
	}
	_ = os.Getenv("HOME")
}
