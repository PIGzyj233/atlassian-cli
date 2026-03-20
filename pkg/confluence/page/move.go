package page

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdMove creates the `page move` command.
func NewCmdMove(f *cmdutil.Factory) *cobra.Command {
	var (
		targetParentID string
		targetSpace    string
		position       string
	)

	cmd := &cobra.Command{
		Use:   "move <page-id>",
		Short: "Move a page to a new parent or space",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pageID := args[0]
			if targetParentID == "" && targetSpace == "" {
				return fmt.Errorf("provide --target-parent-id or --target-space")
			}

			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			targetID := targetParentID
			if targetID == "" {
				targetID = targetSpace
			}

			path := client.ConfluenceAPIPath(
				fmt.Sprintf("content/%s/move/%s/%s", pageID, position, targetID),
			)

			var result map[string]any
			_, err = client.Put(path, nil, &result)
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&targetParentID, "target-parent-id", "", "Target parent page ID")
	cmd.Flags().StringVar(&targetSpace, "target-space", "", "Target space key")
	cmd.Flags().StringVar(&position, "position", "append", "Position: append, above, below")

	return cmd
}
