// SiGG-Satellite-Network-SII  //

package nativetracing

import (
	"github.com/apache/skywalking-satellite/internal/pkg/config"
	module "github.com/apache/skywalking-satellite/internal/satellite/module/api"
	forwarder "github.com/apache/skywalking-satellite/plugins/forwarder/api"
	frowarder_nativetracing "github.com/apache/skywalking-satellite/plugins/forwarder/grpc/nativetracing"
	"github.com/apache/skywalking-satellite/plugins/receiver/grpc"

	v3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	v3_compat "skywalking.apache.org/repo/goapi/collect/language/agent/v3/compat"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const (
	Name     = "grpc-native-tracing-receiver"
	ShowName = "GRPC Native Tracing Receiver"
)

type Receiver struct {
	config.CommonFields
	grpc.CommonGRPCReceiverFields
	traceService           *TraceSegmentReportService
	spanAttachEventService *SpanAttachedEventReportService
}

func (r *Receiver) Name() string {
	return Name
}

func (r *Receiver) ShowName() string {
	return ShowName
}

func (r *Receiver) Description() string {
	return "This is a receiver for SkyWalking native tracing and span attached event format, " +
		"which is defined at https://github.com/apache/skywalking-data-collect-protocol/blob/master/language-agent/Tracing.proto."
}

func (r *Receiver) DefaultConfig() string {
	return ""
}

func (r *Receiver) RegisterHandler(server interface{}) {
	r.CommonGRPCReceiverFields = *grpc.InitCommonGRPCReceiverFields(server)
	r.traceService = &TraceSegmentReportService{receiveChannel: r.OutputChannel}
	r.spanAttachEventService = &SpanAttachedEventReportService{receiveChannel: r.OutputChannel}
	v3.RegisterTraceSegmentReportServiceServer(r.Server, r.traceService)
	v3.RegisterSpanAttachedEventReportServiceServer(r.Server, r.spanAttachEventService)
	v3_compat.RegisterTraceSegmentReportServiceServer(r.Server, &TraceSegmentReportServiceCompat{reportService: r.traceService})
}

func (r *Receiver) RegisterSyncInvoker(_ module.SyncInvoker) {
}

func (r *Receiver) Channel() <-chan *v1.SniffData {
	return r.OutputChannel
}

func (r *Receiver) SupportForwarders() []forwarder.Forwarder {
	return []forwarder.Forwarder{
		new(frowarder_nativetracing.Forwarder),
	}
}
