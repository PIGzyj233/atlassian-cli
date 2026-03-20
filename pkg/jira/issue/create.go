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

// CreateIssueOpts holds the parameters for creating an issue.
type CreateIssueOpts struct {
	Project     string
	Summary     string
	IssueType   string
	Assignee    string
	Description string
	Components  string
	FieldsJSON  string
}

// NewCmdCreate creates the `issue create` command.
func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var opts CreateIssueOpts

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Jira issue",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := createIssue(client, opts)
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&opts.Project, "project", "", "Project key (required)")
	cmd.Flags().StringVar(&opts.Summary, "summary", "", "Issue summary (required)")
	cmd.Flags().StringVar(&opts.IssueType, "type", "Task", "Issue type")
	cmd.Flags().StringVar(&opts.Assignee, "assignee", "", "Assignee username or account ID")
	cmd.Flags().StringVar(&opts.Description, "description", "", "Issue description (Markdown)")
	cmd.Flags().StringVar(&opts.Components, "components", "", "Comma-separated component names")
	cmd.Flags().StringVar(&opts.FieldsJSON, "fields-json", "", "Additional fields as JSON object")
	_ = cmd.MarkFlagRequired("project")
	_ = cmd.MarkFlagRequired("summary")

	return cmd
}

func createIssue(client *api.Client, opts CreateIssueOpts) (map[string]any, error) {
	fields := map[string]any{
		"project":   map[string]any{"key": opts.Project},
		"summary":   opts.Summary,
		"issuetype": map[string]any{"name": opts.IssueType},
	}

	if opts.Assignee != "" {
		if client.IsCloud() {
			fields["assignee"] = map[string]any{"accountId": opts.Assignee}
		} else {
			fields["assignee"] = map[string]any{"name": opts.Assignee}
		}
	}

	if opts.Description != "" {
		fields["description"] = opts.Description
	}

	if opts.Components != "" {
		components := make([]map[string]any, 0)
		for _, name := range strings.Split(opts.Components, ",") {
			trimmed := strings.TrimSpace(name)
			if trimmed == "" {
				continue
			}
			components = append(components, map[string]any{"name": trimmed})
		}
		if len(components) > 0 {
			fields["components"] = components
		}
	}

	if opts.FieldsJSON != "" {
		var extra map[string]any
		if err := json.Unmarshal([]byte(opts.FieldsJSON), &extra); err != nil {
			return nil, fmt.Errorf("parsing --fields-json: %w", err)
		}
		for k, v := range extra {
			fields[k] = v
		}
	}

	body := map[string]any{"fields": fields}
	path := client.JiraAPIPath("issue")

	var result map[string]any
	_, err := client.Post(path, body, &result)
	if err != nil {
		return nil, fmt.Errorf("creating issue: %w", err)
	}

	return result, nil
}
