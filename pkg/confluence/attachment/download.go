package attachment

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdDownload creates the `attachment download` command.
func NewCmdDownload(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download <attachment-id>",
		Short: "Download a Confluence attachment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			// Get attachment metadata to find download URL
			path := client.ConfluenceAPIPath("content/" + args[0])
			var meta map[string]any
			_, err = client.Get(path, &meta)
			if err != nil {
				return fmt.Errorf("getting attachment info: %w", err)
			}

			// Extract download link
			links, _ := meta["_links"].(map[string]any)
			downloadPath, _ := links["download"].(string)
			title, _ := meta["title"].(string)
			if downloadPath == "" {
				return fmt.Errorf("no download URL found for attachment %s", args[0])
			}

			// Download the file
			downloadURL := client.BaseURL + downloadPath
			req, err := http.NewRequest("GET", downloadURL, nil)
			if err != nil {
				return err
			}
			if client.Auth != nil {
				client.Auth.Apply(req)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("downloading: %w", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				return fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
			}

			// Write to local file
			if title == "" {
				title = args[0]
			}
			outFile, err := os.Create(title)
			if err != nil {
				return fmt.Errorf("creating output file: %w", err)
			}
			defer outFile.Close()

			written, err := io.Copy(outFile, resp.Body)
			if err != nil {
				return fmt.Errorf("writing file: %w", err)
			}

			absPath, _ := filepath.Abs(title)
			return output.Print(cmd, map[string]any{
				"downloaded": true,
				"filename":   title,
				"size":       written,
				"path":       absPath,
			})
		},
	}

	return cmd
}

