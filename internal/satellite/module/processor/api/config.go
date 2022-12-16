package api

import (
	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
)

// ProcessorConfig contains all implementation fields.
type ProcessorConfig struct {
	*config.CommonFields

	FilterConfig []plugin.Config `mapstructure:"filters"` // filter plugins
}
