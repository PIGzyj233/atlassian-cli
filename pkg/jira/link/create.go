package link

import (
	"fmt"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/convert"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdCreate creates the `link create` command.
func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var (
		linkType string
		inward   string
		outward  string
		comment  string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an issue link",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			if err := createLink(client, linkType, inward, outward, comment); err != nil {
				return err
			}

			return output.Print(cmd, map[string]any{
				"success": true,
				"message": "Issue link created",
			})
		},
	}

	cmd.Flags().StringVar(&linkType, "type", "", "Link type name (required)")
	cmd.Flags().StringVar(&inward, "inward", "", "Inward issue key (required)")
	cmd.Flags().StringVar(&outward, "outward", "", "Outward issue key (required)")
	cmd.Flags().StringVar(&comment, "comment", "", "Comment body")
	_ = cmd.MarkFlagRequired("type")
	_ = cmd.MarkFlagRequired("inward")
	_ = cmd.MarkFlagRequired("outward")

	return cmd
}

func createLink(client *api.Client, linkType, inward, outward, comment string) error {
	body := map[string]any{
		"type":         map[string]any{"name": linkType},
		"inwardIssue":  map[string]any{"key": inward},
		"outwardIssue": map[string]any{"key": outward},
	}

	if comment != "" {
		// Cloud (v3) requires ADF; Server/DC (v2) uses plain string.
		var bodyValue any = comment
		if client.IsCloud() {
			bodyValue = convert.MarkdownToADF(comment)
		}
		body["comment"] = map[string]any{"body": bodyValue}
	}

	path := client.JiraAPIPath("issueLink")
	_, err := client.Post(path, body, nil)
	if err != nil {
		return fmt.Errorf("creating issue link: %w", err)
	}
	return nil
}
