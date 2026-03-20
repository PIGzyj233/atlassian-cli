package link

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdCreateRemote creates the `link create-remote` command.
func NewCmdCreateRemote(f *cmdutil.Factory) *cobra.Command {
	var (
		url     string
		title   string
		summary string
	)

	cmd := &cobra.Command{
		Use:   "create-remote <issue-key>",
		Short: "Create a remote link on an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := createRemoteLink(client, args[0], url, title, summary)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "Remote link URL (required)")
	cmd.Flags().StringVar(&title, "title", "", "Remote link title (required)")
	cmd.Flags().StringVar(&summary, "summary", "", "Remote link summary")
	_ = cmd.MarkFlagRequired("url")
	_ = cmd.MarkFlagRequired("title")

	return cmd
}

func createRemoteLink(client *api.Client, issueKey, url, title, summary string) (map[string]any, error) {
	object := map[string]any{
		"url":   url,
		"title": title,
	}
	if summary != "" {
		object["summary"] = summary
	}

	body := map[string]any{
		"object": object,
	}

	path := client.JiraAPIPath("issue/" + issueKey + "/remotelink")

	var result map[string]any
	_, err := client.Post(path, body, &result)
	if err != nil {
		return nil, fmt.Errorf("creating remote link on %s: %w", issueKey, err)
	}
	return result, nil
}
