package page

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/convert"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdGet creates the `page get` command.
func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var (
		pageID          string
		title           string
		space           string
		includeMetadata bool
		raw             bool
	)

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a Confluence page",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			if pageID == "" && (title == "" || space == "") {
				return fmt.Errorf("provide --id, or both --title and --space")
			}

			var result map[string]any

			if pageID != "" {
				result, err = fetchPageByID(client, pageID, raw, includeMetadata)
			} else {
				result, err = fetchPageByTitle(client, space, title, raw, includeMetadata)
			}
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&pageID, "id", "", "Page ID")
	cmd.Flags().StringVar(&title, "title", "", "Page title")
	cmd.Flags().StringVar(&space, "space", "", "Space key")
	cmd.Flags().BoolVar(&includeMetadata, "include-metadata", false, "Include page metadata (labels, properties, history)")
	cmd.Flags().BoolVar(&raw, "raw", false, "Return raw storage format (no Markdown conversion)")

	return cmd
}

func fetchPageByID(client *api.Client, pageID string, raw, includeMetadata bool) (map[string]any, error) {
	path := client.ConfluenceAPIPath("content/" + pageID)
	expand := "body.storage,version,space,children.attachment"
	if includeMetadata {
		expand += ",metadata.labels,metadata.properties,history"
	}
	rb := api.NewRequestBuilder(path).
		Query("expand", expand)

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("getting page %s: %w", pageID, err)
	}

	if !raw {
		convertStorageBody(result)
	}

	return result, nil
}

func fetchPageByTitle(client *api.Client, space, title string, raw, includeMetadata bool) (map[string]any, error) {
	path := client.ConfluenceAPIPath("content")
	expand := "body.storage,version,space"
	if includeMetadata {
		expand += ",metadata.labels,metadata.properties,history"
	}
	rb := api.NewRequestBuilder(path).
		Query("spaceKey", space).
		Query("title", title).
		Query("expand", expand)

	var result map[string]any
	_, err := client.Get(rb.BuildPath(), &result)
	if err != nil {
		return nil, fmt.Errorf("searching page '%s' in space %s: %w", title, space, err)
	}

	// Extract first result
	results, ok := result["results"].([]any)
	if !ok || len(results) == 0 {
		return nil, fmt.Errorf("page '%s' not found in space %s", title, space)
	}

	page, ok := results[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected page data format")
	}
	if !raw {
		convertStorageBody(page)
	}

	return page, nil
}

// convertStorageBody extracts body.storage.value and adds content_markdown field.
func convertStorageBody(page map[string]any) {
	body, ok := page["body"].(map[string]any)
	if !ok {
		return
	}
	storage, ok := body["storage"].(map[string]any)
	if !ok {
		return
	}
	value, ok := storage["value"].(string)
	if !ok {
		return
	}
	page["content_markdown"] = convert.StorageToMarkdown(value)
}
