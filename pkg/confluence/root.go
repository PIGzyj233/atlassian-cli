package confluence

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/PigZyj2333/atlassian-cli/pkg/confluence/analytics"
	"github.com/PigZyj2333/atlassian-cli/pkg/confluence/attachment"
	"github.com/PigZyj2333/atlassian-cli/pkg/confluence/auth"
	"github.com/PigZyj2333/atlassian-cli/pkg/confluence/comment"
	"github.com/PigZyj2333/atlassian-cli/pkg/confluence/label"
	"github.com/PigZyj2333/atlassian-cli/pkg/confluence/page"
	"github.com/PigZyj2333/atlassian-cli/pkg/confluence/search"
	"github.com/PigZyj2333/atlassian-cli/pkg/confluence/user"
	"github.com/spf13/cobra"
)

// NewCmdRoot creates the root `confluence` command.
func NewCmdRoot(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "confluence <command> [flags]",
		Short:         "Work with Confluence from the command line",
		Long:          "A CLI tool for interacting with Atlassian Confluence, designed for both humans and LLM automation.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().String("host", "", "Target Atlassian instance hostname")
	cmd.PersistentFlags().String("output", "json", "Output format: json, text, table")
	cmd.PersistentFlags().Bool("verbose", false, "Enable debug logging to stderr")
	cmd.PersistentFlags().Bool("no-color", false, "Disable ANSI colors")

	// Register all subcommand groups
	cmd.AddCommand(auth.NewCmdAuth(f))
	cmd.AddCommand(page.NewCmdPage(f))
	cmd.AddCommand(search.NewCmdSearch(f))
	cmd.AddCommand(comment.NewCmdComment(f))
	cmd.AddCommand(label.NewCmdLabel(f))
	cmd.AddCommand(attachment.NewCmdAttachment(f))
	cmd.AddCommand(user.NewCmdUser(f))
	cmd.AddCommand(analytics.NewCmdAnalytics(f))

	cmd.Version = f.Version

	return cmd
}
