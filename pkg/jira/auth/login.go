package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/PigZyj2333/atlassian-cli/internal/config"
	"github.com/PigZyj2333/atlassian-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdLogin creates the `auth login` command.
func NewCmdLogin(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with an Atlassian instance",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(f)
		},
	}
	return cmd
}

func runLogin(f *cmdutil.Factory) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Atlassian instance hostname (e.g., company.atlassian.net): ")
	hostname, _ := reader.ReadString('\n')
	hostname = strings.TrimSpace(hostname)
	if hostname == "" {
		return fmt.Errorf("hostname is required")
	}

	instType := "server"
	if config.IsCloudHost(hostname) {
		instType = "cloud"
	}
	fmt.Fprintf(os.Stderr, "Detected instance type: %s\n", instType)

	fmt.Print("Auth method (basic/pat) [basic]: ")
	authMethod, _ := reader.ReadString('\n')
	authMethod = strings.TrimSpace(authMethod)
	if authMethod == "" {
		authMethod = "basic"
	}

	host := &config.HostConfig{
		Type:       instType,
		Auth:       authMethod,
		Jira:       config.ServiceConfig{Enabled: true},
		Confluence: config.ServiceConfig{Enabled: true},
	}

	if instType == "cloud" {
		host.Confluence.BasePath = "/wiki"
	}

	switch authMethod {
	case "basic":
		fmt.Print("Username/email: ")
		username, _ := reader.ReadString('\n')
		host.Username = strings.TrimSpace(username)

		fmt.Print("API token: ")
		token, _ := reader.ReadString('\n')
		host.Token = strings.TrimSpace(token)
	case "pat":
		fmt.Print("Personal Access Token: ")
		token, _ := reader.ReadString('\n')
		host.Token = strings.TrimSpace(token)
	default:
		return fmt.Errorf("unsupported auth method: %s", authMethod)
	}

	// Validate credentials by building authenticator
	if _, err := config.NewAuthenticator(host); err != nil {
		return fmt.Errorf("invalid credentials: %w", err)
	}

	f.Config.Hosts[hostname] = host
	if f.Config.Defaults.Host == "" {
		f.Config.Defaults.Host = hostname
	}

	if err := f.Config.Save(f.ConfigPath); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Logged in to %s\n", hostname)
	return nil
}
