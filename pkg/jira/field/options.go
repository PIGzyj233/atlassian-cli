package field

import (
	"fmt"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdOptions creates the `field options` command.
// Ref: Python FieldOptionsMixin — Cloud uses context option API, Server uses createmeta.
func NewCmdOptions(f *cmdutil.Factory) *cobra.Command {
	var (
		contextID string
		project   string
		issueType string
		contains  string
	)

	cmd := &cobra.Command{
		Use:   "options <field-id>",
		Short: "Get allowed option values for a custom field",
		Long: `Get allowed option values for a custom field.

Cloud: Uses the Field Context Option API. If --context-id is not provided,
auto-resolves by fetching contexts and using the global context.

Server/DC: Uses createmeta to get allowedValues. Requires --project and
--issue-type parameters.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			var result any
			if client.IsCloud() {
				result, err = getFieldOptionsCloud(client, args[0], contextID, contains)
			} else {
				if project == "" || issueType == "" {
					return fmt.Errorf("Server/DC requires --project and --issue-type to retrieve field options")
				}
				result, err = getFieldOptionsServer(client, args[0], project, issueType, contains)
			}
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&contextID, "context-id", "", "Context ID (Cloud only, auto-resolved if not provided)")
	cmd.Flags().StringVar(&project, "project", "", "Project key (required for Server/DC)")
	cmd.Flags().StringVar(&issueType, "issue-type", "", "Issue type name (required for Server/DC)")
	cmd.Flags().StringVar(&contains, "contains", "", "Filter options containing this string")
	return cmd
}

// getFieldOptionsCloud fetches field options via Cloud API.
// Ref: Python _get_field_options_cloud() — auto-resolves context, paginates.
func getFieldOptionsCloud(client *api.Client, fieldID, contextID, contains string) (any, error) {
	// Auto-resolve context if not provided
	if contextID == "" {
		resolvedID, err := resolveCloudContext(client, fieldID)
		if err != nil {
			return nil, err
		}
		contextID = resolvedID
	}

	// Paginate through all options
	var allOptions []any
	startAt := 0
	maxResults := 100

	for {
		rb := api.NewRequestBuilder(fmt.Sprintf("/rest/api/3/field/%s/context/%s/option", fieldID, contextID)).
			QueryInt("startAt", startAt).
			QueryInt("maxResults", maxResults)

		var result map[string]any
		_, err := client.Get(rb.BuildPath(), &result)
		if err != nil {
			return nil, fmt.Errorf("fetching field options: %w", err)
		}

		values, _ := result["values"].([]any)
		allOptions = append(allOptions, values...)

		total := int(getFloat(result, "total", float64(len(values))))
		startAt += len(values)

		if startAt >= total || len(values) == 0 {
			break
		}
	}

	// Filter by --contains if specified
	if contains != "" {
		allOptions = filterOptions(allOptions, contains)
	}

	return allOptions, nil
}

// resolveCloudContext auto-resolves the global context for a field.
// Ref: Python _get_field_options_cloud() — prefers global context.
func resolveCloudContext(client *api.Client, fieldID string) (string, error) {
	rb := api.NewRequestBuilder(fmt.Sprintf("/rest/api/3/field/%s/context", fieldID)).
		QueryInt("maxResults", 100)

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return "", fmt.Errorf("fetching field contexts for %s: %w", fieldID, err)
	}

	values, _ := result["values"].([]any)
	if len(values) == 0 {
		return "", fmt.Errorf("no contexts found for field %s; cannot retrieve options without a context", fieldID)
	}

	// Prefer global context (isGlobalContext == true)
	for _, v := range values {
		ctx, ok := v.(map[string]any)
		if !ok {
			continue
		}
		if isGlobal, _ := ctx["isGlobalContext"].(bool); isGlobal {
			if id, ok := ctx["id"].(string); ok {
				return id, nil
			}
			// id might be a float from JSON
			if idNum, ok := ctx["id"].(float64); ok {
				return fmt.Sprintf("%.0f", idNum), nil
			}
		}
	}

	// Fallback to first context
	firstCtx, _ := values[0].(map[string]any)
	if id, ok := firstCtx["id"].(string); ok {
		return id, nil
	}
	if idNum, ok := firstCtx["id"].(float64); ok {
		return fmt.Sprintf("%.0f", idNum), nil
	}

	return "", fmt.Errorf("could not resolve context ID for field %s", fieldID)
}

// getFieldOptionsServer fetches field options via Server/DC createmeta.
// Ref: Python _get_field_options_server() — resolves issue type ID, then extracts allowedValues.
func getFieldOptionsServer(client *api.Client, fieldID, projectKey, issueType, contains string) (any, error) {
	// Step 1: Get issue types for the project to resolve issueType → ID
	issueTypeID, err := resolveIssueTypeID(client, projectKey, issueType)
	if err != nil {
		return nil, err
	}

	// Step 2: Paginate through createmeta fields
	startAt := 0
	maxResults := 50

	for {
		rb := api.NewRequestBuilder(client.JiraAPIPath(
			fmt.Sprintf("issue/createmeta/%s/issuetypes/%s", projectKey, issueTypeID))).
			QueryInt("startAt", startAt).
			QueryInt("maxResults", maxResults)

		var meta map[string]any
		_, err := client.Get(rb.BuildPath(), &meta)
		if err != nil {
			return nil, fmt.Errorf("fetching createmeta: %w", err)
		}

		fieldEntries, _ := meta["values"].([]any)
		for _, entry := range fieldEntries {
			fm, ok := entry.(map[string]any)
			if !ok {
				continue
			}
			if fid, _ := fm["fieldId"].(string); fid == fieldID {
				options, _ := fm["allowedValues"].([]any)
				if contains != "" {
					options = filterOptions(options, contains)
				}
				return options, nil
			}
		}

		total := int(getFloat(meta, "total", float64(len(fieldEntries))))
		startAt += len(fieldEntries)
		if startAt >= total || len(fieldEntries) == 0 {
			break
		}
	}

	return []any{}, nil
}

// resolveIssueTypeID looks up an issue type name → ID for a project.
func resolveIssueTypeID(client *api.Client, projectKey, issueTypeName string) (string, error) {
	path := client.JiraAPIPath("project/" + projectKey)
	var project map[string]any
	_, err := client.Get(path, &project)
	if err != nil {
		return "", fmt.Errorf("fetching project %s: %w", projectKey, err)
	}

	issueTypes, _ := project["issueTypes"].([]any)
	lowerName := strings.ToLower(issueTypeName)
	for _, it := range issueTypes {
		itMap, ok := it.(map[string]any)
		if !ok {
			continue
		}
		name, _ := itMap["name"].(string)
		if strings.ToLower(name) == lowerName {
			if id, ok := itMap["id"].(string); ok {
				return id, nil
			}
			if idNum, ok := itMap["id"].(float64); ok {
				return fmt.Sprintf("%.0f", idNum), nil
			}
		}
	}

	return "", fmt.Errorf("issue type '%s' not found in project '%s'", issueTypeName, projectKey)
}

func filterOptions(options []any, contains string) []any {
	lower := strings.ToLower(contains)
	var filtered []any
	for _, opt := range options {
		optMap, ok := opt.(map[string]any)
		if !ok {
			continue
		}
		// Check "value" field
		if val, _ := optMap["value"].(string); strings.Contains(strings.ToLower(val), lower) {
			filtered = append(filtered, opt)
			continue
		}
		// Check "name" field (some options use name instead of value)
		if name, _ := optMap["name"].(string); strings.Contains(strings.ToLower(name), lower) {
			filtered = append(filtered, opt)
		}
	}
	return filtered
}

func getFloat(m map[string]any, key string, defaultVal float64) float64 {
	if v, ok := m[key].(float64); ok {
		return v
	}
	return defaultVal
}
