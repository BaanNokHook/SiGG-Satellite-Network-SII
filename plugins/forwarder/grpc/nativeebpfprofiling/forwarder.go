// SiGG-Satellite-Network-SII  //

package nativeebpfprofiling

import (
	"context"
	"fmt"
	"io"
	"reflect"

	"google.golang.org/grpc"

	profiling "skywalking.apache.org/repo/goapi/collect/ebpf/profiling/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/pkg/log"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
)

const (
	Name     = "native-ebpf-profiling-grpc-forwarder"
	ShowName = "Native EBPF Profiling GRPC Forwarder"
)

type Forwarder struct {
	config.CommonFields

	profilingClient profiling.EBPFProfilingServiceClient
}

func (f *Forwarder) Name() string {
	return Name
}

func (f *Forwarder) ShowName() string {
	return ShowName
}

func (f *Forwarder) Description() string {
	return "This is a synchronization grpc forwarder with the SkyWalking native process protocol."
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
	f.profilingClient = profiling.NewEBPFProfilingServiceClient(client)
	return nil
}

func (f *Forwarder) Forward(batch event.BatchEvents) error {
	stream, err := f.profilingClient.CollectProfilingData(context.Background())
	if err != nil {
		return err
	}
	for _, e := range batch {
		data, ok := e.GetData().(*v1.SniffData_EBPFProfilingDataList)
		if !ok {
			continue
		}
		for _, d := range data.EBPFProfilingDataList.DataList {
			err := stream.Send(d)
			if err != nil {
				log.Logger.Errorf("%s send log data error: %v", f.Name(), err)
				err = closeStream(stream)
				if err != nil {
					log.Logger.Errorf("%s close stream error: %v", f.Name(), err)
				}
				return err
			}
		}
	}
	return closeStream(stream)
}

func closeStream(stream profiling.EBPFProfilingService_CollectProfilingDataClient) error {
	_, err := stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

func (f *Forwarder) ForwardType() v1.SniffType {
	return v1.SniffType_EBPFProfilingType
}

func (f *Forwarder) SyncForward(e *v1.SniffData) (*v1.SniffData, error) {
	query := e.GetEBPFProfilingTaskQuery()
	if query == nil {
		return nil, fmt.Errorf("unsupport data")
	}
	commands, err := f.profilingClient.QueryTasks(context.Background(), query)
	if err != nil {
		return nil, err
	}
	return &v1.SniffData{Data: &v1.SniffData_Commands{Commands: commands}}, nil
}

func (f *Forwarder) SupportedSyncInvoke() bool {
	return true
}
