package dev

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdDev creates the `dev` subcommand group.
func NewCmdDev(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev <command>",
		Short: "View development information for issues",
	}
	cmd.AddCommand(NewCmdInfo(f))
	cmd.AddCommand(NewCmdBatchInfo(f))
	return cmd
}
