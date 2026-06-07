package cli

import (
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:           "cc-notify",
		Short:         "Send Codex, Claude, and Droid hook notifications",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.PersistentFlags().String("config", "", "config file path")
	root.AddCommand(newHookCommand())
	root.AddCommand(newConfigCommand())
	root.AddCommand(newTestCommand())
	root.AddCommand(newInstallCommand())
	root.AddCommand(newVersionCommand())
	return root
}
