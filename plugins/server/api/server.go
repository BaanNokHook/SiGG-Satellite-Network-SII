// SiGG-Satellite-Network-SII  //

package api

import (
	"reflect"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
)

// Server is a plugin interface, that defines new servers, such as gRPC server and http server.
type Server interface {
	plugin.SharingPlugin
	// GetServer returns the listener server.
	GetServer() interface{}
}

// GetServer gets an initialized server plugin.
func GetServer(config plugin.Config) Server {
	return plugin.Get(reflect.TypeOf((*Server)(nil)).Elem(), config).(Server)
}
