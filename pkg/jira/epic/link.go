package epic

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdLink creates the `epic link` command.
func NewCmdLink(f *cmdutil.Factory) *cobra.Command {
	var (
		epicKey string
		field   string
	)

	cmd := &cobra.Command{
		Use:   "link <issue-key>",
		Short: "Link an issue to an epic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			// Determine the field to use based on instance type if not overridden.
			linkField := field
			if linkField == "" {
				if client.IsCloud() {
					linkField = "parent"
				} else {
					linkField = "customfield_10014"
				}
			}

			if err := linkEpic(client, args[0], epicKey, linkField); err != nil {
				return err
			}

			return output.Print(cmd, map[string]any{
				"success": true,
				"key":     args[0],
				"epic":    epicKey,
				"field":   linkField,
				"message": fmt.Sprintf("Issue %s linked to epic %s", args[0], epicKey),
			})
		},
	}

	cmd.Flags().StringVar(&epicKey, "epic", "", "Epic issue key to link to (required)")
	cmd.Flags().StringVar(&field, "field", "", "Field name for epic link (default: \"parent\" for Cloud, \"customfield_10014\" for Server/DC)")
	_ = cmd.MarkFlagRequired("epic")

	return cmd
}

// linkEpic sets the epic link field on the given issue.
func linkEpic(client *api.Client, issueKey, epicKey, field string) error {
	var fieldValue any
	if field == "parent" {
		fieldValue = map[string]any{"key": epicKey}
	} else {
		fieldValue = epicKey
	}

	body := map[string]any{
		"fields": map[string]any{
			field: fieldValue,
		},
	}

	path := client.JiraAPIPath("issue/" + issueKey)
	_, err := client.Put(path, body, nil)
	if err != nil {
		return fmt.Errorf("linking issue %s to epic %s: %w", issueKey, epicKey, err)
	}
	return nil
}
