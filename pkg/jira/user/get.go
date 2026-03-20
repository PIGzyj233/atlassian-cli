package user

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdGet creates the `user get` command.
func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <identifier>",
		Short: "Get a Jira user by account ID (Cloud) or username (Server/DC)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := getUser(client, args[0])
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	return cmd
}

func getUser(client *api.Client, identifier string) (map[string]any, error) {
	path := client.JiraAPIPath("user")

	rb := api.NewRequestBuilder(path)
	if client.IsCloud() {
		rb.Query("accountId", identifier)
	} else {
		rb.Query("username", identifier)
	}

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("getting user %s: %w", identifier, err)
	}

	return result, nil
}
