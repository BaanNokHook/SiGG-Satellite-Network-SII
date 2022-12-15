// SiGG-Satellite-Network-SII  //

//go:build windows

package queue

import (
	"github.com/apache/skywalking-satellite/plugins/queue/api"
	"github.com/apache/skywalking-satellite/plugins/queue/memory"
	"github.com/apache/skywalking-satellite/plugins/queue/none"
)

var queues = []api.Queue{
	// Please register the queue plugins available on Windows platform here.
	new(memory.Queue),
	new(none.Queue),
}
