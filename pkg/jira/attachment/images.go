package attachment

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdImages creates the `attachment images` command.
func NewCmdImages(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images <issue-key>",
		Short: "Get base64-encoded image attachments from an issue",
		Long:  "Fetches image attachments from a Jira issue and returns them as base64-encoded data, useful for LLM tools that need to process images.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := fetchImageAttachments(client, args[0])
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	return cmd
}

func fetchImageAttachments(client *api.Client, issueKey string) (map[string]any, error) {
	// Fetch issue with attachment field only
	path := api.NewRequestBuilder(client.JiraAPIPath("issue/" + issueKey)).
		Query("fields", "attachment").
		BuildPath()

	var issue map[string]any
	_, err := client.Get(path, &issue)
	if err != nil {
		return nil, fmt.Errorf("fetching attachments for %s: %w", issueKey, err)
	}

	// Extract attachments from the response
	fields, ok := issue["fields"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unexpected response: missing fields")
	}

	rawAttachments, ok := fields["attachment"].([]any)
	if !ok || len(rawAttachments) == 0 {
		return map[string]any{
			"issue":  issueKey,
			"count":  0,
			"images": []map[string]any{},
		}, nil
	}

	// Filter to image attachments and download+encode each
	var images []map[string]any
	for _, raw := range rawAttachments {
		att, ok := raw.(map[string]any)
		if !ok {
			continue
		}

		mimeType, _ := att["mimeType"].(string)
		if !strings.HasPrefix(mimeType, "image/") {
			continue
		}

		filename, _ := att["filename"].(string)
		contentURL, _ := att["content"].(string)
		size, _ := att["size"].(float64) // JSON numbers decode as float64

		if filename == "" || contentURL == "" {
			continue
		}

		b64Data, err := downloadAndEncode(client, contentURL)
		if err != nil {
			// Include the image entry with an error instead of failing entirely
			images = append(images, map[string]any{
				"filename": filename,
				"mimeType": mimeType,
				"size":     int64(size),
				"error":    err.Error(),
			})
			continue
		}

		images = append(images, map[string]any{
			"filename":   filename,
			"mimeType":   mimeType,
			"size":       int64(size),
			"base64Data": b64Data,
		})
	}

	return map[string]any{
		"issue":  issueKey,
		"count":  len(images),
		"images": images,
	}, nil
}

// downloadAndEncode fetches the content from the given URL using the
// client's authentication and returns the base64-encoded data.
func downloadAndEncode(client *api.Client, contentURL string) (string, error) {
	req, err := http.NewRequest("GET", contentURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating download request: %w", err)
	}

	if client.Auth != nil {
		if err := client.Auth.Apply(req); err != nil {
			return "", fmt.Errorf("applying auth: %w", err)
		}
	}

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("downloading image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download returned HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading image data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}
