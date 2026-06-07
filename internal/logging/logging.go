package logging

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipfred/cc-hooks-notify/internal/config"
)

func New(cfg config.Config) (*slog.Logger, func(), error) {
	level := slog.LevelInfo
	switch strings.ToUpper(cfg.LogLevel) {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN", "WARNING":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	}

	logDir := cfg.LogDir
	if logDir == "" {
		logDir = filepath.Join(configDirFallback(), "logs")
	}
	if err := os.MkdirAll(logDir, 0o700); err != nil {
		return slog.Default(), func() {}, err
	}
	file, err := os.OpenFile(filepath.Join(logDir, "cc-notify.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return slog.Default(), func() {}, err
	}
	handler := slog.NewTextHandler(io.Writer(file), &slog.HandlerOptions{Level: level})
	return slog.New(handler), func() { _ = file.Close() }, nil
}

func configDirFallback() string {
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, config.AppDirName)
	}
	return "."
}
