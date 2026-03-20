package transition

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdTransition creates the `transition` subcommand group.
func NewCmdTransition(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transition <command>",
		Short: "Manage issue transitions",
	}
	cmd.AddCommand(NewCmdList(f))
	return cmd
}
