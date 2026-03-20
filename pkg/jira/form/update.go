package form

import (
	"encoding/json"
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdUpdate creates the `form update` command.
func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var answersJSON string

	cmd := &cobra.Command{
		Use:   "update <issue-key> <form-id>",
		Short: "Update answers on a ProForma form",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := updateForm(client, args[0], args[1], answersJSON)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&answersJSON, "answers-json", "", "JSON string with answer data (required)")
	_ = cmd.MarkFlagRequired("answers-json")
	return cmd
}

func updateForm(client *api.Client, issueKey, formID, answersJSON string) (map[string]any, error) {
	// Parse the answers JSON.
	var answers map[string]any
	if err := json.Unmarshal([]byte(answersJSON), &answers); err != nil {
		return nil, fmt.Errorf("parsing --answers-json: %w", err)
	}

	// GET the current property value.
	path := client.JiraAPIPath("issue/" + issueKey + "/properties/proforma.forms")

	var property map[string]any
	_, err := client.Get(path, &property)
	if err != nil {
		return nil, fmt.Errorf("getting forms for %s: %w", issueKey, err)
	}

	value, ok := property["value"]
	if !ok {
		return nil, fmt.Errorf("no form data found for issue %s", issueKey)
	}

	// Find the form by ID and merge answers into it.
	updated, err := mergeAnswers(value, formID, answers)
	if err != nil {
		return nil, fmt.Errorf("form %s not found on issue %s", formID, issueKey)
	}

	// PUT the updated property value back.
	_, err = client.Put(path, updated, nil)
	if err != nil {
		return nil, fmt.Errorf("updating form %s on %s: %w", formID, issueKey, err)
	}

	return map[string]any{
		"success": true,
		"message": fmt.Sprintf("Form %s updated on issue %s", formID, issueKey),
	}, nil
}

// mergeAnswers finds the form by ID within the property value and merges the
// provided answers into it. It returns the updated property value to PUT back.
func mergeAnswers(value any, formID string, answers map[string]any) (any, error) {
	switch v := value.(type) {
	case []any:
		found := false
		for i, item := range v {
			formMap, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if fmt.Sprintf("%v", formMap["id"]) == formID {
				// Merge answers into existing answers map or create one.
				existing, _ := formMap["answers"].(map[string]any)
				if existing == nil {
					existing = make(map[string]any)
				}
				for k, val := range answers {
					existing[k] = val
				}
				formMap["answers"] = existing
				v[i] = formMap
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("form not found")
		}
		return v, nil

	case map[string]any:
		// Check for a "forms" key containing a list.
		if forms, ok := v["forms"]; ok {
			formsList, ok := forms.([]any)
			if !ok {
				return nil, fmt.Errorf("form not found")
			}
			found := false
			for i, item := range formsList {
				formMap, ok := item.(map[string]any)
				if !ok {
					continue
				}
				if fmt.Sprintf("%v", formMap["id"]) == formID {
					existing, _ := formMap["answers"].(map[string]any)
					if existing == nil {
						existing = make(map[string]any)
					}
					for k, val := range answers {
						existing[k] = val
					}
					formMap["answers"] = existing
					formsList[i] = formMap
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("form not found")
			}
			v["forms"] = formsList
			return v, nil
		}

		// Maybe keyed by form ID directly.
		if formData, ok := v[formID]; ok {
			fd, ok := formData.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("form not found")
			}
			existing, _ := fd["answers"].(map[string]any)
			if existing == nil {
				existing = make(map[string]any)
			}
			for k, val := range answers {
				existing[k] = val
			}
			fd["answers"] = existing
			v[formID] = fd
			return v, nil
		}

		return nil, fmt.Errorf("form not found")
	}

	return nil, fmt.Errorf("form not found")
}
