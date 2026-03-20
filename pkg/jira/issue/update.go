package issue

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var (
		fieldsJSON  string
		components  string
		attachments string
	)

	cmd := &cobra.Command{
		Use:   "update <issue-key>",
		Short: "Update a Jira issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fieldsJSON == "" && components == "" && attachments == "" {
				return fmt.Errorf("at least one of --fields-json, --components, or --attachments is required")
			}
			if attachments != "" {
				return fmt.Errorf("--attachments is not implemented yet")
			}

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			fields := make(map[string]any)
			if fieldsJSON != "" {
				if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
					return fmt.Errorf("parsing --fields-json: %w", err)
				}
			}
			if components != "" {
				fields["components"] = parseComponentNames(components)
			}

			if err := updateIssue(client, args[0], fields); err != nil {
				return err
			}

			return output.Print(cmd, map[string]any{
				"success": true,
				"key":     args[0],
				"message": "Issue updated",
			})
		},
	}

	cmd.Flags().StringVar(&fieldsJSON, "fields-json", "", "Fields to update as JSON")
	cmd.Flags().StringVar(&components, "components", "", "Comma-separated component names")
	cmd.Flags().StringVar(&attachments, "attachments", "", "Comma-separated file paths to attach")

	return cmd
}

func updateIssue(client *api.Client, issueKey string, fields map[string]any) error {
	body := map[string]any{"fields": fields}
	path := client.JiraAPIPath("issue/" + issueKey)
	_, err := client.Put(path, body, nil)
	if err != nil {
		return fmt.Errorf("updating issue %s: %w", issueKey, err)
	}
	return nil
}

func parseComponentNames(value string) []map[string]any {
	components := make([]map[string]any, 0)
	for _, name := range strings.Split(value, ",") {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		components = append(components, map[string]any{"name": trimmed})
	}
	return components
}
