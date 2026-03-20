package page

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/convert"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdUpdate creates the `page update` command.
func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var (
		title     string
		content   string
		minorEdit bool
		comment   string
		format    string
		emoji     string
	)

	cmd := &cobra.Command{
		Use:   "update <page-id>",
		Short: "Update an existing Confluence page",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pageID := args[0]
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			// Fetch current version
			currentVersion, currentTitle, err := fetchCurrentVersion(client, pageID)
			if err != nil {
				return fmt.Errorf("fetching current version: %w", err)
			}

			if title == "" {
				title = currentTitle
			}

			// Convert markdown to storage format if needed
			body := content
			if format == "markdown" {
				body = convert.MarkdownToStorage(content)
			}

			payload := map[string]any{
				"type":  "page",
				"title": title,
				"body": map[string]any{
					"storage": map[string]any{
						"value":          body,
						"representation": "storage",
					},
				},
				"version": map[string]any{
					"number":    currentVersion + 1,
					"minorEdit": minorEdit,
					"message":   comment,
				},
			}

			path := client.ConfluenceAPIPath("content/" + pageID)
			var result map[string]any
			_, err = client.Put(path, payload, &result)
			if err != nil {
				return err
			}

			// Set or remove emoji via content property
			if cmd.Flags().Changed("emoji") {
				if emoji != "" {
					if err := setPageEmoji(client, pageID, emoji); err != nil {
						return fmt.Errorf("setting emoji: %w", err)
					}
				} else {
					// Empty string means remove emoji
					deletePageEmoji(client, pageID)
				}
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "New page title")
	cmd.Flags().StringVar(&content, "content", "", "New page content (required)")
	cmd.Flags().BoolVar(&minorEdit, "minor-edit", false, "Mark as minor edit")
	cmd.Flags().StringVar(&comment, "comment", "", "Version comment")
	cmd.Flags().StringVar(&format, "format", "markdown", "Content format: markdown or storage")
	cmd.Flags().StringVar(&emoji, "emoji", "", "Page title emoji (empty to remove)")
	cmd.MarkFlagRequired("content")

	return cmd
}

// fetchCurrentVersion fetches the current version number and title of a page.
func fetchCurrentVersion(client *api.Client, pageID string) (int, string, error) {
	path := client.ConfluenceAPIPath("content/" + pageID)
	rb := api.NewRequestBuilder(path).Query("expand", "version")

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return 0, "", err
	}

	version, ok := result["version"].(map[string]any)
	if !ok {
		return 0, "", fmt.Errorf("missing version info for page %s", pageID)
	}

	number, ok := version["number"].(float64)
	if !ok {
		return 0, "", fmt.Errorf("missing version number for page %s", pageID)
	}

	title, _ := result["title"].(string)
	return int(number), title, nil
}
