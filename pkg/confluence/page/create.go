package page

import (
	"fmt"

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
					if err := setPageEmoji(client, id, emoji); err != nil {
						return fmt.Errorf("setting emoji: %w", err)
					}
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

// setPageEmoji sets both published and draft emoji content properties.
// Matches Python reference: pages.py:257-296 (_set_page_emoji).
func setPageEmoji(client *api.Client, pageID, emoji string) error {
	for _, key := range []string{"emoji-title-published", "emoji-title-draft"} {
		if err := setContentProperty(client, pageID, key, emoji); err != nil {
			return fmt.Errorf("setting %s: %w", key, err)
		}
	}
	return nil
}

// deletePageEmoji removes both emoji content properties.
func deletePageEmoji(client *api.Client, pageID string) {
	for _, key := range []string{"emoji-title-published", "emoji-title-draft"} {
		path := client.ConfluenceAPIPath("content/" + pageID + "/property/" + key)
		client.Delete(path, nil) //nolint:errcheck // best-effort removal
	}
}

// setContentProperty creates or updates a content property with version tracking.
// Matches Python reference: pages.py:200-255 (_set_single_property).
func setContentProperty(client *api.Client, pageID, key, value string) error {
	propPath := client.ConfluenceAPIPath("content/" + pageID + "/property/" + key)

	// Try to GET existing property for version
	var existing map[string]any
	_, getErr := client.Get(propPath, &existing)

	payload := map[string]any{
		"key":   key,
		"value": value,
	}

	if getErr == nil && existing != nil {
		// Property exists — PUT with incremented version
		if ver, ok := existing["version"].(map[string]any); ok {
			if num, ok := ver["number"].(float64); ok {
				payload["version"] = map[string]any{"number": int(num) + 1}
			}
		}
		_, err := client.Put(propPath, payload, nil)
		return err
	}

	// Property doesn't exist — POST to create
	propsPath := client.ConfluenceAPIPath("content/" + pageID + "/property")
	_, err := client.Post(propsPath, payload, nil)
	return err
}
