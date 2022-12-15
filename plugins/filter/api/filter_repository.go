// SiGG-Satellite-Network-SII  //

package api

import (
	"reflect"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
)

// GetFilter an initialized filter plugin.
func GetFilter(config plugin.Config) Filter {
	return plugin.Get(reflect.TypeOf((*Filter)(nil)).Elem(), config).(Filter)
}

// RegisterFilterPlugins register the used filter plugins.
func RegisterFilterPlugins() {
	plugin.RegisterPluginCategory(reflect.TypeOf((*Filter)(nil)).Elem())
	filters := []Filter{
		// Please register the filter plugins at here.
	}
	for _, filter := range filters {
		plugin.RegisterPlugin(filter)
	}
}
