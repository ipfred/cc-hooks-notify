package hook

import (
	"strings"
	"testing"
)

func TestParseValidJSON(t *testing.T) {
	payload, err := Parse(strings.NewReader(`{"hook_event_name":"Stop","turn_id":"t","cwd":"/a/b"}`))
	if err != nil {
		t.Fatal(err)
	}
	if String(payload, "hook_event_name") != "Stop" {
		t.Fatalf("unexpected hook_event_name: %v", payload)
	}
	if !IsCodexStop(payload) {
		t.Fatal("expected Codex Stop payload")
	}
}

func TestParseInvalidJSONReturnsParseError(t *testing.T) {
	payload, err := Parse(strings.NewReader(`{bad`))
	if err != nil {
		t.Fatal(err)
	}
	if String(payload, "parse_error") == "" {
		t.Fatalf("expected parse_error: %v", payload)
	}
}

func TestProjectName(t *testing.T) {
	got := ProjectName(`/Users/fred/project/demo`)
	if got != "project/demo" {
		t.Fatalf("ProjectName() = %q", got)
	}
}
