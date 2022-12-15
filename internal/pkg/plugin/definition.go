package plugin

// NameField is a required field in Config.
const NameField = "plugin_name"

// Plugin defines the plugin model in Satellite.
type Plugin interface {
	// Name returns the name of the specific plugin.
	Name() string
	// ShowName returns the specific name show on documentations.
	ShowName() string
	// Description returns the description of the specific plugin.
	Description() string
	// DefaultConfig returns the default config, that is a YAML pattern.
	DefaultConfig() string
}

// SharingPlugin the plugins cloud be sharing with different modules in different namespaces.
type SharingPlugin interface {
	Plugin

	// Prepare the sharing plugins, such as build the connection with the external services.
	Prepare() error
	// Start a server to receive the input APM data.
	Start() error
	// Close the sharing plugin.
	Close() error
}

// Config is used to initialize the DefaultInitializingPlugin.
type Config map[string]interface{}
