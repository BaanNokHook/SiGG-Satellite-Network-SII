// SiGG-Satellite-Network-SII  //

package fetcher

import (
	"reflect"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	"github.com/apache/skywalking-satellite/plugins/fetcher/api"
)

// RegisterFetcherPlugins register the used fetcher plugins.
func RegisterFetcherPlugins() {
	plugin.RegisterPluginCategory(reflect.TypeOf((*api.Fetcher)(nil)).Elem())
	fetchers := []api.Fetcher{}
	for _, fetcher := range fetchers {
		plugin.RegisterPlugin(fetcher)
	}
}
