package cmdutil

// Factory holds shared dependencies for all commands.
type Factory struct {
	// Will be populated in later tasks:
	// Config  *config.Config
	// Client  *api.Client
	// Output  string
	Version string
}

// NewFactory creates a Factory with defaults.
func NewFactory(version string) *Factory {
	return &Factory{
		Version: version,
	}
}
