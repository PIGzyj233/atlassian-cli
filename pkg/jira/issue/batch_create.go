package issue

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdBatchCreate(f *cmdutil.Factory) *cobra.Command {
	var (
		file         string
		validateOnly bool
	)

	cmd := &cobra.Command{
		Use:   "batch-create",
		Short: "Create multiple Jira issues from a JSON file",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("reading file %s: %w", file, err)
			}

			var issues []map[string]any
			if err := json.Unmarshal(data, &issues); err != nil {
				return fmt.Errorf("parsing JSON: %w", err)
			}

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			if validateOnly {
				for i, issueData := range issues {
					if _, err := normalizeBatchIssueFields(client, issueData); err != nil {
						return fmt.Errorf("validating issue %d: %w", i, err)
					}
				}
				return output.Print(cmd, map[string]any{
					"valid": true,
					"count": len(issues),
				})
			}

			result, err := batchCreateIssues(client, issues)
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "JSON file with issue data (required)")
	cmd.Flags().BoolVar(&validateOnly, "validate-only", false, "Validate without creating")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}

func batchCreateIssues(client *api.Client, issues []map[string]any) (map[string]any, error) {
	issueUpdates := make([]map[string]any, 0, len(issues))
	for i, issueData := range issues {
		fields, err := normalizeBatchIssueFields(client, issueData)
		if err != nil {
			return nil, fmt.Errorf("preparing issue %d: %w", i, err)
		}
		issueUpdates = append(issueUpdates, map[string]any{"fields": fields})
	}

	body := map[string]any{"issueUpdates": issueUpdates}
	path := client.JiraAPIPath("issue/bulk")

	var result map[string]any
	_, err := client.Post(path, body, &result)
	if err != nil {
		return nil, fmt.Errorf("batch creating issues: %w", err)
	}
	return result, nil
}

func normalizeBatchIssueFields(client *api.Client, issueData map[string]any) (map[string]any, error) {
	if isFieldPayload(issueData) {
		cloned := make(map[string]any, len(issueData))
		for k, v := range issueData {
			cloned[k] = v
		}
		return cloned, nil
	}

	projectKey, _ := issueData["project_key"].(string)
	summary, _ := issueData["summary"].(string)
	issueType, _ := issueData["issue_type"].(string)
	if strings.TrimSpace(projectKey) == "" || strings.TrimSpace(summary) == "" || strings.TrimSpace(issueType) == "" {
		return nil, fmt.Errorf("project_key, summary, and issue_type are required")
	}

	fields := map[string]any{
		"project":   map[string]any{"key": projectKey},
		"summary":   summary,
		"issuetype": map[string]any{"name": issueType},
	}

	if assignee, _ := issueData["assignee"].(string); strings.TrimSpace(assignee) != "" {
		if client != nil && client.IsCloud() {
			fields["assignee"] = map[string]any{"accountId": assignee}
		} else {
			fields["assignee"] = map[string]any{"name": assignee}
		}
	}
	if description, _ := issueData["description"].(string); description != "" {
		fields["description"] = description
	}
	if components, ok := issueData["components"]; ok {
		normalized, err := normalizeComponentsValue(components)
		if err != nil {
			return nil, err
		}
		if len(normalized) > 0 {
			fields["components"] = normalized
		}
	}

	for k, v := range issueData {
		switch k {
		case "project_key", "summary", "issue_type", "assignee", "description", "components":
			continue
		default:
			fields[k] = v
		}
	}

	return fields, nil
}

func normalizeComponentsValue(value any) ([]map[string]any, error) {
	switch v := value.(type) {
	case string:
		return parseComponentNames(v), nil
	case []string:
		components := make([]map[string]any, 0, len(v))
		for _, name := range v {
			components = append(components, map[string]any{"name": name})
		}
		return components, nil
	case []any:
		components := make([]map[string]any, 0, len(v))
		for _, item := range v {
			name, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("components must contain only strings")
			}
			components = append(components, map[string]any{"name": name})
		}
		return components, nil
	default:
		return nil, fmt.Errorf("unsupported components type %T", value)
	}
}

func isFieldPayload(issueData map[string]any) bool {
	_, hasProject := issueData["project"]
	_, hasIssueType := issueData["issuetype"]
	_, hasSummary := issueData["summary"]
	return hasProject && hasIssueType && hasSummary
}
