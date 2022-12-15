// SiGG-Satellite-Network-SII  //

package server

import (
	"reflect"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	"github.com/apache/skywalking-satellite/plugins/server/api"
	"github.com/apache/skywalking-satellite/plugins/server/grpc"
	"github.com/apache/skywalking-satellite/plugins/server/http"
)

// RegisterServerPlugins register the used server plugins.
func RegisterServerPlugins() {
	plugin.RegisterPluginCategory(reflect.TypeOf((*api.Server)(nil)).Elem())
	servers := []api.Server{
		// Please register the server plugins at here.
		new(grpc.Server),
		new(http.Server),
	}
	for _, server := range servers {
		plugin.RegisterPlugin(server)
	}
}
