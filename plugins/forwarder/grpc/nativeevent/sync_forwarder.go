// SiGG-Satellite-Network-SII  //

package nativeevent

import (
	"context"
	"fmt"
	"io"
	"reflect"

	"google.golang.org/grpc"

	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/pkg/log"
	"github.com/apache/skywalking-satellite/internal/satellite/event"

	nativeevent "skywalking.apache.org/repo/goapi/collect/event/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const (
	Name     = "native-event-grpc-forwarder"
	ShowName = "Native Event GRPC Forwarder"
)

type Forwarder struct {
	config.CommonFields
	client nativeevent.EventServiceClient
}

func (f *Forwarder) Name() string {
	return Name
}

func (f *Forwarder) ShowName() string {
	return ShowName
}

func (f *Forwarder) Description() string {
	return "This is a synchronization grpc forwarder with the SkyWalking native event protocol."
}

func (f *Forwarder) DefaultConfig() string {
	return ``
}

func (f *Forwarder) Prepare(connection interface{}) error {
	client, ok := connection.(*grpc.ClientConn)
	if !ok {
		return fmt.Errorf("the %s only accepts a grpc client, but received a %s",
			f.Name(), reflect.TypeOf(connection).String())
	}
	f.client = nativeevent.NewEventServiceClient(client)
	return nil
}

func (f *Forwarder) Forward(batch event.BatchEvents) error {
	stream, err := f.client.Collect(context.Background())
	if err != nil {
		log.Logger.Errorf("open grpc stream error %v", err)
		return err
	}
	for _, e := range batch {
		data, ok := e.GetData().(*v1.SniffData_Event)
		if !ok {
			continue
		}
		err := stream.Send(data.Event)
		if err != nil {
			log.Logger.Errorf("%s send log data error: %v", f.Name(), err)
			err = closeStream(stream)
			if err != nil {
				log.Logger.Errorf("%s close stream error: %v", f.Name(), err)
			}
			return err
		}
	}
	return closeStream(stream)
}

func closeStream(stream nativeevent.EventService_CollectClient) error {
	_, err := stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (f *Forwarder) ForwardType() v1.SniffType {
	return v1.SniffType_EventType
}

func (f *Forwarder) SyncForward(_ *v1.SniffData) (*v1.SniffData, error) {
	return nil, fmt.Errorf("unsupport sync forward")
}

func (f *Forwarder) SupportedSyncInvoke() bool {
	return false
}
