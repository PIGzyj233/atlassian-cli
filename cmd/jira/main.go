package main

import (
	"fmt"
	"os"

	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	jiracmd "github.com/PigZyj2333/atlassian-cli/pkg/jira"
)

var version = "dev"

func main() {
	f := cmdutil.NewFactory(version)
	rootCmd := jiracmd.NewCmdRoot(f)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
