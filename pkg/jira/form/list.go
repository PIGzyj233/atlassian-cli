package form

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdList creates the `form list` command.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list <issue-key>",
		Short: "List ProForma forms on an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := listForms(client, args[0])
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}
	return cmd
}

func listForms(client *api.Client, issueKey string) (map[string]any, error) {
	path := client.JiraAPIPath("issue/" + issueKey + "/properties/proforma.forms")

	var result map[string]any
	_, err := client.Get(path, &result)
	if err != nil {
		return nil, fmt.Errorf("listing forms for %s: %w", issueKey, err)
	}
	return result, nil
}
