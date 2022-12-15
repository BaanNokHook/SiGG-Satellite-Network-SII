// SiGG-Satellite-Network-SII  //

package nativeclr

import (
	"context"
	"fmt"
	"reflect"

	"google.golang.org/grpc"

	agent "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
)

const (
	Name     = "native-clr-grpc-forwarder"
	ShowName = "Native CLR GRPC Forwarder"
)

type Forwarder struct {
	config.CommonFields

	client agent.CLRMetricReportServiceClient
}

func (f *Forwarder) Name() string {
	return Name
}

func (f *Forwarder) ShowName() string {
	return ShowName
}

func (f *Forwarder) Description() string {
	return "This is a synchronization grpc forwarder with the SkyWalking native Configuration Discovery Service protocol."
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
	f.client = agent.NewCLRMetricReportServiceClient(client)
	return nil
}

func (f *Forwarder) Forward(batch event.BatchEvents) error {
	for _, e := range batch {
		clr := e.GetClr()
		_, err := f.client.Collect(context.Background(), clr)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Forwarder) ForwardType() v1.SniffType {
	return v1.SniffType_CLRMetricType
}

func (f *Forwarder) SyncForward(_ *v1.SniffData) (*v1.SniffData, error) {
	return nil, fmt.Errorf("unsupport sync forward")
}

func (f *Forwarder) SupportedSyncInvoke() bool {
	return false
}
