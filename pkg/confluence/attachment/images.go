package attachment

import (
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// imageMediaTypes are MIME types considered as images.
var imageMediaTypes = []string{
	"image/png", "image/jpeg", "image/gif", "image/svg+xml",
	"image/webp", "image/bmp", "image/tiff",
}

// NewCmdImages creates the `attachment images` command.
func NewCmdImages(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images <content-id>",
		Short: "List image attachments for Confluence content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.ConfluenceClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			path := client.ConfluenceAPIPath("content/" + args[0] + "/child/attachment")
			rb := api.NewRequestBuilder(path).QueryInt("limit", 200)

			var result map[string]any
			_, err = client.Get(rb.BuildPath(), &result)
			if err != nil {
				return err
			}

			// Filter to image media types only
			attachments, _ := result["results"].([]any)
			var images []any
			for _, attAny := range attachments {
				att, ok := attAny.(map[string]any)
				if !ok {
					continue
				}
				mediaType, _ := att["mediaType"].(string)
				if isImageType(mediaType) {
					images = append(images, att)
				}
			}

			result["results"] = images
			result["size"] = len(images)

			return output.Print(cmd, result)
		},
	}

	return cmd
}

func isImageType(mediaType string) bool {
	mt := strings.ToLower(mediaType)
	for _, t := range imageMediaTypes {
		if mt == t {
			return true
		}
	}
	return false
}
