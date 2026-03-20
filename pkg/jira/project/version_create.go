package project

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/PigZyj2333/atlassian-cli/internal/api"
	"github.com/PigZyj2333/atlassian-cli/internal/output"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdVersionCreate creates the `project version-create` command.
func NewCmdVersionCreate(f *cmdutil.Factory) *cobra.Command {
	var (
		name        string
		startDate   string
		releaseDate string
		description string
		file        string
	)

	cmd := &cobra.Command{
		Use:   "version-create <project-key>",
		Short: "Create one or more versions in a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectKey := args[0]

			client, err := f.JiraClient(cmdutil.GetHost(cmd))
			if err != nil {
				return err
			}

			if file != "" {
				result, err := batchCreateVersions(client, projectKey, file)
				if err != nil {
					return err
				}
				return output.Print(cmd, result)
			}

			if name == "" {
				return fmt.Errorf("--name is required when --file is not provided")
			}

			result, err := createVersion(client, projectKey, name, startDate, releaseDate, description)
			if err != nil {
				return err
			}
			return output.Print(cmd, result)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Version name (required unless --file is provided)")
	cmd.Flags().StringVar(&startDate, "start-date", "", "Version start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&releaseDate, "release-date", "", "Version release date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&description, "description", "", "Version description")
	cmd.Flags().StringVar(&file, "file", "", "JSON file with array of version objects for batch creation")

	return cmd
}

func createVersion(client *api.Client, projectKey, name, startDate, releaseDate, description string) (map[string]any, error) {
	path := client.JiraAPIPath("version")

	body := map[string]any{
		"name":    name,
		"project": projectKey,
	}
	if startDate != "" {
		body["startDate"] = startDate
	}
	if releaseDate != "" {
		body["releaseDate"] = releaseDate
	}
	if description != "" {
		body["description"] = description
	}

	var result map[string]any
	_, err := client.Post(path, body, &result)
	if err != nil {
		return nil, fmt.Errorf("creating version: %w", err)
	}
	return result, nil
}

func batchCreateVersions(client *api.Client, projectKey, filePath string) (map[string]any, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", filePath, err)
	}

	var versions []map[string]any
	if err := json.Unmarshal(data, &versions); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}

	path := client.JiraAPIPath("version")
	created := make([]any, 0, len(versions))
	var errors []map[string]any

	for i, ver := range versions {
		ver["project"] = projectKey

		var result map[string]any
		_, err := client.Post(path, ver, &result)
		if err != nil {
			errors = append(errors, map[string]any{
				"index": i,
				"error": err.Error(),
			})
			continue
		}
		created = append(created, result)
	}

	return map[string]any{
		"created": created,
		"errors":  errors,
		"total":   len(versions),
	}, nil
}
