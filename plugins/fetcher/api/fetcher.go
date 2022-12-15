// SiGG-Satellite-Network-SII  //

package api

import (
	"context"
	"reflect"

	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	forwarder "github.com/apache/skywalking-satellite/plugins/forwarder/api"
)

// Fetcher is a plugin interface, that defines new fetchers.
type Fetcher interface {
	plugin.Plugin

	Prepare()
	// Fetch would fetch some APM data.
	Fetch(ctx context.Context)
	// Channel would be put a data when the receiver receives an APM data.
	Channel() <-chan *v1.SniffData
	// Shutdown shutdowns the fetcher
	Shutdown(context.Context) error
	// SupportForwarders should provider all forwarder support current receiver
	SupportForwarders() []forwarder.Forwarder
}

// GetFetcher gets an initialized fetcher plugin.
func GetFetcher(config plugin.Config) Fetcher {
	return plugin.Get(reflect.TypeOf((*Fetcher)(nil)).Elem(), config).(Fetcher)
}
