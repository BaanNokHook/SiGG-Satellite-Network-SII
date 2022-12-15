// SiGG-Satellite-Network-SII  //

package nativeevent

import (
	"io"
	"time"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	nativeevent "skywalking.apache.org/repo/goapi/collect/event/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const eventName = "grpc-nativeevent-event"

type EventService struct {
	receiveChannel chan *v1.SniffData
	nativeevent.UnimplementedEventServiceServer
}

func (e *EventService) Collect(stream nativeevent.EventService_CollectServer) error {
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&common.Commands{})
		}
		if err != nil {
			return err
		}
		d := &v1.SniffData{
			Name:      eventName,
			Timestamp: time.Now().UnixNano() / 1e6,
			Meta:      nil,
			Type:      v1.SniffType_EventType,
			Remote:    true,
			Data: &v1.SniffData_Event{
				Event: item,
			},
		}
		e.receiveChannel <- d
	}
}
