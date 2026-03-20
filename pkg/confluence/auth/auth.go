package auth

import (
	jiraauth "github.com/PigZyj2333/atlassian-cli/pkg/jira/auth"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAuth creates the `auth` subcommand group (shared with jira).
func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	// Reuse jira auth — same config file, same commands
	return jiraauth.NewCmdAuth(f)
}
