// SiGG-Satellite-Network-SII  //

package plugins

import (
	"github.com/apache/skywalking-satellite/plugins/client"
	"github.com/apache/skywalking-satellite/plugins/fallbacker"
	fetcher "github.com/apache/skywalking-satellite/plugins/fetcher"
	filter "github.com/apache/skywalking-satellite/plugins/filter/api"
	"github.com/apache/skywalking-satellite/plugins/forwarder"
	parser "github.com/apache/skywalking-satellite/plugins/parser/api"
	"github.com/apache/skywalking-satellite/plugins/queue"
	"github.com/apache/skywalking-satellite/plugins/receiver"
	"github.com/apache/skywalking-satellite/plugins/server"
)

// RegisterPlugins register the whole supported plugin category and plugin types to the registry.
func RegisterPlugins() {
	// plugins
	filter.RegisterFilterPlugins()
	forwarder.RegisterForwarderPlugins()
	parser.RegisterParserPlugins()
	queue.RegisterQueuePlugins()
	receiver.RegisterReceiverPlugins()
	fetcher.RegisterFetcherPlugins()
	fallbacker.RegisterFallbackerPlugins()
	// sharing plugins
	server.RegisterServerPlugins()
	client.RegisterClientPlugins()
}
