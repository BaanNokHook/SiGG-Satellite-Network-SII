// SiGG-Satellite-Network-SII  //

package nativeclr

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	agent "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	receiver_grpc "github.com/apache/skywalking-satellite/plugins/receiver/grpc"
)

func TestReceiver_RegisterHandler(t *testing.T) {
	receiver_grpc.TestReceiver(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, ctx context.Context) string {
		client := agent.NewCLRMetricReportServiceClient(conn)
		data := initData()
		_, err := client.Collect(ctx, data)
		if err != nil {
			t.Fatalf("cannot send data: %v", err)
		}
		return data.String()
	}, func(data *v1.SniffData) string {
		return data.GetClr().String()
	}, t)
}

func initData() *agent.CLRMetricCollection {
	return &agent.CLRMetricCollection{
		Service:         "demo-service",
		ServiceInstance: "demo-instance",
		Metrics: []*agent.CLRMetric{
			{
				Time: time.Now().Unix() / 1e6,
				Cpu: &common.CPU{
					UsagePercent: 99.9,
				},
				Gc: &agent.ClrGC{
					Gen0CollectCount: 1,
					Gen1CollectCount: 2,
					Gen2CollectCount: 3,
					HeapMemory:       1024 * 1024 * 1024,
				},
				Thread: &agent.ClrThread{
					AvailableWorkerThreads:         10,
					AvailableCompletionPortThreads: 10,
					MaxWorkerThreads:               64,
					MaxCompletionPortThreads:       64,
				},
			},
		},
	}
}
