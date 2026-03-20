package label

import (
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdList creates the `label list` command.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list <page-id>",
		Short: "List labels on a page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			path := client.ConfluenceAPIPath("content/" + args[0] + "/label")
			var result map[string]any
			_, err = client.Get(path, &result)
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}
}
