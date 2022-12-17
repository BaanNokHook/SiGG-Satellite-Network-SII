package api

import (
	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
)

type SenderConfig struct {
	*config.CommonFields
	// plugins config
	ForwardersConfig []plugin.Config `mapstructure:"forwarders"`  // forwarder plugins config
	FallbackerConfig plugin.Config   `mapstructure:"fallbacker"`  // fallbacker plugins config
	ClientName       string          `mapstructure:"client_name"` // client plugin name

	MaxBufferSize  int `mapstructure:"max_buffer_size"`  // the max buffer capacity
	MinFlushEvents int `mapstructure:"min_flush_events"` // the min flush events when receives a timer flush signal
	FlushTime      int `mapstructure:"flush_time"`       // the period flush time
}
