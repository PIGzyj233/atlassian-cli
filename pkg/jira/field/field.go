package field

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdField creates the `field` subcommand group.
func NewCmdField(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "field <command>",
		Short: "Discover and inspect Jira fields",
	}
	cmd.AddCommand(NewCmdSearch(f))
	cmd.AddCommand(NewCmdOptions(f))
	return cmd
}
