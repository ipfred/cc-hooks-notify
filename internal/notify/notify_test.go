package notify

import (
	"testing"

	"github.com/ipfred/cc-hooks-notify/internal/hook"
)

func TestBuildStopNotification(t *testing.T) {
	data := Build("notification", hook.Payload{
		"hook_event_name":        "Stop",
		"cwd":                    "/Users/fred/project/demo",
		"last_assistant_message": "finished",
		"notification_type":      "",
	}, "Codex")
	if data.Event != "stop" {
		t.Fatalf("event = %q", data.Event)
	}
	if data.Title != "Codex stop | project/demo" {
		t.Fatalf("title = %q", data.Title)
	}
	if data.Body != "finished" {
		t.Fatalf("body = %q", data.Body)
	}
}

func TestBuildPermissionNotification(t *testing.T) {
	data := Build("notification", hook.Payload{
		"notification_type": "permission_prompt",
		"cwd":               "/repo/app",
		"message":           "approve command",
	}, "Claude")
	if data.Event != "permission" {
		t.Fatalf("event = %q", data.Event)
	}
	if data.Title != "Claude 需要权限确认 | repo/app" {
		t.Fatalf("title = %q", data.Title)
	}
}
