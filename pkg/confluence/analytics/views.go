package analytics

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdViews creates the `analytics views` command.
func NewCmdViews(f *cmdutil.Factory) *cobra.Command {
	var includeTitle bool

	cmd := &cobra.Command{
		Use:   "views <page-id>",
		Short: "Get page view statistics (Cloud only)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pageID := args[0]
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			if !client.IsCloud() {
				return fmt.Errorf("analytics views is only available for Confluence Cloud")
			}

			// Use the analytics API endpoint
			path := "/rest/api/analytics/content/" + pageID + "/views"
			var viewsData map[string]any
			_, err = client.Get(path, &viewsData)
			if err != nil {
				return fmt.Errorf("getting page views: %w", err)
			}

			result := map[string]any{
				"page_id":     pageID,
				"total_views": viewsData["count"],
				"last_seen":   viewsData["lastSeen"],
			}

			// Optionally fetch page title
			if includeTitle {
				pagePath := client.ConfluenceAPIPath("content/" + pageID)
				var pageData map[string]any
				if _, err := client.Get(pagePath, &pageData); err == nil {
					result["title"] = pageData["title"]
				}
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().BoolVar(&includeTitle, "include-title", true, "Include page title")

	return cmd
}
