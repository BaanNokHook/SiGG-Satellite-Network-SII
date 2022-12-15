// SiGG-Satellite-Network-SII  //

package otlpmetricsv1

import (
	"context"
	"testing"

	metrics "skywalking.apache.org/repo/goapi/proto/opentelemetry/proto/collector/metrics/v1"
	common "skywalking.apache.org/repo/goapi/proto/opentelemetry/proto/common/v1"
	v1 "skywalking.apache.org/repo/goapi/proto/opentelemetry/proto/metrics/v1"
	resource "skywalking.apache.org/repo/goapi/proto/opentelemetry/proto/resource/v1"
	sniffer "skywalking.apache.org/repo/goapi/satellite/data/v1"

	"google.golang.org/grpc"

	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	receiver_grpc "github.com/apache/skywalking-satellite/plugins/receiver/grpc"
)

func TestReceiver_RegisterHandler(t *testing.T) {
	recConf := make(map[string]string, 2)
	receiver_grpc.TestReceiverWithConfig(new(Receiver), recConf, func(t *testing.T, sequence int, conn *grpc.ClientConn, ctx context.Context) string {
		client := metrics.NewMetricsServiceClient(conn)
		data := initData()
		_, err := client.Export(ctx, data)
		if err != nil {
			t.Fatalf("cannot open the stream send mode: %v", err)
		}
		return data.String()
	}, func(data *sniffer.SniffData) string {
		return data.GetOpenTelementryMetricsV1Request().String()
	}, t)
}

func initData() *metrics.ExportMetricsServiceRequest {
	return &metrics.ExportMetricsServiceRequest{
		ResourceMetrics: []*v1.ResourceMetrics{
			{
				Resource: &resource.Resource{
					Attributes: []*common.KeyValue{
						{
							Key: "test",
							Value: &common.AnyValue{
								Value: &common.AnyValue_StringValue{
									StringValue: "1",
								},
							},
						},
					},
				},
			},
		},
	}
}
