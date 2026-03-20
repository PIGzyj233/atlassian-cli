package watcher

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdWatcher creates the `watcher` subcommand group.
func NewCmdWatcher(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watcher <command>",
		Short: "Manage issue watchers",
	}
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdAdd(f))
	cmd.AddCommand(NewCmdRemove(f))
	return cmd
}
