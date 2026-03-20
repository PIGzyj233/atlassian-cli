package analytics

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAnalytics creates the `analytics` subcommand group.
func NewCmdAnalytics(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analytics <command>",
		Short: "Confluence analytics (Cloud only)",
	}
	cmd.AddCommand(NewCmdViews(f))
	return cmd
}
