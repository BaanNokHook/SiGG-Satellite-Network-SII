// SiGG-Satellite-Network-SII  //

package client

import (
	"reflect"

	"github.com/apache/skywalking-satellite/plugins/client/grpc"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	"github.com/apache/skywalking-satellite/plugins/client/api"
	"github.com/apache/skywalking-satellite/plugins/client/kafka"
)

// RegisterClientPlugins register the used client plugins.
func RegisterClientPlugins() {
	plugin.RegisterPluginCategory(reflect.TypeOf((*api.Client)(nil)).Elem())
	clients := []api.Client{
		// Please register the client plugins at here.
		new(kafka.Client),
		new(grpc.Client),
	}
	for _, client := range clients {
		plugin.RegisterPlugin(client)
	}
}
