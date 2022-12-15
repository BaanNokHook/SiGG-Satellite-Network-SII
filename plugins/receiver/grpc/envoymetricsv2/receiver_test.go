// SiGG-Satellite-Network-SII  //

package envoymetricsv2

import (
	"context"
	"strconv"
	"testing"

	"skywalking.apache.org/repo/goapi/proto/envoy/api/v2/core"

	v2 "skywalking.apache.org/repo/goapi/proto/envoy/service/metrics/v2"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	"google.golang.org/grpc"

	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	receiver_grpc "github.com/apache/skywalking-satellite/plugins/receiver/grpc"
)

func TestReceiver_RegisterHandler(t *testing.T) {
	recConf := make(map[string]string, 2)
	recConf["limit_count"] = "1"
	recConf["flush_time"] = "1000"
	receiver_grpc.TestReceiverWithConfig(new(Receiver), recConf, func(t *testing.T, sequence int, conn *grpc.ClientConn, ctx context.Context) string {
		client := v2.NewMetricsServiceClient(conn)
		data := initData(sequence)
		collect, err := client.StreamMetrics(ctx)
		if err != nil {
			t.Fatalf("cannot open the stream send mode: %v", err)
		}
		if err := collect.Send(data); err != nil {
			t.Fatalf("cannot send the data to the server: %v", err)
		}
		if err := collect.CloseSend(); err != nil {
			t.Fatalf("cannot close the stream mode: %v", err)
		}
		return data.String()
	}, func(data *v1.SniffData) string {
		return data.GetEnvoyMetricsV2List().Messages[0].String()
	}, t)
}

func initData(sequence int) *v2.StreamMetricsMessage {
	return &v2.StreamMetricsMessage{
		Identifier: &v2.StreamMetricsMessage_Identifier{
			Node: &core.Node{
				Id: "test" + strconv.Itoa(sequence),
			},
		},
	}
}
