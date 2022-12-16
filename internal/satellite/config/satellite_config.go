package config

import (
	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/pkg/log"
	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	gatherer "github.com/apache/skywalking-satellite/internal/satellite/module/gatherer/api"
	processor "github.com/apache/skywalking-satellite/internal/satellite/module/processor/api"
	sender "github.com/apache/skywalking-satellite/internal/satellite/module/sender/api"
	"github.com/apache/skywalking-satellite/internal/satellite/telemetry"
)

// SatelliteConfig is to initialize Satellite.
type SatelliteConfig struct {
	Logger    *log.LoggerConfig `mapstructure:"logger"`
	Pipes     []*PipeConfig     `mapstructure:"pipes"`
	Sharing   *SharingConfig    `mapstructure:"sharing"`
	Telemetry *telemetry.Config `mapstructure:"telemetry"`
}

// SharingConfig contains some plugins,which could be shared by every namespace. That is useful to reduce resources cost.
type SharingConfig struct {
	Clients             []plugin.Config      `mapstructure:"clients"`
	Servers             []plugin.Config      `mapstructure:"servers"`
	SharingCommonConfig *config.CommonFields `mapstructure:"common_config"`
}

// PipeConfig initializes the different module in different namespace.
type PipeConfig struct {
	PipeCommonConfig *config.CommonFields       `mapstructure:"common_config"`
	Gatherer         *gatherer.GathererConfig   `mapstructure:"gatherer"`
	Processor        *processor.ProcessorConfig `mapstructure:"processor"`
	Sender           *sender.SenderConfig       `mapstructure:"sender"`
}

// NewDefaultSatelliteConfig creates a satellite config with default value.
func NewDefaultSatelliteConfig() *SatelliteConfig {
	return &SatelliteConfig{
		Logger: &log.LoggerConfig{
			LogPattern:  "%time [%level][%field] - %msg",
			TimePattern: "2006-01-02 15:04:05.000",
			Level:       "info",
		},
		Telemetry: &telemetry.Config{
			Cluster:  "default_cluster",
			Service:  "default_service",
			Instance: "default_instance",
		},
		Sharing: &SharingConfig{
			SharingCommonConfig: &config.CommonFields{
				PipeName: "sharing",
			},
		},
	}
}
