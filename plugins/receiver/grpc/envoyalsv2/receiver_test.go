// SiGG-Satellite-Network-SII  //

package envoyalsv2

import (
	"context"
	"strconv"
	"testing"

	"google.golang.org/protobuf/proto"

	v2 "skywalking.apache.org/repo/goapi/proto/envoy/service/accesslog/v2"
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
		client := v2.NewAccessLogServiceClient(conn)
		data := initData(sequence)
		collect, err := client.StreamAccessLogs(ctx)
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
		m := new(v2.StreamAccessLogsMessage)
		_ = proto.Unmarshal(data.GetEnvoyALSV2List().Messages[0], m)
		return m.String()
	}, t)
}

func initData(sequence int) *v2.StreamAccessLogsMessage {
	return &v2.StreamAccessLogsMessage{
		Identifier: &v2.StreamAccessLogsMessage_Identifier{
			LogName: "test" + strconv.Itoa(sequence),
		},
	}
}
