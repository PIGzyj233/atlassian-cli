package label

import (
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAdd creates the `label add` command.
func NewCmdAdd(f *cmdutil.Factory) *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "add <page-id>",
		Short: "Add a label to a page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			payload := []map[string]any{
				{
					"prefix": "global",
					"name":   name,
				},
			}

			path := client.ConfluenceAPIPath("content/" + args[0] + "/label")
			var result map[string]any
			_, err = client.Post(path, payload, &result)
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Label name (required)")
	cmd.MarkFlagRequired("name")

	return cmd
}
