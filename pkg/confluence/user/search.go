package user

import (
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
				// Server/DC: use user search API
				path := api.NewRequestBuilder("/rest/api/user/search").
					Query("username", query).
					QueryInt("maxResults", limit).
					BuildPath()

				_, err = client.Get(path, &result)
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
