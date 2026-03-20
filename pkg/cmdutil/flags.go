package cmdutil

import "github.com/spf13/cobra"

// GetHost returns the --host flag value from the command hierarchy.
func GetHost(cmd *cobra.Command) string {
	flag := cmd.Flag("host")
	if flag == nil {
		return ""
	}
	return flag.Value.String()
}

// GetOutput returns the --output flag value.
func GetOutput(cmd *cobra.Command) string {
	flag := cmd.Flag("output")
	if flag == nil || flag.Value.String() == "" {
		return "json"
	}
	return flag.Value.String()
}
