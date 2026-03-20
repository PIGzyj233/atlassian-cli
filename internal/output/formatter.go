package output

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

// Formatter writes structured data to an output stream.
type Formatter interface {
	Write(w io.Writer, data any) error
}

// NewFormatter creates a Formatter based on the format name.
func NewFormatter(format string) Formatter {
	switch format {
	case "table":
		return &TableFormatter{}
	case "text":
		return &TextFormatter{}
	default:
		return &JSONFormatter{}
	}
}

// Print writes data to stdout using the format from --output flag.
func Print(cmd *cobra.Command, data any) error {
	format, _ := cmd.Flags().GetString("output")
	if format == "" {
		format = "json"
	}
	f := NewFormatter(format)
	return f.Write(os.Stdout, data)
}

// PrintError writes an error as JSON to stderr.
func PrintError(cmd *cobra.Command, err error) {
	format, _ := cmd.Flags().GetString("output")
	if format == "json" {
		fmt.Fprintf(os.Stderr, `{"error": %q}`+"\n", err.Error())
	} else {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	}
}
