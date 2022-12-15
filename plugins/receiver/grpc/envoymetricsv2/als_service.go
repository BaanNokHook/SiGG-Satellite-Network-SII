// SiGG-Satellite-Network-SII  //

package envoymetricsv2

import (
	"context"
	"io"
	"time"

	"github.com/apache/skywalking-satellite/internal/satellite/module/buffer"

	v2 "skywalking.apache.org/repo/goapi/proto/envoy/service/metrics/v2"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const eventName = "grpc-envoy-metrics-v2-event"

type MetricsService struct {
	receiveChannel chan *v1.SniffData
	limiterConfig  buffer.LimiterConfig
	v2.UnimplementedMetricsServiceServer
}

func (m *MetricsService) StreamMetrics(stream v2.MetricsService_StreamMetricsServer) error {
	messages := make(chan *v2.StreamMetricsMessage, m.limiterConfig.LimitCount*2)
	limiter := buffer.NewLimiter(m.limiterConfig, func() int {
		return len(messages)
	})

	var identity *v2.StreamMetricsMessage_Identifier

	defer limiter.Stop()
	limiter.Start(context.Background(), func() {
		count := len(messages)
		if count == 0 {
			return
		}
		metricsMessages := make([]*v2.StreamMetricsMessage, 0)
		for i := 0; i < count; i++ {
			metricsMessages = append(metricsMessages, <-messages)
		}
		metricsMessages[0].Identifier = identity

		d := &v1.SniffData{
			Name:      eventName,
			Timestamp: time.Now().UnixNano() / 1e6,
			Meta:      nil,
			Type:      v1.SniffType_EnvoyMetricsV2Type,
			Remote:    true,
			Data: &v1.SniffData_EnvoyMetricsV2List{
				EnvoyMetricsV2List: &v1.EnvoyMetricsV2List{
					Messages: metricsMessages,
				},
			},
		}
		m.receiveChannel <- d
	})

	var err1 error
	for {
		item, err := stream.Recv()
		if err != nil {
			err1 = err
			break
		}
		if item.Identifier != nil {
			identity = item.Identifier
		}
		messages <- item
		limiter.Check()
	}

	if err1 != io.EOF {
		return err1
	}
	return stream.SendAndClose(&v2.StreamMetricsResponse{})
}
