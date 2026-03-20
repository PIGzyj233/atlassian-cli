package search

import (
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdSearch creates the `search` command.
func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	var (
		limit        int
		spacesFilter string
	)

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search Confluence content using CQL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			// If query doesn't look like CQL, wrap as text search
			cql := query
			if !isCQL(query) {
				cql = `text ~ "` + query + `"`
			}

			// Apply spaces filter
			if spacesFilter != "" {
				spaces := strings.Split(spacesFilter, ",")
				var spaceParts []string
				for _, s := range spaces {
					s = strings.TrimSpace(s)
					if s != "" {
						spaceParts = append(spaceParts, `space = "`+s+`"`)
					}
				}
				if len(spaceParts) > 0 {
					cql = "(" + cql + ") AND (" + strings.Join(spaceParts, " OR ") + ")"
				}
			}

			path := client.ConfluenceAPIPath("content/search")
			rb := api.NewRequestBuilder(path).
				Query("cql", cql).
				QueryInt("limit", limit)

			var result map[string]any
			_, err = client.Get(rb.BuildPath(), &result)
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum results to return")
	cmd.Flags().StringVar(&spacesFilter, "spaces-filter", "", "Comma-separated space keys to filter")

	return cmd
}

// isCQL checks if a query looks like CQL syntax.
func isCQL(query string) bool {
	cqlKeywords := []string{"=", "~", " AND ", " OR ", " IN ", " NOT "}
	for _, kw := range cqlKeywords {
		if strings.Contains(query, kw) {
			return true
		}
	}
	return false
}
