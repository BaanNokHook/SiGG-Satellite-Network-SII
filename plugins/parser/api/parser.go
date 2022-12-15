// SiGG-Satellite-Network-SII  //

package api

import (
	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
)

// Parser is a plugin interface, that defines new Parsers for Collector plugin.
type Parser interface {
	plugin.Plugin

	// ParseBytes parse the byte buffer into events.
	ParseBytes(bytes []byte) (event.BatchEvents, error)

	// ParseStr parse the string into events.
	ParseStr(str string) (event.BatchEvents, error)
}
