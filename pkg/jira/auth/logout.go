package auth

import (
	"fmt"
	"os"

	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdLogout(f *cmdutil.Factory) *cobra.Command {
	var host string
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Remove stored credentials for a host",
		RunE: func(cmd *cobra.Command, args []string) error {
			target := host
			if target == "" {
				target = f.Config.Defaults.Host
			}
			if target == "" {
				return fmt.Errorf("specify --host or configure a default host")
			}
			if _, ok := f.Config.Hosts[target]; !ok {
				return fmt.Errorf("host %q not found in config", target)
			}
			delete(f.Config.Hosts, target)
			if f.Config.Defaults.Host == target {
				f.Config.Defaults.Host = ""
			}
			if err := f.Config.Save(f.ConfigPath); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Logged out of %s\n", target)
			return nil
		},
	}
	cmd.Flags().StringVar(&host, "host", "", "Host to log out of")
	return cmd
}
