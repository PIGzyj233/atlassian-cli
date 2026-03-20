package issue

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdGet creates the `issue get` command.
func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var (
		fields       string
		expand       string
		commentLimit int
	)

	cmd := &cobra.Command{
		Use:   "get <issue-key>",
		Short: "Get a Jira issue by key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := fetchIssue(client, args[0], fields, expand)
			if err != nil {
				return err
			}

			_ = commentLimit
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&fields, "fields", "", "Comma-separated list of fields to include")
	cmd.Flags().StringVar(&expand, "expand", "", "Comma-separated list of fields to expand")
	cmd.Flags().IntVar(&commentLimit, "comment-limit", 0, "Max comments to include (0 = default)")

	return cmd
}

func fetchIssue(client *api.Client, issueKey, fields, expand string) (map[string]any, error) {
	path := client.JiraAPIPath("issue/" + issueKey)

	rb := api.NewRequestBuilder(path)
	if fields != "" {
		rb.Query("fields", fields)
	}
	if expand != "" {
		rb.Query("expand", expand)
	}

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("getting issue %s: %w", issueKey, err)
	}

	return result, nil
}
