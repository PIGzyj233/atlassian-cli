package comment

import (
	"github.com/PigZyj2333/atlassian-cli/internal/convert"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdReply creates the `comment reply` command.
func NewCmdReply(f *cmdutil.Factory) *cobra.Command {
	var body string

	cmd := &cobra.Command{
		Use:   "reply <comment-id>",
		Short: "Reply to an existing comment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			// Convert markdown to storage format
			storageBody := convert.MarkdownToStorage(body)

			payload := map[string]any{
				"type": "comment",
				"container": map[string]any{
					"id":   args[0],
					"type": "comment",
				},
				"body": map[string]any{
					"storage": map[string]any{
						"value":          storageBody,
						"representation": "storage",
					},
				},
			}

			path := client.ConfluenceAPIPath("content")
			var result map[string]any
			_, err = client.Post(path, payload, &result)
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&body, "body", "", "Reply body in Markdown (required)")
	cmd.MarkFlagRequired("body")

	return cmd
}
