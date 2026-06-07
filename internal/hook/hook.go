package hook

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
)

type Payload map[string]any

func Parse(r io.Reader) (Payload, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return Payload{"parse_error": err.Error()}, nil
	}
	text := strings.TrimSpace(string(raw))
	if text == "" {
		return Payload{}, nil
	}
	var payload Payload
	if err := json.Unmarshal([]byte(text), &payload); err != nil {
		return Payload{"raw_input": text, "parse_error": err.Error()}, nil
	}
	return payload, nil
}

func ProjectName(cwd string) string {
	if strings.TrimSpace(cwd) == "" {
		return "unknown"
	}
	normalized := strings.TrimRight(strings.ReplaceAll(cwd, "\\", "/"), "/")
	parts := strings.Split(normalized, "/")
	if len(parts) >= 2 {
		return filepath.ToSlash(filepath.Join(parts[len(parts)-2], parts[len(parts)-1]))
	}
	return parts[0]
}

func String(payload Payload, key string) string {
	value, ok := payload[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return ""
	}
}

func IsCodexStop(payload Payload) bool {
	return strings.EqualFold(String(payload, "hook_event_name"), "stop") && payload["turn_id"] != nil
}
