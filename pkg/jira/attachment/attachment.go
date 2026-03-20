package attachment

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAttachment creates the `attachment` subcommand group.
func NewCmdAttachment(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachment <command>",
		Short: "Manage issue attachments",
	}
	cmd.AddCommand(NewCmdDownload(f))
	cmd.AddCommand(NewCmdImages(f))
	return cmd
}
