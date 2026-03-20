package transition

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdList creates the `transition list` command.
// Ref: Python TransitionsMixin.get_transitions() — returns raw transitions array.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <issue-key>",
		Short: "List available transitions for an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := listTransitions(client, args[0])
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}
	return cmd
}

func listTransitions(client *api.Client, issueKey string) ([]any, error) {
	path := client.JiraAPIPath("issue/" + issueKey + "/transitions")

	var result map[string]any
	_, err := client.Get(path, &result)
	if err != nil {
		return nil, fmt.Errorf("getting transitions for %s: %w", issueKey, err)
	}

	// Return the transitions array directly
	// Ref: Python get_transitions() returns response["transitions"]
	transitions, _ := result["transitions"].([]any)
	return transitions, nil
}
