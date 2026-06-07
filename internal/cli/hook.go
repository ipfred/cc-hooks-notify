package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ipfred/cc-hooks-notify/internal/config"
	"github.com/ipfred/cc-hooks-notify/internal/hook"
	"github.com/ipfred/cc-hooks-notify/internal/logging"
	"github.com/ipfred/cc-hooks-notify/internal/notify"
	"github.com/spf13/cobra"
)

func newHookCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hook",
		Short: "Run as a hook handler",
	}
	for _, mode := range []string{"codex", "claude", "droid"} {
		mode := mode
		cmd.AddCommand(&cobra.Command{
			Use:   mode,
			Short: fmt.Sprintf("Handle %s hook input from stdin", mode),
			RunE: func(cmd *cobra.Command, args []string) error {
				return runHook(cmd, mode)
			},
		})
	}
	return cmd
}

func runHook(cmd *cobra.Command, mode string) error {
	configPath, _ := cmd.Root().PersistentFlags().GetString("config")
	cfg, _, err := config.Load(configPath)
	if err != nil {
		return err
	}
	logger, closeLogger, err := logging.New(cfg)
	if err != nil {
		return err
	}
	defer closeLogger()

	payload, err := hook.Parse(os.Stdin)
	if err != nil {
		return err
	}
	logger.Info("hook input parsed", "mode", mode, "payload", payload)
	prefix := map[string]string{
		"codex":  "Codex",
		"claude": "Claude",
		"droid":  "Droid",
	}[mode]
	if err := notify.Send("notification", payload, cfg, prefix, logger); err != nil {
		logger.Error("notify failed", "error", err)
	}
	if mode == "codex" && hook.IsCodexStop(payload) {
		out, _ := json.Marshal(map[string]any{})
		fmt.Println(string(out))
	}
	return nil
}
