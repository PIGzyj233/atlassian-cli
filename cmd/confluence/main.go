package main

import (
	"fmt"
	"os"

	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	confluencecmd "github.com/PigZyj2333/atlassian-cli/pkg/confluence"
)

var version = "dev"

func main() {
	f := cmdutil.NewFactory(version)
	rootCmd := confluencecmd.NewCmdRoot(f)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
