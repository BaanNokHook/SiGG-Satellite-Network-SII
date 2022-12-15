// SiGG-Satellite-Network-SII  //

package nativelog

import (
	"context"
	"strconv"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"google.golang.org/grpc"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	logging "skywalking.apache.org/repo/goapi/collect/logging/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	receiver_grpc "github.com/apache/skywalking-satellite/plugins/receiver/grpc"
)

func TestReceiver_RegisterHandler(t *testing.T) {
	receiver_grpc.TestReceiver(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, ctx context.Context) string {
		client := logging.NewLogReportServiceClient(conn)
		data := initData(sequence)
		collect, err := client.Collect(ctx)
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
		d := new(logging.LogData)
		_ = proto.Unmarshal(data.GetLogList().Logs[0], d)
		return d.String()
	}, t)
}

func initData(sequence int) *logging.LogData {
	seq := strconv.Itoa(sequence)
	return &logging.LogData{
		Timestamp:       time.Now().Unix(),
		Service:         "demo-service" + seq,
		ServiceInstance: "demo-instance" + seq,
		Endpoint:        "demo-endpoint" + seq,
		TraceContext: &logging.TraceContext{
			TraceSegmentId: "mock-segmentId" + seq,
			TraceId:        "mock-traceId" + seq,
			SpanId:         1,
		},
		Tags: &logging.LogTags{
			Data: []*common.KeyStringValuePair{
				{
					Key:   "mock-key" + seq,
					Value: "mock-value" + seq,
				},
			},
		},
		Body: &logging.LogDataBody{
			Type: "mock-type" + seq,
			Content: &logging.LogDataBody_Text{
				Text: &logging.TextLog{
					Text: "this is a mock text mock log" + seq,
				},
			},
		},
	}
}
