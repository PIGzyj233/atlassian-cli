package comment

import (
	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdList creates the `comment list` command.
func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list <page-id>",
		Short: "List comments on a Confluence page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			path := client.ConfluenceAPIPath("content/" + args[0] + "/child/comment")
			rb := api.NewRequestBuilder(path).
				Query("expand", "body.view.value,version").
				Query("depth", "all")

			var result map[string]any
			_, err = client.Get(rb.BuildPath(), &result)
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}
}
