package sla

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdSLA creates the `sla` subcommand group.
func NewCmdSLA(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sla <command>",
		Short: "View SLA information",
	}
	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdDates(f))
	return cmd
}
