package jira

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/attachment"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/auth"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/board"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/changelog"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/comment"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/dev"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/epic"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/field"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/form"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/issue"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/link"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/project"
	searchcmd "github.com/PigZyj2333/atlassian-cli/pkg/jira/search"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/servicedesk"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/sla"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/sprint"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/transition"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/user"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/watcher"
	"github.com/PigZyj2333/atlassian-cli/pkg/jira/worklog"
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

	cmd.AddCommand(attachment.NewCmdAttachment(f))
	cmd.AddCommand(auth.NewCmdAuth(f))
	cmd.AddCommand(issue.NewCmdIssue(f))
	cmd.AddCommand(searchcmd.NewCmdSearch(f))
	cmd.AddCommand(changelog.NewCmdChangelog(f))
	cmd.AddCommand(comment.NewCmdComment(f))
	cmd.AddCommand(dev.NewCmdDev(f))
	cmd.AddCommand(epic.NewCmdEpic(f))
	cmd.AddCommand(field.NewCmdField(f))
	cmd.AddCommand(form.NewCmdForm(f))
	cmd.AddCommand(transition.NewCmdTransition(f))
	cmd.AddCommand(board.NewCmdBoard(f))
	cmd.AddCommand(sprint.NewCmdSprint(f))
	cmd.AddCommand(link.NewCmdLink(f))
	cmd.AddCommand(sla.NewCmdSLA(f))
	cmd.AddCommand(watcher.NewCmdWatcher(f))
	cmd.AddCommand(worklog.NewCmdWorklog(f))
	cmd.AddCommand(user.NewCmdUser(f))
	cmd.AddCommand(project.NewCmdProject(f))
	cmd.AddCommand(servicedesk.NewCmdServiceDesk(f))

	return cmd
}
