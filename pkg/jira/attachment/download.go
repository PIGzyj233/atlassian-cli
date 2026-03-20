package attachment

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdDownload creates the `attachment download` command.
func NewCmdDownload(f *cmdutil.Factory) *cobra.Command {
	var outputDir string

	cmd := &cobra.Command{
		Use:   "download <issue-key>",
		Short: "Download all attachments from an issue",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := downloadAttachments(client, args[0], outputDir)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&outputDir, "output-dir", ".", "Directory to save downloaded files")
	return cmd
}

func downloadAttachments(client *api.Client, issueKey, outputDir string) (map[string]any, error) {
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
			"issue":       issueKey,
			"downloaded":  0,
			"attachments": []map[string]any{},
		}, nil
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("creating output directory: %w", err)
	}

	var results []map[string]any
	for _, raw := range rawAttachments {
		att, ok := raw.(map[string]any)
		if !ok {
			continue
		}

		filename, _ := att["filename"].(string)
		contentURL, _ := att["content"].(string)
		if filename == "" || contentURL == "" {
			continue
		}

		err := downloadFile(client, contentURL, outputDir, filename)
		status := "ok"
		errMsg := ""
		if err != nil {
			status = "error"
			errMsg = err.Error()
		}

		entry := map[string]any{
			"filename": filename,
			"status":   status,
		}
		if errMsg != "" {
			entry["error"] = errMsg
		}
		results = append(results, entry)
	}

	return map[string]any{
		"issue":       issueKey,
		"outputDir":   outputDir,
		"downloaded":  len(results),
		"attachments": results,
	}, nil
}

// downloadFile performs an authenticated GET request to the content URL
// and saves the response body to a file in the specified directory.
func downloadFile(client *api.Client, contentURL, outputDir, filename string) error {
	req, err := http.NewRequest("GET", contentURL, nil)
	if err != nil {
		return fmt.Errorf("creating download request: %w", err)
	}

	if client.Auth != nil {
		if err := client.Auth.Apply(req); err != nil {
			return fmt.Errorf("applying auth: %w", err)
		}
	}

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("downloading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned HTTP %d", resp.StatusCode)
	}

	filePath := filepath.Join(outputDir, filename)
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("creating file %s: %w", filePath, err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("writing file %s: %w", filePath, err)
	}
	return nil
}
