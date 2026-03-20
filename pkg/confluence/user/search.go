package user

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdSearch creates the `user search` command.
func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	var (
		limit int
		group string
	)

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search Confluence users",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			var result map[string]any

			if client.IsCloud() {
				// Cloud: use CQL search endpoint
				cql := `user.fullname ~ "` + query + `"`
				path := api.NewRequestBuilder("/rest/api/search/user").
					Query("cql", cql).
					QueryInt("limit", limit).
					BuildPath()

				_, err = client.Get(path, &result)
			} else {
				// Server/DC: search via group member API with client-side filtering.
				// Matches Python reference: search.py:138-195 (_search_user_server_dc).
				members, searchErr := searchGroupMembers(client, query, group, limit)
				if searchErr != nil {
					return searchErr
				}
				result = map[string]any{"results": members, "size": len(members)}
			}

			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 10, "Maximum results to return")
	cmd.Flags().StringVar(&group, "group", "confluence-users", "Group to search within (Server/DC)")

	return cmd
}

// searchGroupMembers fetches group members and filters by query, with pagination.
func searchGroupMembers(client *api.Client, query, group string, limit int) ([]any, error) {
	queryLower := strings.ToLower(query)
	var matches []any
	start := 0
	pageSize := 200

	for len(matches) < limit {
		path := api.NewRequestBuilder(
			fmt.Sprintf("/rest/api/group/%s/member", url.PathEscape(group)),
		).
			QueryInt("start", start).
			QueryInt("limit", pageSize).
			BuildPath()

		var page map[string]any
		_, err := client.Get(path, &page)
		if err != nil {
			return nil, err
		}

		members, _ := page["results"].([]any)
		for _, m := range members {
			member, ok := m.(map[string]any)
			if !ok {
				continue
			}
			display, _ := member["displayName"].(string)
			username, _ := member["username"].(string)
			if strings.Contains(strings.ToLower(display), queryLower) ||
				strings.Contains(strings.ToLower(username), queryLower) {
				matches = append(matches, m)
				if len(matches) >= limit {
					break
				}
			}
		}

		// Check for pagination
		links, _ := page["_links"].(map[string]any)
		if _, hasNext := links["next"]; !hasNext || len(members) == 0 {
			break
		}
		start += len(members)
	}

	return matches, nil
}
