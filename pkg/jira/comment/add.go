package comment

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/convert"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAdd creates the `comment add` command.
func NewCmdAdd(f *cmdutil.Factory) *cobra.Command {
	var (
		body       string
		visibility string
		public     bool
	)

	cmd := &cobra.Command{
		Use:   "add <issue-key>",
		Short: "Add a comment to an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate mutual exclusivity (ref: Python add_comment)
			if cmd.Flags().Changed("public") && visibility != "" {
				return fmt.Errorf("cannot use both --public and --visibility; --public uses the ServiceDesk API which does not support Jira visibility restrictions")
			}

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			// ServiceDesk API path for internal/public comments
			if cmd.Flags().Changed("public") {
				result, err := addServiceDeskComment(client, args[0], body, public)
				if err != nil {
					return err
				}
				return output.Print(cmd, result)
			}

			result, err := addComment(client, args[0], body, visibility)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&body, "body", "", "Comment body (required)")
	cmd.Flags().StringVar(&visibility, "visibility", "", "Visibility role (restricts who can see the comment)")
	cmd.Flags().BoolVar(&public, "public", false, "Make comment public/customer-visible (Service Desk issues only)")
	_ = cmd.MarkFlagRequired("body")
	return cmd
}

func addComment(client *api.Client, issueKey, body, visibility string) (map[string]any, error) {
	// Cloud (v3) requires ADF; Server/DC (v2) uses plain string.
	var bodyValue any = body
	if client.IsCloud() {
		bodyValue = convert.MarkdownToADF(body)
	}

	payload := map[string]any{"body": bodyValue}
	if visibility != "" {
		payload["visibility"] = map[string]any{
			"type":  "role",
			"value": visibility,
		}
	}

	path := client.JiraAPIPath("issue/" + issueKey + "/comment")
	var result map[string]any
	_, err := client.Post(path, payload, &result)
	if err != nil {
		return nil, fmt.Errorf("adding comment to %s: %w", issueKey, err)
	}
	return result, nil
}

// addServiceDeskComment uses the ServiceDesk API for public/internal comments.
// Ref: Python _add_servicedesk_comment() — uses X-ExperimentalApi header.
func addServiceDeskComment(client *api.Client, issueKey, body string, public bool) (map[string]any, error) {
	payload := map[string]any{
		"body":   body,
		"public": public,
	}

	path := "/rest/servicedeskapi/request/" + issueKey + "/comment"

	// The ServiceDesk API requires the X-ExperimentalApi header
	headers := map[string]string{"X-ExperimentalApi": "opt-in"}

	var result map[string]any
	_, err := client.PostWithHeaders(path, payload, &result, headers)
	if err != nil {
		return nil, fmt.Errorf("adding ServiceDesk comment to %s: %w", issueKey, err)
	}
	return result, nil
}
