package attachment

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAttachment creates the `attachment` subcommand group.
func NewCmdAttachment(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attachment <command>",
		Short: "Manage Confluence attachments",
	}
	cmd.AddCommand(NewCmdList(f))
	cmd.AddCommand(NewCmdUpload(f))
	cmd.AddCommand(NewCmdUploadBatch(f))
	cmd.AddCommand(NewCmdDownload(f))
	cmd.AddCommand(NewCmdDownloadAll(f))
	cmd.AddCommand(NewCmdDelete(f))
	cmd.AddCommand(NewCmdImages(f))
	return cmd
}
