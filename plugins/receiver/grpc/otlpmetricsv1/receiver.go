// SiGG-Satellite-Network-SII  //

package otlpmetricsv1

import (
	"github.com/apache/skywalking-satellite/internal/pkg/config"
	module "github.com/apache/skywalking-satellite/internal/satellite/module/api"
	forwarder "github.com/apache/skywalking-satellite/plugins/forwarder/api"
	"github.com/apache/skywalking-satellite/plugins/forwarder/grpc/otlpmetricsv1"
	"github.com/apache/skywalking-satellite/plugins/receiver/grpc"

	metrics "skywalking.apache.org/repo/goapi/proto/opentelemetry/proto/collector/metrics/v1"
	sniffer "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const (
	Name     = "grpc-otlp-metrics-v1-receiver"
	ShowName = "GRPC OpenTelemetry Metrics v1 Receiver"
)

type Receiver struct {
	config.CommonFields
	grpc.CommonGRPCReceiverFields
	service *MetricsService
}

func (r *Receiver) Name() string {
	return Name
}

func (r *Receiver) ShowName() string {
	return ShowName
}

func (r *Receiver) Description() string {
	return "This is a receiver for OpenTelemetry Metrics v1 format, " +
		"which is defined at https://github.com/open-telemetry/opentelemetry-proto/blob/" +
		"724e427879e3d2bae2edc0218fff06e37b9eb46e/opentelemetry/proto/collector/metrics/v1/metrics_service.proto."
}

func (r *Receiver) DefaultConfig() string {
	return ` `
}

func (r *Receiver) RegisterHandler(server interface{}) {
	r.CommonGRPCReceiverFields = *grpc.InitCommonGRPCReceiverFields(server)
	r.service = &MetricsService{receiveChannel: r.OutputChannel}
	metrics.RegisterMetricsServiceServer(r.Server, r.service)
}

func (r *Receiver) RegisterSyncInvoker(_ module.SyncInvoker) {
}

func (r *Receiver) Channel() <-chan *sniffer.SniffData {
	return r.OutputChannel
}

func (r *Receiver) SupportForwarders() []forwarder.Forwarder {
	return []forwarder.Forwarder{
		new(otlpmetricsv1.Forwarder),
	}
}
