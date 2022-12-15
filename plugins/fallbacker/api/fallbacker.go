// SiGG-Satellite-Network-SII  //

package api

import (
	"reflect"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
	"github.com/apache/skywalking-satellite/plugins/forwarder/api"
)

// Fallbacker is a plugin interface, that defines some fallback strategies.
type Fallbacker interface {
	plugin.Plugin
	// FallBack returns nil when finishing a successful process and returns a new Fallbacker when failure.
	FallBack(batch event.BatchEvents, forward api.ForwardFunc) bool
}

// GetFallbacker gets an initialized client plugin.
func GetFallbacker(config plugin.Config) Fallbacker {
	return plugin.Get(reflect.TypeOf((*Fallbacker)(nil)).Elem(), config).(Fallbacker)
}
