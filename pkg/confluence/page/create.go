package page

import (
	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/convert"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdCreate creates the `page create` command.
func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var (
		space    string
		title    string
		content  string
		parentID string
		format   string
		emoji    string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new Confluence page",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			// Convert markdown to storage format if needed
			body := content
			if format == "markdown" {
				body = convert.MarkdownToStorage(content)
			}

			payload := map[string]any{
				"type":  "page",
				"title": title,
				"space": map[string]any{"key": space},
				"body": map[string]any{
					"storage": map[string]any{
						"value":          body,
						"representation": "storage",
					},
				},
			}

			if parentID != "" {
				payload["ancestors"] = []map[string]any{{"id": parentID}}
			}

			path := client.ConfluenceAPIPath("content")
			var result map[string]any
			_, err = client.Post(path, payload, &result)
			if err != nil {
				return err
			}

			// Set emoji via content property if provided
			if emoji != "" {
				if id, ok := result["id"].(string); ok {
					setPageEmoji(client, id, emoji)
				}
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&space, "space", "", "Space key (required)")
	cmd.Flags().StringVar(&title, "title", "", "Page title (required)")
	cmd.Flags().StringVar(&content, "content", "", "Page content (required)")
	cmd.Flags().StringVar(&parentID, "parent-id", "", "Parent page ID")
	cmd.Flags().StringVar(&format, "format", "markdown", "Content format: markdown or storage")
	cmd.Flags().StringVar(&emoji, "emoji", "", "Page title emoji")
	cmd.MarkFlagRequired("space")
	cmd.MarkFlagRequired("title")
	cmd.MarkFlagRequired("content")

	return cmd
}

// setPageEmoji sets the emoji-title-published content property on a page.
func setPageEmoji(client *api.Client, pageID, emoji string) {
	path := client.ConfluenceAPIPath("content/" + pageID + "/property/emoji-title-published")
	payload := map[string]any{
		"key":   "emoji-title-published",
		"value": emoji,
	}
	client.Post(path, payload, nil)
}

