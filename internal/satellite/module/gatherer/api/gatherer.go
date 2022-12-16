package api

import (
	"github.com/apache/skywalking-satellite/internal/satellite/event"
	"github.com/apache/skywalking-satellite/internal/satellite/module/api"
	queue "github.com/apache/skywalking-satellite/plugins/queue/api"
)

// Gatherer is the APM data collection module in Satellite.
type Gatherer interface {
	api.Module
	// PartitionCount is the all partition counter of gatherer. All event is partitioned.
	PartitionCount() int
	// OutputDataChannel is a blocking channel to transfer the apm data to the upstream processor module.
	OutputDataChannel(partition int) <-chan *queue.SequenceEvent
	// Ack the sent offset.
	Ack(lastOffset *event.Offset)
	// Inject the Processor module.
	SetProcessor(processor api.Module) error
}
