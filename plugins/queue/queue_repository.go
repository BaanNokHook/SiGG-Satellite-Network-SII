// SiGG-Satellite-Network-SII  //

package queue

import (
	"reflect"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	"github.com/apache/skywalking-satellite/plugins/queue/api"
)

// RegisterQueuePlugins register the used queue plugins.
func RegisterQueuePlugins() {
	plugin.RegisterPluginCategory(reflect.TypeOf((*api.Queue)(nil)).Elem())
	for _, q := range queues {
		plugin.RegisterPlugin(q)
	}
}
