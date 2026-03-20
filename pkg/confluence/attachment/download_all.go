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

// NewCmdDownloadAll creates the `attachment download-all` command.
func NewCmdDownloadAll(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-all <content-id>",
		Short: "Download all attachments for Confluence content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contentID := args[0]
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			// List all attachments
			listPath := client.ConfluenceAPIPath("content/" + contentID + "/child/attachment")
			rb := api.NewRequestBuilder(listPath).QueryInt("limit", 200)

			var listResult map[string]any
			_, err = client.Get(rb.BuildPath(), &listResult)
			if err != nil {
				return fmt.Errorf("listing attachments: %w", err)
			}

			attachments, ok := listResult["results"].([]any)
			if !ok || len(attachments) == 0 {
				return output.Print(cmd, map[string]any{
					"message":    "No attachments found",
					"content_id": contentID,
				})
			}

			// Create output directory
			dir := "attachments_" + contentID
			os.MkdirAll(dir, 0755)

			var downloaded []map[string]any
			var failed []map[string]any

			for _, attAny := range attachments {
				att, ok := attAny.(map[string]any)
				if !ok {
					continue
				}

				title, _ := att["title"].(string)
				links, _ := att["_links"].(map[string]any)
				downloadPath, _ := links["download"].(string)

				if downloadPath == "" {
					failed = append(failed, map[string]any{
						"filename": title,
						"error":    "no download URL",
					})
					continue
				}

				downloadURL := client.BaseURL + downloadPath
				outPath := filepath.Join(dir, title)

				if err := downloadToFile(downloadURL, outPath, client.Auth); err != nil {
					failed = append(failed, map[string]any{
						"filename": title,
						"error":    err.Error(),
					})
					continue
				}

				fi, _ := os.Stat(outPath)
				downloaded = append(downloaded, map[string]any{
					"filename": title,
					"path":     outPath,
					"size":     fi.Size(),
				})
			}

			return output.Print(cmd, map[string]any{
				"success":    true,
				"content_id": contentID,
				"total":      len(attachments),
				"downloaded": downloaded,
				"failed":     failed,
			})
		},
	}

	return cmd
}

func downloadToFile(url, outPath string, auth interface{ Apply(req *http.Request) error }) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	if auth != nil {
		auth.Apply(req)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}
