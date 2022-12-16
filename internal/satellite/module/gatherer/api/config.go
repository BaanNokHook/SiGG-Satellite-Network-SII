package api

import (
	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
)

// GathererConfig contains all implementation fields.
type GathererConfig struct {
	// common config
	*config.CommonFields
	QueueConfig plugin.Config `mapstructure:"queue"` // queue plugin config

	// ReceiverGatherer
	ReceiverConfig plugin.Config `mapstructure:"receiver"`    // collector plugin config
	ServerName     string        `mapstructure:"server_name"` // depends on which server

	// FetcherGatherer
	FetcherConfig plugin.Config `mapstructure:"fetcher"` // fetcher plugin config
}
