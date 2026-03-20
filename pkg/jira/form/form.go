package form

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdForm creates the `form` subcommand group.
func NewCmdForm(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "form <command>",
		Short: "Manage ProForma forms on issues",
	}
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdUpdate(f))
	return cmd
}
