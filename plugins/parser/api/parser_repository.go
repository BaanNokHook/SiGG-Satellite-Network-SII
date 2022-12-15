// SiGG-Satellite-Network-SII  //

package api

import (
	"reflect"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
)

// GetParser an initialized filter plugin.
func GetParser(config plugin.Config) Parser {
	return plugin.Get(reflect.TypeOf((*Parser)(nil)).Elem(), config).(Parser)
}

// RegisterParserPlugins register the used parser plugins.
func RegisterParserPlugins() {
	plugin.RegisterPluginCategory(reflect.TypeOf((*Parser)(nil)).Elem())
	parsers := []Parser{
		// Please register the parser plugins at here.
	}
	for _, parser := range parsers {
		plugin.RegisterPlugin(parser)
	}
}
