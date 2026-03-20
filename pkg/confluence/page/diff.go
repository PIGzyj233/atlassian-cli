package page

import (
	"fmt"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/convert"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"
)

// NewCmdDiff creates the `page diff` command.
func NewCmdDiff(f *cmdutil.Factory) *cobra.Command {
	var (
		fromVersion int
		toVersion   int
	)

	cmd := &cobra.Command{
		Use:   "diff <page-id>",
		Short: "Show diff between two page versions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pageID := args[0]
			if fromVersion <= 0 || toVersion <= 0 {
				return fmt.Errorf("both --from and --to version numbers are required")
			}

			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			// Fetch both versions
			fromContent, err := fetchVersionContent(client, pageID, fromVersion)
			if err != nil {
				return fmt.Errorf("fetching version %d: %w", fromVersion, err)
			}

			toContent, err := fetchVersionContent(client, pageID, toVersion)
			if err != nil {
				return fmt.Errorf("fetching version %d: %w", toVersion, err)
			}

			// Convert to markdown
			fromMD := convert.StorageToMarkdown(fromContent)
			toMD := convert.StorageToMarkdown(toContent)

			// Compute proper unified diff using Myers algorithm
			// (matches Python reference: difflib.unified_diff at pages.py:1187)
			diff := unifiedDiff(
				strings.Split(fromMD, "\n"),
				strings.Split(toMD, "\n"),
				fmt.Sprintf("v%d", fromVersion),
				fmt.Sprintf("v%d", toVersion),
			)

			result := map[string]any{
				"page_id":      pageID,
				"from_version": fromVersion,
				"to_version":   toVersion,
				"diff":         diff,
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().IntVar(&fromVersion, "from", 0, "Source version number (required)")
	cmd.Flags().IntVar(&toVersion, "to", 0, "Target version number (required)")

	return cmd
}

func fetchVersionContent(client *api.Client, pageID string, version int) (string, error) {
	path := client.ConfluenceAPIPath("content/" + pageID)
	rb := api.NewRequestBuilder(path).
		Query("expand", "body.storage").
		Query("status", "historical").
		QueryInt("version", version)

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return "", err
	}

	body, _ := result["body"].(map[string]any)
	storage, _ := body["storage"].(map[string]any)
	value, _ := storage["value"].(string)
	return value, nil
}

// unifiedDiff produces a proper unified diff with hunk headers using Myers algorithm.
func unifiedDiff(from, to []string, fromLabel, toLabel string) string {
	// Append newline to each line for difflib compatibility
	fromLines := make([]string, len(from))
	toLines := make([]string, len(to))
	for i, l := range from {
		fromLines[i] = l + "\n"
	}
	for i, l := range to {
		toLines[i] = l + "\n"
	}

	diff := difflib.UnifiedDiff{
		A:        fromLines,
		B:        toLines,
		FromFile: fromLabel,
		ToFile:   toLabel,
		Context:  3,
	}
	text, _ := difflib.GetUnifiedDiffString(diff)
	return text
}
