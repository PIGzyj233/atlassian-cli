package jira

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/auth"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/board"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/comment"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/field"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/issue"
	searchcmd "github.com/PigZyj2333/atlassian-cli/pkg/jira/search"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/sprint"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/transition"
	"github.com/spf13/cobra"
)

// NewCmdRoot creates the root `jira` command.
func NewCmdRoot(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "jira <command> [flags]",
		Short:         "Work with Jira from the command line",
		Long:          "A CLI tool for interacting with Atlassian Jira, designed for both humans and LLM automation.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().String("host", "", "Target Atlassian instance hostname")
	cmd.PersistentFlags().String("output", "json", "Output format: json, text, table")
	cmd.PersistentFlags().Bool("verbose", false, "Enable debug logging to stderr")
	cmd.PersistentFlags().Bool("no-color", false, "Disable ANSI colors")

	cmd.Version = f.Version

	cmd.AddCommand(auth.NewCmdAuth(f))
	cmd.AddCommand(issue.NewCmdIssue(f))
	cmd.AddCommand(searchcmd.NewCmdSearch(f))
	cmd.AddCommand(comment.NewCmdComment(f))
	cmd.AddCommand(field.NewCmdField(f))
	cmd.AddCommand(transition.NewCmdTransition(f))
	cmd.AddCommand(board.NewCmdBoard(f))
	cmd.AddCommand(sprint.NewCmdSprint(f))

	return cmd
}
