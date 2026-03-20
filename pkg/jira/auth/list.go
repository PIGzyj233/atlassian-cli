package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all configured instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")
			if output == "json" {
				type hostEntry struct {
					Host    string `json:"host"`
					Type    string `json:"type"`
					Auth    string `json:"auth"`
					Default bool   `json:"default"`
				}
				var entries []hostEntry
				for name, host := range f.Config.Hosts {
					entries = append(entries, hostEntry{
						Host:    name,
						Type:    host.Type,
						Auth:    host.Auth,
						Default: name == f.Config.Defaults.Host,
					})
				}
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(entries)
			}

			for name, host := range f.Config.Hosts {
				marker := "  "
				if name == f.Config.Defaults.Host {
					marker = "* "
				}
				fmt.Fprintf(os.Stdout, "%s%s (%s, %s)\n", marker, name, host.Type, host.Auth)
			}
			return nil
		},
	}
	return cmd
}
