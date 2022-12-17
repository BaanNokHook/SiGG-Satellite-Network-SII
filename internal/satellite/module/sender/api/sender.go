package api

import (
	"github.com/apache/skywalking-satellite/internal/satellite/event"
	"github.com/apache/skywalking-satellite/internal/satellite/module/api"
)

type Sender interface {
	api.Module
	api.SyncInvoker

	// InputDataChannel is a blocking channel to receive the apm data from the downstream processor module.
	InputDataChannel(partition int) chan<- *event.OutputEventContext

	// Inject the gatherer module.
	SetGatherer(g api.Module) error
}
