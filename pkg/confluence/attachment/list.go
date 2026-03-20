package attachment

import (
	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdList creates the `attachment list` command.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	var (
		limit     int
		filename  string
		mediaType string
	)

	cmd := &cobra.Command{
		Use:   "list <content-id>",
		Short: "List attachments on Confluence content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			path := client.ConfluenceAPIPath("content/" + args[0] + "/child/attachment")
			rb := api.NewRequestBuilder(path).
				QueryInt("limit", limit)

			if filename != "" {
				rb.Query("filename", filename)
			}
			if mediaType != "" {
				rb.Query("mediaType", mediaType)
			}

			var result map[string]any
			_, err = client.Get(rb.BuildPath(), &result)
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum results to return")
	cmd.Flags().StringVar(&filename, "filename", "", "Filter by filename")
	cmd.Flags().StringVar(&mediaType, "media-type", "", "Filter by media type")

	return cmd
}
