package cli

import (
	"log/slog"

	"github.com/ipfred/cc-hooks-notify/internal/config"
	"github.com/ipfred/cc-hooks-notify/internal/hook"
	"github.com/ipfred/cc-hooks-notify/internal/logging"
	"github.com/ipfred/cc-hooks-notify/internal/notify"
	"github.com/spf13/cobra"
)

func newTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test [dingtalk|feishu|all]",
		Short: "Send a test notification",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Root().PersistentFlags().GetString("config")
			cfg, _, err := config.Load(path)
			if err != nil {
				return err
			}
			if len(args) == 1 && args[0] != "all" {
				cfg.Events["notification"] = config.EventConfig{Channels: []string{args[0]}}
			}
			logger, closeLogger, err := logging.New(cfg)
			if err != nil {
				return err
			}
			defer closeLogger()
			payload := hook.Payload{
				"hook_event_name": "Notification",
				"cwd":             ".",
				"message":         "cc-notify test message",
			}
			if err := notify.Send("notification", payload, cfg, "cc-notify", logger); err != nil {
				logger.Error("test notify failed", slog.Any("error", err))
			}
			return nil
		},
	}
	return cmd
}
