// SiGG-Satellite-Network-SII  //

package api

import (
	"reflect"

	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
)

// Forwarder is a plugin interface, that defines new forwarders.
type Forwarder interface {
	plugin.Plugin
	// Prepare do some preparation works, such as create a stub in gRPC and create a producer in Kafka.
	Prepare(connection interface{}) error
	// Forward the batch events to the external services, such as Kafka MQ and SkyWalking OAP cluster.
	Forward(batch event.BatchEvents) error
	// SyncForward the single event to the external service with sync forward
	SyncForward(event *v1.SniffData) (*v1.SniffData, error)
	// ForwardType returns the supported event type.
	ForwardType() v1.SniffType
	// SupportedSyncInvoke return is support SyncForward
	SupportedSyncInvoke() bool
}

// ForwardFunc represent the Forward() in Forwarder
type ForwardFunc func(batch event.BatchEvents) error

// GetForwarder an initialized filter plugin.
func GetForwarder(config plugin.Config) Forwarder {
	return plugin.Get(reflect.TypeOf((*Forwarder)(nil)).Elem(), config).(Forwarder)
}
