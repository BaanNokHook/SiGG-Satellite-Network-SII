// SiGG-Satellite-Network-SII  //

package api

import (
	"reflect"

	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	module "github.com/apache/skywalking-satellite/internal/satellite/module/api"
	forwarder "github.com/apache/skywalking-satellite/plugins/forwarder/api"
)

// Receiver is a plugin interface, that defines new collectors.
type Receiver interface {
	plugin.Plugin

	// RegisterHandler register  a handler to the server, such as to handle a gRPC or an HTTP request
	RegisterHandler(server interface{})

	// RegisterSyncInvoker register the sync invoker, receive event and sync invoke to sender
	RegisterSyncInvoker(invoker module.SyncInvoker)

	// Channel would be put a data when the receiver receives an APM data.
	Channel() <-chan *v1.SniffData

	// SupportForwarders should provider all forwarder support current receiver
	SupportForwarders() []forwarder.Forwarder
}

// GetReceiver gets an initialized receiver plugin.
func GetReceiver(config plugin.Config) Receiver {
	return plugin.Get(reflect.TypeOf((*Receiver)(nil)).Elem(), config).(Receiver)
}
