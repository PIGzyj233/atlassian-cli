package page

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdTree creates the `page tree` command.
func NewCmdTree(f *cmdutil.Factory) *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "tree <space-key>",
		Short: "Get hierarchical page tree for a space",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			spaceKey := args[0]
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			// Paginate to fetch all pages with ancestors
			pageSize := 200
			start := 0
			var allPages []any
			hasMore := false

			for len(allPages) < limit {
				fetchLimit := pageSize
				if remaining := limit - len(allPages); remaining < fetchLimit {
					fetchLimit = remaining
				}

				path := client.ConfluenceAPIPath("content")
				rb := api.NewRequestBuilder(path).
					Query("spaceKey", spaceKey).
					Query("type", "page").
					Query("expand", "ancestors").
					QueryInt("start", start).
					QueryInt("limit", fetchLimit)

				var result map[string]any
				_, err := client.Get(rb.BuildPath(), &result)
				if err != nil {
					return fmt.Errorf("fetching pages: %w", err)
				}

				batch, ok := result["results"].([]any)
				if !ok || len(batch) == 0 {
					break
				}
				allPages = append(allPages, batch...)

				// Check for next link
				links, _ := result["_links"].(map[string]any)
				if _, hasNext := links["next"]; !hasNext {
					break
				}
				hasMore = true
				start += len(batch)
			}

			// Build flat list with parent_id and depth
			var resultPages []map[string]any
			for _, pageAny := range allPages {
				page, ok := pageAny.(map[string]any)
				if !ok {
					continue
				}

				pageID, _ := page["id"].(string)
				title, _ := page["title"].(string)

				// Position from extensions
				var position any
				if ext, ok := page["extensions"].(map[string]any); ok {
					position = ext["position"]
				}

				// Parent and depth from ancestors
				var parentID string
				depth := 0
				if ancestors, ok := page["ancestors"].([]any); ok && len(ancestors) > 0 {
					depth = len(ancestors)
					if lastAncestor, ok := ancestors[len(ancestors)-1].(map[string]any); ok {
						parentID, _ = lastAncestor["id"].(string)
					}
				}

				entry := map[string]any{
					"id":        pageID,
					"title":     title,
					"parent_id": parentID,
					"position":  position,
					"depth":     depth,
				}
				resultPages = append(resultPages, entry)
			}

			treeResult := map[string]any{
				"space_key":   spaceKey,
				"total_pages": len(resultPages),
				"has_more":    hasMore && len(allPages) >= limit,
				"pages":       resultPages,
			}

			return output.Print(cmd, treeResult)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 500, "Maximum pages to fetch")

	return cmd
}
