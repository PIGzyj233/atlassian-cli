package search

import (
	"fmt"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// SearchOpts holds search parameters.
type SearchOpts struct {
	JQL            string
	Fields         string
	Limit          int
	StartAt        int
	ProjectsFilter string
	Expand         string
}

// NewCmdSearch creates the `search` command.
func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	var opts SearchOpts

	cmd := &cobra.Command{
		Use:   "search <jql>",
		Short: "Search Jira issues using JQL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.JQL = args[0]

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := searchIssues(client, opts)
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&opts.Fields, "fields", "", "Comma-separated fields to return")
	cmd.Flags().IntVar(&opts.Limit, "limit", 50, "Maximum results to return")
	cmd.Flags().IntVar(&opts.StartAt, "start-at", 0, "Starting index for pagination")
	cmd.Flags().StringVar(&opts.ProjectsFilter, "projects-filter", "", "Comma-separated project keys to filter")
	cmd.Flags().StringVar(&opts.Expand, "expand", "", "Fields to expand")

	return cmd
}

func searchIssues(client *api.Client, opts SearchOpts) (map[string]any, error) {
	if client.IsCloud() {
		return searchCloud(client, opts)
	}
	return searchServerDC(client, opts)
}

func searchCloud(client *api.Client, opts SearchOpts) (map[string]any, error) {
	body := map[string]any{
		"jql":        buildSearchJQL(opts.JQL, opts.ProjectsFilter),
		"maxResults": opts.Limit,
	}
	if opts.Fields != "" {
		body["fields"] = opts.Fields
	}
	if opts.Expand != "" {
		body["expand"] = opts.Expand
	}

	path := client.JiraAPIPath("search/jql")
	var result map[string]any
	_, err := client.Post(path, body, &result)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	return result, nil
}

func searchServerDC(client *api.Client, opts SearchOpts) (map[string]any, error) {
	rb := api.NewRequestBuilder(client.JiraAPIPath("search")).
		Query("jql", buildSearchJQL(opts.JQL, opts.ProjectsFilter)).
		Query("fields", opts.Fields).
		Query("expand", opts.Expand)

	if opts.Limit > 0 {
		rb.QueryInt("maxResults", opts.Limit)
	}
	if opts.StartAt > 0 {
		rb.QueryInt("startAt", opts.StartAt)
	}

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	return result, nil
}

func buildSearchJQL(jql, projectsFilter string) string {
	if strings.TrimSpace(projectsFilter) == "" {
		return jql
	}
	return fmt.Sprintf("(%s) AND project in (%s)", jql, projectsFilter)
}
