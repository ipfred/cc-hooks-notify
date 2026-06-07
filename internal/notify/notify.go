package notify

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/ipfred/cc-hooks-notify/internal/channels"
	"github.com/ipfred/cc-hooks-notify/internal/config"
	"github.com/ipfred/cc-hooks-notify/internal/hook"
)

type Data struct {
	Title string
	Body  string
	Event string
}

func Build(eventType string, payload hook.Payload, prefix string) Data {
	if prefix == "" {
		prefix = "Claude"
	}
	hookEvent := hook.String(payload, "hook_event_name")
	project := hook.ProjectName(hook.String(payload, "cwd"))
	notificationType := hook.String(payload, "notification_type")
	message := hook.String(payload, "message")
	lastMessage := hook.String(payload, "last_assistant_message")

	bodySource := message
	if bodySource == "" {
		bodySource = lastMessage
	}
	if bodySource == "" {
		bodySource = "none msg"
	}

	evt := "notification"
	lowerType := strings.ToLower(notificationType)
	switch {
	case strings.Contains(lowerType, "permission") || lowerType == "permission_prompt":
		evt = "permission"
	case strings.Contains(lowerType, "idle") || lowerType == "idle_prompt":
		evt = "idle"
	case strings.EqualFold(hookEvent, "stop"):
		evt = "stop"
	}

	titles := map[string]string{
		"permission":   fmt.Sprintf("%s 需要权限确认 | %s", prefix, project),
		"idle":         fmt.Sprintf("%s 正在等待 | %s", prefix, project),
		"stop":         fmt.Sprintf("%s stop | %s", prefix, project),
		"notification": fmt.Sprintf("%s 通知 | %s", prefix, project),
	}

	return Data{
		Title: titles[evt],
		Body:  preview(bodySource, 100),
		Event: evt,
	}
}

func Send(eventType string, payload hook.Payload, cfg config.Config, prefix string, logger *slog.Logger) error {
	if prefix == "" {
		prefix = cfg.Prefix
	}
	if prefix == "" {
		prefix = "Claude"
	}
	data := Build(eventType, payload, prefix)
	targets := enabledChannels(cfg)
	if eventCfg, ok := cfg.Events[eventType]; ok && len(eventCfg.Channels) > 0 {
		targets = eventCfg.Channels
	}
	if len(targets) == 0 {
		logger.Warn("no notification channels enabled")
		return nil
	}
	for _, name := range targets {
		chCfg, ok := cfg.Channels[name]
		if !ok {
			logger.Warn("unknown channel", "channel", name)
			continue
		}
		sender, err := channels.New(name, chCfg)
		if err != nil {
			logger.Warn("unknown channel", "channel", name, "error", err)
			continue
		}
		if err := sender.Send(data.Title, data.Body, data.Event); err != nil {
			logger.Error("channel send failed", "channel", name, "error", err)
			continue
		}
		logger.Info("channel send success", "channel", name)
	}
	return nil
}

func enabledChannels(cfg config.Config) []string {
	names := []string{}
	for name, ch := range cfg.Channels {
		if (ch.Enabled != nil && *ch.Enabled) || strings.TrimSpace(ch.Webhook) != "" {
			names = append(names, name)
		}
	}
	return names
}

func preview(value string, limit int) string {
	if value == "" {
		return ""
	}
	runes := []rune(value)
	if len(runes) <= limit {
		return value
	}
	return string(runes[:limit-3]) + "..."
}
