package confluence

import (
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdRoot creates the root `confluence` command.
func NewCmdRoot(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "confluence <command> [flags]",
		Short:         "Work with Confluence from the command line",
		Long:          "A CLI tool for interacting with Atlassian Confluence, designed for both humans and LLM automation.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().String("host", "", "Target Atlassian instance hostname")
	cmd.PersistentFlags().String("output", "json", "Output format: json, text, table")
	cmd.PersistentFlags().Bool("verbose", false, "Enable debug logging to stderr")
	cmd.PersistentFlags().Bool("no-color", false, "Disable ANSI colors")

	cmd.Version = f.Version

	return cmd
}
