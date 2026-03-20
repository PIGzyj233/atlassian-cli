package attachment

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdDelete creates the `attachment delete` command.
func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	var confirm bool

	cmd := &cobra.Command{
		Use:   "delete <attachment-id>",
		Short: "Delete a Confluence attachment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm {
				return fmt.Errorf("use --confirm to delete attachment %s", args[0])
			}

			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			path := client.ConfluenceAPIPath("content/" + args[0])
			_, err = client.Delete(path, nil)
			if err != nil {
				return err
			}

			return output.Print(cmd, map[string]any{
				"deleted":       true,
				"attachment_id": args[0],
			})
		},
	}

	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm deletion")

	return cmd
}
