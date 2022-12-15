// SiGG-Satellite-Network-SII  //

package nativeprocess

import (
	"context"
	"fmt"
	"reflect"

	"google.golang.org/grpc"

	process "skywalking.apache.org/repo/goapi/collect/ebpf/profiling/process/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
)

const (
	Name     = "native-process-grpc-forwarder"
	ShowName = "Native Process GRPC Forwarder"
)

type Forwarder struct {
	config.CommonFields

	processClient process.EBPFProcessServiceClient
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
	f.processClient = process.NewEBPFProcessServiceClient(client)
	return nil
}

func (f *Forwarder) Forward(batch event.BatchEvents) error {
	for _, e := range batch {
		data, ok := e.GetData().(*v1.SniffData_EBPFProcessPingPkgList)
		if !ok {
			continue
		}
		_, err := f.processClient.KeepAlive(context.Background(), data.EBPFProcessPingPkgList)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Forwarder) ForwardType() v1.SniffType {
	return v1.SniffType_EBPFProcessType
}

func (f *Forwarder) SyncForward(e *v1.SniffData) (*v1.SniffData, error) {
	report := e.GetEBPFProcessReportList()
	if report == nil {
		return nil, fmt.Errorf("unsupport data")
	}
	processes, err := f.processClient.ReportProcesses(context.Background(), report)
	if err != nil {
		return nil, err
	}
	return &v1.SniffData{Data: &v1.SniffData_EBPFReportProcessDownstream{EBPFReportProcessDownstream: processes}}, nil
}

func (f *Forwarder) SupportedSyncInvoke() bool {
	return true
}
