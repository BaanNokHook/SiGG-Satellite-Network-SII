package gatherer

import (
	"github.com/apache/skywalking-satellite/internal/satellite/module/gatherer/api"
	"github.com/apache/skywalking-satellite/internal/satellite/sharing"
	fetcher "github.com/apache/skywalking-satellite/plugins/fetcher/api"
	"github.com/apache/skywalking-satellite/plugins/queue/partition"
	receiver "github.com/apache/skywalking-satellite/plugins/receiver/api"
	server "github.com/apache/skywalking-satellite/plugins/server/api"
)

// NewGatherer returns a gatherer module
func NewGatherer(cfg *api.GathererConfig) api.Gatherer {
	if cfg.ReceiverConfig != nil {
		return newReceiverGatherer(cfg)
	} else if cfg.FetcherConfig != nil {
		return newFetcherGatherer(cfg)
	}
	return nil
}

// newFetcherGatherer crates a gatherer with the fetcher role.
func newFetcherGatherer(cfg *api.GathererConfig) *FetcherGatherer {
	return &FetcherGatherer{
		config:         cfg,
		runningQueue:   partition.NewPartitionQueue(cfg.QueueConfig),
		runningFetcher: fetcher.GetFetcher(cfg.FetcherConfig),
	}
}

// newReceiverGatherer crates a gatherer with the receiver role.
func newReceiverGatherer(cfg *api.GathererConfig) *ReceiverGatherer {
	return &ReceiverGatherer{
		config:          cfg,
		runningQueue:    partition.NewPartitionQueue(cfg.QueueConfig),
		runningReceiver: receiver.GetReceiver(cfg.ReceiverConfig),
		runningServer:   sharing.Manager[cfg.ServerName].(server.Server),
	}
}
