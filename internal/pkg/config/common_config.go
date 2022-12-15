package config

const (
	CommonFieldsName = "CommonFields"
	TagName          = "mapstructure"
)

// CommonFields defines some common fields for every module or plugin.
type CommonFields struct {
	// PipeName indicates which pipe it belongs to.
	PipeName string `mapstructure:"pipe_name"`
}
