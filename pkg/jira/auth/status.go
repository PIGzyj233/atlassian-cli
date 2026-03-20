package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdStatus(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current authentication state",
		RunE: func(cmd *cobra.Command, args []string) error {
			host, name, err := f.Config.DefaultHost()
			if err != nil {
				return fmt.Errorf("not authenticated: %w", err)
			}

			output, _ := cmd.Flags().GetString("output")
			if output == "json" {
				data := map[string]any{
					"host":     name,
					"type":     host.Type,
					"auth":     host.Auth,
					"username": host.Username,
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(data)
			}

			fmt.Fprintf(os.Stdout, "Host:     %s\n", name)
			fmt.Fprintf(os.Stdout, "Type:     %s\n", host.Type)
			fmt.Fprintf(os.Stdout, "Auth:     %s\n", host.Auth)
			if host.Username != "" {
				fmt.Fprintf(os.Stdout, "Username: %s\n", host.Username)
			}
			return nil
		},
	}
	return cmd
}
