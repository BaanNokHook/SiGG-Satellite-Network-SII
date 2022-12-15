// SiGG-Satellite-Network-SII  //

package otlpmetricsv1

import (
	"context"
	"time"

	metrics "skywalking.apache.org/repo/goapi/proto/opentelemetry/proto/collector/metrics/v1"
	sniffer "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const eventName = "grpc-envoy-metrics-v3-event"

type MetricsService struct {
	receiveChannel chan *sniffer.SniffData
	metrics.MetricsServiceServer
}

func (m *MetricsService) Export(ctx context.Context, req *metrics.ExportMetricsServiceRequest) (*metrics.ExportMetricsServiceResponse, error) {
	e := &sniffer.SniffData{
		Name:      eventName,
		Timestamp: time.Now().UnixNano() / 1e6,
		Meta:      nil,
		Type:      sniffer.SniffType_OpenTelementryMetricsV1Type,
		Remote:    true,
		Data: &sniffer.SniffData_OpenTelementryMetricsV1Request{
			OpenTelementryMetricsV1Request: req,
		},
	}

	m.receiveChannel <- e
	return &metrics.ExportMetricsServiceResponse{}, nil
}
