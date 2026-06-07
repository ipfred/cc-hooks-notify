package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

func newInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install hook integrations",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "codex",
		Short: "Install Codex plugin hooks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return installCodexPlugin()
		},
	})
	return cmd
}

func installCodexPlugin() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	root := filepath.Join(home, ".cc-notify", "codex-marketplace")
	pluginRoot := filepath.Join(root, "plugins", "cc-notify")
	if err := os.MkdirAll(filepath.Join(pluginRoot, ".codex-plugin"), 0o700); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(pluginRoot, "hooks"), 0o700); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(root, ".agents", "plugins"), 0o700); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(pluginRoot, ".codex-plugin", "plugin.json"), codexPluginManifest()); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(pluginRoot, "hooks", "hooks.json"), codexHooksConfig()); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(root, ".agents", "plugins", "marketplace.json"), codexMarketplace()); err != nil {
		return err
	}
	if _, err := exec.LookPath("codex"); err != nil {
		fmt.Printf("Codex plugin files written to %s\n", root)
		fmt.Println("Install manually: codex plugin marketplace add " + root)
		fmt.Println("Then run: codex plugin add cc-notify@cc-notify-local")
		return nil
	}
	_ = exec.Command("codex", "plugin", "marketplace", "remove", "cc-notify-local").Run()
	if output, err := exec.Command("codex", "plugin", "marketplace", "add", root).CombinedOutput(); err != nil {
		return fmt.Errorf("codex plugin marketplace add failed: %w\n%s", err, string(output))
	}
	if output, err := exec.Command("codex", "plugin", "add", "cc-notify@cc-notify-local").CombinedOutput(); err != nil {
		return fmt.Errorf("codex plugin add failed: %w\n%s", err, string(output))
	}
	fmt.Println("Installed Codex plugin cc-notify@cc-notify-local")
	fmt.Println("Open a new Codex session and run /hooks to trust the hook.")
	return nil
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

func codexPluginManifest() map[string]any {
	return map[string]any{
		"name":        "cc-notify",
		"version":     pluginVersion(),
		"description": "Codex notification hooks for DingTalk and Feishu.",
		"author":      map[string]string{"name": "fred"},
		"repository":  "https://github.com/ipfred/cc-hooks-notify",
		"license":     "MIT",
		"keywords":    []string{"notification", "dingtalk", "feishu", "hooks", "webhook", "codex"},
		"interface": map[string]any{
			"displayName":      "cc-notify",
			"shortDescription": "Send Codex hook notifications to DingTalk or Feishu.",
			"longDescription":  "Bundles Codex lifecycle hooks that call the globally installed cc-notify command.",
			"developerName":    "fred",
			"category":         "Productivity",
			"capabilities":     []string{"Hooks"},
			"defaultPrompt":    []string{"Use cc-notify for Codex lifecycle notifications."},
			"brandColor":       "#2563EB",
		},
	}
}

func pluginVersion() string {
	if regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)`).MatchString(version) {
		return version
	}
	return "0.4.0"
}

func codexHooksConfig() map[string]any {
	handler := map[string]any{
		"type":          "command",
		"command":       "cc-notify hook codex",
		"timeout":       10,
		"statusMessage": "cc-notify sending",
	}
	return map[string]any{
		"description": "cc-notify Codex Hooks",
		"hooks": map[string]any{
			"PermissionRequest": []map[string]any{{"matcher": "*", "hooks": []map[string]any{handler}}},
			"Stop":              []map[string]any{{"hooks": []map[string]any{handler}}},
		},
	}
}

func codexMarketplace() map[string]any {
	return map[string]any{
		"name":      "cc-notify-local",
		"interface": map[string]string{"displayName": "cc-notify local"},
		"plugins": []map[string]any{
			{
				"name":     "cc-notify",
				"source":   map[string]string{"source": "local", "path": "./plugins/cc-notify"},
				"policy":   map[string]string{"installation": "AVAILABLE", "authentication": "ON_INSTALL"},
				"category": "Productivity",
			},
		},
	}
}
