package page

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdHistory creates the `page history` command.
func NewCmdHistory(f *cmdutil.Factory) *cobra.Command {
	var (
		version int
		raw     bool
	)

	cmd := &cobra.Command{
		Use:   "history <page-id>",
		Short: "Get page history or a specific version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pageID := args[0]
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			if version > 0 {
				// Fetch specific historical version
				path := client.ConfluenceAPIPath("content/" + pageID)
				rb := api.NewRequestBuilder(path).
					Query("expand", "body.storage,version,space").
					Query("status", "historical").
					QueryInt("version", version)

				var result map[string]any
				_, err := client.Get(rb.BuildPath(), &result)
				if err != nil {
					return fmt.Errorf("getting page version %d: %w", version, err)
				}

				if !raw {
					convertStorageBody(result)
				}

				return output.Print(cmd, result)
			}

			// Fetch page history metadata
			path := client.ConfluenceAPIPath("content/" + pageID + "/history")
			var result map[string]any
			_, err = client.Get(path, &result)
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().IntVar(&version, "version", 0, "Specific version number to retrieve")
	cmd.Flags().BoolVar(&raw, "raw", false, "Return raw storage format")

	return cmd
}
