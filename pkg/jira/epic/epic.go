package epic

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdEpic creates the `epic` subcommand group.
func NewCmdEpic(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "epic <command>",
		Short: "Manage epics",
	}
	cmd.AddCommand(NewCmdLink(f))
	return cmd
}
