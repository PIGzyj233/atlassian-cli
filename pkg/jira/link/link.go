package link

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdLink creates the `link` subcommand group.
func NewCmdLink(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link <command>",
		Short: "Manage issue links",
	}
	cmd.AddCommand(NewCmdTypes(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdCreateRemote(f))
	cmd.AddCommand(NewCmdRemove(f))
	return cmd
}
