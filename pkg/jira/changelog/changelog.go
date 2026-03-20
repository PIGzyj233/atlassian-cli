package changelog

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdChangelog creates the `changelog` subcommand group.
func NewCmdChangelog(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changelog <command>",
		Short: "View issue changelog",
	}
	cmd.AddCommand(NewCmdBatch(f))
	return cmd
}
