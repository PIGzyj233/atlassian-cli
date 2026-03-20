package attachment

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdUpload creates the `attachment upload` command.
func NewCmdUpload(f *cmdutil.Factory) *cobra.Command {
	var (
		filePath  string
		comment   string
		minorEdit bool
	)

	cmd := &cobra.Command{
		Use:   "upload <content-id>",
		Short: "Upload an attachment to Confluence content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contentID := args[0]
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			result, err := uploadFile(client.BaseURL, client.Auth, contentID, filePath, comment, minorEdit)
			if err != nil {
				return err
			}

			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "Path to file to upload (required)")
	cmd.Flags().StringVar(&comment, "comment", "", "Attachment comment")
	cmd.Flags().BoolVar(&minorEdit, "minor-edit", true, "Mark as minor edit")
	cmd.MarkFlagRequired("file")

	return cmd
}

// uploadFile performs a multipart form upload to the Confluence attachment API.
func uploadFile(baseURL string, auth interface{ Apply(req *http.Request) error }, contentID, filePath, comment string, minorEdit bool) (map[string]any, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("creating form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("copying file: %w", err)
	}

	if comment != "" {
		writer.WriteField("comment", comment)
	}
	writer.WriteField("minorEdit", fmt.Sprintf("%t", minorEdit))
	writer.Close()

	url := baseURL + "/rest/api/content/" + contentID + "/child/attachment"
	req, err := http.NewRequest("PUT", url, &buf)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Atlassian-Token", "nocheck")

	if auth != nil {
		if err := auth.Apply(req); err != nil {
			return nil, fmt.Errorf("applying auth: %w", err)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("uploading: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	fi, _ := os.Stat(filePath)
	return map[string]any{
		"success":    true,
		"content_id": contentID,
		"filename":   filepath.Base(filePath),
		"size":       fi.Size(),
	}, nil
}
