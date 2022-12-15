// SiGG-Satellite-Network-SII  //

package envoymetricsv2

import (
	"context"
	"fmt"
	"io"
	"reflect"

	v2 "skywalking.apache.org/repo/goapi/proto/envoy/service/metrics/v2"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	"google.golang.org/grpc"

	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/pkg/log"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
)

const (
	Name     = "envoy-metrics-v2-grpc-forwarder"
	ShowName = "Envoy Metrics v2 GRPC Forwarder"
)

type Forwarder struct {
	config.CommonFields
	metricsClient v2.MetricsServiceClient
}

func (f *Forwarder) Name() string {
	return Name
}

func (f *Forwarder) ShowName() string {
	return ShowName
}

func (f *Forwarder) Description() string {
	return "This is a synchronization Metrics v2 grpc forwarder with the Envoy metrics protocol."
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
	f.metricsClient = v2.NewMetricsServiceClient(client)
	return nil
}

func (f *Forwarder) Forward(batch event.BatchEvents) error {
	for _, e := range batch {
		data, ok := e.GetData().(*v1.SniffData_EnvoyMetricsV2List)
		if !ok {
			continue
		}
		stream, err := f.metricsClient.StreamMetrics(context.Background())
		if err != nil {
			log.Logger.Errorf("open grpc stream error %v", err)
			return err
		}
		for _, message := range data.EnvoyMetricsV2List.Messages {
			err := stream.Send(message)
			if err != nil {
				log.Logger.Errorf("%s send envoy metrics v2 data error: %v", f.Name(), err)
				f.closeStream(stream)
				return err
			}
		}
		f.closeStream(stream)
	}
	return nil
}

func (f *Forwarder) closeStream(stream v2.MetricsService_StreamMetricsClient) {
	_, err := stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		log.Logger.Warnf("%s close stream error: %v", f.Name(), err)
	}
}

func (f *Forwarder) ForwardType() v1.SniffType {
	return v1.SniffType_EnvoyMetricsV2Type
}

func (f *Forwarder) SyncForward(_ *v1.SniffData) (*v1.SniffData, error) {
	return nil, fmt.Errorf("unsupport sync forward")
}

func (f *Forwarder) SupportedSyncInvoke() bool {
	return false
}
