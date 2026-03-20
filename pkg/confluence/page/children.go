package page

import (
	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdChildren creates the `page children` command.
func NewCmdChildren(f *cmdutil.Factory) *cobra.Command {
	var (
		limit          int
		includeContent bool
		includeFolders bool
	)

	cmd := &cobra.Command{
		Use:   "children <page-id>",
		Short: "List child pages of a Confluence page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			expand := "version"
			if includeContent {
				expand = "body.storage,version"
			}

			path := client.ConfluenceAPIPath("content/" + args[0] + "/child/page")
			rb := api.NewRequestBuilder(path).
				Query("expand", expand)
			if limit > 0 {
				rb.QueryInt("limit", limit)
			}

			var result map[string]any
			_, err = client.Get(rb.BuildPath(), &result)
			if err != nil {
				return err
			}

			// Also fetch folders if requested
			if includeFolders {
				folderPath := client.ConfluenceAPIPath("content/" + args[0] + "/child/folder")
				frb := api.NewRequestBuilder(folderPath).
					Query("expand", expand)
				if limit > 0 {
					frb.QueryInt("limit", limit)
				}
				var folderResult map[string]any
				if _, ferr := client.Get(frb.BuildPath(), &folderResult); ferr == nil {
					// Merge folder results into page results
					if folders, ok := folderResult["results"].([]any); ok {
						if pages, ok := result["results"].([]any); ok {
							result["results"] = append(pages, folders...)
						}
					}
				}
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 25, "Maximum results to return")
	cmd.Flags().BoolVar(&includeContent, "include-content", false, "Include page content")
	cmd.Flags().BoolVar(&includeFolders, "include-folders", true, "Include child folders")

	return cmd
}
