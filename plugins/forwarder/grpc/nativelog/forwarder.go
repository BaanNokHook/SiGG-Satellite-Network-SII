// SiGG-Satellite-Network-SII  //

package nativelog

import (
	"context"
	"fmt"
	"io"
	"reflect"

	"google.golang.org/grpc"

	logging "skywalking.apache.org/repo/goapi/collect/logging/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/pkg/log"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
	server_grpc "github.com/apache/skywalking-satellite/plugins/server/grpc"
)

const (
	Name     = "native-log-grpc-forwarder"
	ShowName = "Native Log GRPC Forwarder"
)

type Forwarder struct {
	config.CommonFields

	logClient logging.LogReportServiceClient
}

func (f *Forwarder) Name() string {
	return Name
}

func (f *Forwarder) ShowName() string {
	return ShowName
}

func (f *Forwarder) Description() string {
	return "This is a synchronization grpc forwarder with the SkyWalking native log protocol."
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
	f.logClient = logging.NewLogReportServiceClient(client)
	return nil
}

func (f *Forwarder) Forward(batch event.BatchEvents) error {
	for _, e := range batch {
		data, ok := e.GetData().(*v1.SniffData_LogList)
		if !ok {
			continue
		}
		stream, err := f.logClient.Collect(context.Background())
		if err != nil {
			log.Logger.Errorf("open grpc stream error %v", err)
			return err
		}
		streamClosed := false
		for _, logData := range data.LogList.Logs {
			err := stream.SendMsg(server_grpc.NewOriginalData(logData))
			if err != nil {
				log.Logger.Errorf("%s send log data error: %v", f.Name(), err)
				f.closeStream(stream)
				streamClosed = true
				break
			}
		}

		if !streamClosed {
			f.closeStream(stream)
		}
	}
	return nil
}

func (f *Forwarder) closeStream(stream logging.LogReportService_CollectClient) {
	_, err := stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		if err != nil {
			log.Logger.Errorf("%s close stream error: %v", f.Name(), err)
		}
	}
}

func (f *Forwarder) ForwardType() v1.SniffType {
	return v1.SniffType_Logging
}

func (f *Forwarder) SyncForward(*v1.SniffData) (*v1.SniffData, error) {
	return nil, fmt.Errorf("unsupport sync forward")
}

func (f *Forwarder) SupportedSyncInvoke() bool {
	return false
}
