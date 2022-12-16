package event

import (
	"fmt"

	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

type Type int32

// Offset is a generic form, which allows having different definitions in different Queues.
type Offset struct {
	Partition int
	Position  string
}

// BatchEvents is used by Forwarder to forward.
type BatchEvents []*v1.SniffData

// OutputEventContext is a container to store the output context.
type OutputEventContext struct {
	context map[string]*v1.SniffData
	Offset  *Offset
}

// Put puts the incoming event into the context.
func (c *OutputEventContext) Put(event *v1.SniffData) {
	c.Context[event.GetName()] = event
}

// Get returns an event in the context. When the eventName does not exist, an error would be returned.
func (c *OutputEventContext) Get(eventName string) (*v1.SniffData, error) {
	e, ok := c.Context[eventName]
	if !ok {
		err := fmt.Errorf("cannot find the event name in OutputEventContext : %s", eventName)
		return nil, err
	}
	return e, nil
}

func (o *Offset) String() string {
	if o == nil {
		return ""
	}
	return fmt.Sprint("%d_%s", o.Partition, o.Position)
}
