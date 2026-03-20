package link

import (
	"fmt"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdTypes creates the `link types` command.
func NewCmdTypes(f *cmdutil.Factory) *cobra.Command {
	var filter string

	cmd := &cobra.Command{
		Use:   "types",
		Short: "List issue link types",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := listLinkTypes(client, filter)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&filter, "filter", "", "Filter link types by name (case-insensitive contains match)")

	return cmd
}

func listLinkTypes(client *api.Client, filter string) (map[string]any, error) {
	path := client.JiraAPIPath("issueLinkType")

	var result map[string]any
	_, err := client.Get(path, &result)
	if err != nil {
		return nil, fmt.Errorf("listing issue link types: %w", err)
	}

	if filter != "" {
		result = filterLinkTypes(result, filter)
	}

	return result, nil
}

func filterLinkTypes(result map[string]any, filter string) map[string]any {
	linkTypes, ok := result["issueLinkTypes"].([]any)
	if !ok {
		return result
	}

	filterLower := strings.ToLower(filter)
	filtered := make([]any, 0)

	for _, lt := range linkTypes {
		ltMap, ok := lt.(map[string]any)
		if !ok {
			continue
		}
		name, _ := ltMap["name"].(string)
		if strings.Contains(strings.ToLower(name), filterLower) {
			filtered = append(filtered, lt)
		}
	}

	return map[string]any{
		"issueLinkTypes": filtered,
	}
}
