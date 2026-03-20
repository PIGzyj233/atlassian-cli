package page

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
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
			if targetID == "" && targetSpace != "" {
				// Cross-space move to root: look up the space homepage
				homepageID, err := fetchSpaceHomepage(client, targetSpace)
				if err != nil {
					return fmt.Errorf("looking up space %q: %w", targetSpace, err)
				}
				targetID = homepageID
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
	cmd.Flags().StringVar(&targetSpace, "target-space", "", "Target space key (moves to space root if no parent specified)")
	cmd.Flags().StringVar(&position, "position", "append", "Position: append, above, below")

	return cmd
}

// fetchSpaceHomepage looks up the space homepage ID for cross-space moves.
func fetchSpaceHomepage(client *api.Client, spaceKey string) (string, error) {
	path := client.ConfluenceAPIPath("space/" + spaceKey)
	rb := api.NewRequestBuilder(path).Query("expand", "homepage")

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return "", err
	}

	hp, ok := result["homepage"].(map[string]any)
	if !ok {
		return "", fmt.Errorf("space %s has no homepage", spaceKey)
	}

	id, ok := hp["id"].(string)
	if !ok {
		return "", fmt.Errorf("space %s homepage has no ID", spaceKey)
	}

	return id, nil
}
