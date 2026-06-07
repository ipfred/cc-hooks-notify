package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ipfred/cc-hooks-notify/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage cc-notify config",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Create or update config interactively",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Root().PersistentFlags().GetString("config")
			return runConfigInit(path)
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show redacted config",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Root().PersistentFlags().GetString("config")
			cfg, loaded, err := config.Load(path)
			if err != nil {
				return err
			}
			if loaded == "" {
				loaded = config.DefaultPath()
			}
			data, _ := json.MarshalIndent(config.Redacted(cfg), "", "  ")
			fmt.Printf("Config: %s\n%s\n", loaded, data)
			return nil
		},
	})
	return cmd
}

func runConfigInit(path string) error {
	if path == "" {
		path = config.DefaultPath()
	}
	cfg, _, _ := config.Load(path)
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Config path: %s\n", path)

	dingWebhook := prompt(reader, "DingTalk webhook", cfg.Channels["dingtalk"].Webhook)
	dingSecret := prompt(reader, "DingTalk secret", cfg.Channels["dingtalk"].Secret)
	feishuWebhook := prompt(reader, "Feishu webhook", cfg.Channels["feishu"].Webhook)
	feishuSecret := prompt(reader, "Feishu secret", cfg.Channels["feishu"].Secret)

	if cfg.Channels == nil {
		cfg.Channels = map[string]config.ChannelConfig{}
	}
	markdown := true
	card := true
	cfg.Channels["dingtalk"] = config.ChannelConfig{Webhook: dingWebhook, Secret: dingSecret, Markdown: &markdown}
	cfg.Channels["feishu"] = config.ChannelConfig{Webhook: feishuWebhook, Secret: feishuSecret, Card: &card}
	if cfg.Events == nil {
		cfg.Events = map[string]config.EventConfig{}
	}
	cfg.Events["notification"] = config.EventConfig{Channels: enabledChannelNames(cfg)}
	cfg.Events["stop"] = config.EventConfig{Channels: enabledChannelNames(cfg)}
	if err := config.Save(path, cfg); err != nil {
		return err
	}
	fmt.Printf("Saved config: %s\n", path)
	return nil
}

func prompt(reader *bufio.Reader, label, current string) string {
	if current == "" {
		fmt.Printf("%s: ", label)
	} else {
		fmt.Printf("%s [%s]: ", label, current)
	}
	line, _ := reader.ReadString('\n')
	value := strings.TrimSpace(line)
	if value == "" {
		return current
	}
	return value
}

func enabledChannelNames(cfg config.Config) []string {
	names := []string{}
	for name, ch := range cfg.Channels {
		if strings.TrimSpace(ch.Webhook) != "" {
			names = append(names, name)
		}
	}
	return names
}
