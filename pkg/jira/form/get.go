package form

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdGet creates the `form get` command.
func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <issue-key> <form-id>",
		Short: "Get a specific ProForma form on an issue",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := getForm(client, args[0], args[1])
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}
	return cmd
}

func getForm(client *api.Client, issueKey, formID string) (map[string]any, error) {
	path := client.JiraAPIPath("issue/" + issueKey + "/properties/proforma.forms")

	var property map[string]any
	_, err := client.Get(path, &property)
	if err != nil {
		return nil, fmt.Errorf("getting forms for %s: %w", issueKey, err)
	}

	// The property value contains the forms data.
	value, ok := property["value"]
	if !ok {
		return nil, fmt.Errorf("no form data found for issue %s", issueKey)
	}

	// The value may be a list of forms or a map containing forms.
	switch v := value.(type) {
	case []any:
		for _, item := range v {
			formMap, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if fmt.Sprintf("%v", formMap["id"]) == formID {
				return formMap, nil
			}
		}
	case map[string]any:
		// If value is a map, check for a "forms" key containing a list.
		forms, ok := v["forms"]
		if !ok {
			// Maybe the value itself is keyed by form ID.
			if formData, ok := v[formID]; ok {
				if fd, ok := formData.(map[string]any); ok {
					return fd, nil
				}
			}
			return nil, fmt.Errorf("form %s not found on issue %s", formID, issueKey)
		}
		formsList, ok := forms.([]any)
		if !ok {
			return nil, fmt.Errorf("unexpected forms data format for issue %s", issueKey)
		}
		for _, item := range formsList {
			formMap, ok := item.(map[string]any)
			if !ok {
				continue
			}
			if fmt.Sprintf("%v", formMap["id"]) == formID {
				return formMap, nil
			}
		}
	}

	return nil, fmt.Errorf("form %s not found on issue %s", formID, issueKey)
}
