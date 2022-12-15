// SiGG-Satellite-Network-SII  //

package fallbacker

import (
	"reflect"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	"github.com/apache/skywalking-satellite/plugins/fallbacker/api"
	"github.com/apache/skywalking-satellite/plugins/fallbacker/none"
	"github.com/apache/skywalking-satellite/plugins/fallbacker/timer"
)

// RegisterFallbackerPlugins register the used fallbacker plugins.
func RegisterFallbackerPlugins() {
	plugin.RegisterPluginCategory(reflect.TypeOf((*api.Fallbacker)(nil)).Elem())
	fallbackers := []api.Fallbacker{
		// Please register the fallbacker plugins at here.
		new(none.Fallbacker),
		new(timer.Fallbacker),
	}
	for _, fallbacker := range fallbackers {
		plugin.RegisterPlugin(fallbacker)
	}
}
