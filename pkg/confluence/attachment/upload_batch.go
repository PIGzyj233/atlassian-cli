package attachment

import (
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdUploadBatch creates the `attachment upload-batch` command.
func NewCmdUploadBatch(f *cmdutil.Factory) *cobra.Command {
	var (
		files     []string
		comment   string
		minorEdit bool
	)

	cmd := &cobra.Command{
		Use:   "upload-batch <content-id>",
		Short: "Upload multiple attachments to Confluence content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contentID := args[0]
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			var uploaded []map[string]any
			var failed []map[string]any

			for _, filePath := range files {
				result, err := uploadFile(client.BaseURL, client.Auth, contentID, filePath, comment, minorEdit)
				if err != nil {
					failed = append(failed, map[string]any{
						"filename": filePath,
						"error":    err.Error(),
					})
					continue
				}
				uploaded = append(uploaded, result)
			}

			return output.Print(cmd, map[string]any{
				"success":    true,
				"content_id": contentID,
				"total":      len(files),
				"uploaded":   uploaded,
				"failed":     failed,
			})
		},
	}

	cmd.Flags().StringSliceVar(&files, "files", nil, "Paths to files to upload (required)")
	cmd.Flags().StringVar(&comment, "comment", "", "Attachment comment")
	cmd.Flags().BoolVar(&minorEdit, "minor-edit", true, "Mark as minor edit")
	cmd.MarkFlagRequired("files")

	return cmd
}
