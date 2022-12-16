package api

import (
	"github.com/apache/skywalking-satellite/internal/satellite/module/api"
)

// Processor is the APM data processing module in Satellite.
type Processor interface {
	api.Module
	api.SyncInvoker

	// Inject the gatherer module.
	SetGatherer(gatherer api.Module) error
	// Inject the sender module.
	SetSender(sender api.Module) error
}
