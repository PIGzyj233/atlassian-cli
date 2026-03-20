package auth

import (
	"fmt"
	"os"

	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSwitch(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch <host>",
		Short: "Change default instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			if _, ok := f.Config.Hosts[target]; !ok {
				return fmt.Errorf("host %q not found in config", target)
			}
			f.Config.Defaults.Host = target
			if err := f.Config.Save(f.ConfigPath); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "Switched default host to %s\n", target)
			return nil
		},
	}
	return cmd
}
