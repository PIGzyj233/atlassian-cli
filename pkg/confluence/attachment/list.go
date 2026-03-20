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

			var result map[string]any
			_, err = client.Get(rb.BuildPath(), &result)
			if err != nil {
				return err
			}

			// V1 API doesn't support server-side filtering for filename/mediaType.
			// Apply client-side filtering (matches Python reference: attachments.py:399-420).
			if filename != "" || mediaType != "" {
				result = filterAttachments(result, filename, mediaType)
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 50, "Maximum results to return")
	cmd.Flags().StringVar(&filename, "filename", "", "Filter by filename (exact match, client-side)")
	cmd.Flags().StringVar(&mediaType, "media-type", "", "Filter by media type (exact match, client-side)")

	return cmd
}

// filterAttachments applies client-side filtering on attachment results.
func filterAttachments(result map[string]any, filename, mediaType string) map[string]any {
	results, ok := result["results"].([]any)
	if !ok {
		return result
	}

	var filtered []any
	for _, item := range results {
		att, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if filename != "" {
			if title, _ := att["title"].(string); title != filename {
				continue
			}
		}
		if mediaType != "" {
			if mt, _ := att["mediaType"].(string); mt != mediaType {
				continue
			}
		}
		filtered = append(filtered, item)
	}

	result["results"] = filtered
	result["size"] = len(filtered)
	return result
}
