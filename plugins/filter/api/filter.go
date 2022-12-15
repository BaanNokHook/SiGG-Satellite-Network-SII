// SiGG-Satellite-Network-SII  //

package api

import (
	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
)

// Filter is a plugin interface, that defines new pipeline filters.
type Filter interface {
	plugin.Plugin

	// Process would put the needed event to the OutputEventContext.
	Process(context *event.OutputEventContext)
}
