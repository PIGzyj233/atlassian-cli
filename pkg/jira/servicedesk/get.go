package servicedesk

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdGet creates the `servicedesk get` command.
func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <project-key>",
		Short: "Get a service desk by project key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := getServiceDesk(client, args[0])
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	return cmd
}

func getServiceDesk(client *api.Client, projectKey string) (map[string]any, error) {
	path := fmt.Sprintf("/rest/servicedeskapi/servicedesk/projectKey:%s", projectKey)

	// ServiceDesk API requires X-ExperimentalApi header.
	headers := map[string]string{"X-ExperimentalApi": "opt-in"}

	var result map[string]any
	_, err := client.GetWithHeaders(path, &result, headers)
	if err != nil {
		return nil, fmt.Errorf("getting service desk for project %s: %w", projectKey, err)
	}
	return result, nil
}
