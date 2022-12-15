// SiGG-Satellite-Network-SII  //

package nativetracing

import (
	"context"
	"strconv"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"google.golang.org/grpc"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	agent "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	receiver_grpc "github.com/apache/skywalking-satellite/plugins/receiver/grpc"
)

func TestReceiver_RegisterHandlerStream(t *testing.T) {
	receiver_grpc.TestReceiver(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, ctx context.Context) string {
		client := agent.NewTraceSegmentReportServiceClient(conn)
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
		d := new(agent.SegmentObject)
		_ = proto.Unmarshal(data.GetSegment(), d)
		return d.String()
	}, t)
}

func TestReceiver_RegisterHandlerSync(t *testing.T) {
	receiver_grpc.TestReceiver(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, ctx context.Context) string {
		client := agent.NewTraceSegmentReportServiceClient(conn)
		data := initData(sequence)
		segmentCollection := &agent.SegmentCollection{
			Segments: []*agent.SegmentObject{
				data,
			},
		}
		_, err := client.CollectInSync(ctx, segmentCollection)
		if err != nil {
			t.Fatalf("cannot send data in sync mode: %v", err)
		}
		return data.String()
	}, func(data *v1.SniffData) string {
		d := new(agent.SegmentObject)
		_ = proto.Unmarshal(data.GetSegment(), d)
		return d.String()
	}, t)
}

func initData(sequence int) *agent.SegmentObject {
	seq := strconv.Itoa(sequence)
	return &agent.SegmentObject{
		TraceId:         "trace" + seq,
		TraceSegmentId:  "trace-segment" + seq,
		Service:         "demo-service" + seq,
		ServiceInstance: "demo-instance" + seq,
		IsSizeLimited:   false,
		Spans: []*agent.SpanObject{
			{
				SpanId:        1,
				ParentSpanId:  0,
				StartTime:     time.Now().Unix(),
				EndTime:       time.Now().Unix() + 10,
				OperationName: "operation" + seq,
				Peer:          "127.0.0.1:6666",
				SpanType:      agent.SpanType_Entry,
				SpanLayer:     agent.SpanLayer_Http,
				ComponentId:   1,
				IsError:       false,
				SkipAnalysis:  false,
				Tags: []*common.KeyStringValuePair{
					{
						Key:   "mock-key" + seq,
						Value: "mock-value" + seq,
					},
				},
				Logs: []*agent.Log{
					{
						Time: time.Now().Unix(),
						Data: []*common.KeyStringValuePair{
							{
								Key:   "log-key" + seq,
								Value: "log-value" + seq,
							},
						},
					},
				},
				Refs: []*agent.SegmentReference{
					{
						RefType:                  agent.RefType_CrossThread,
						TraceId:                  "trace" + seq,
						ParentTraceSegmentId:     "parent-trace-segment" + seq,
						ParentSpanId:             0,
						ParentService:            "parent" + seq,
						ParentServiceInstance:    "parent" + seq,
						ParentEndpoint:           "parent" + seq,
						NetworkAddressUsedAtPeer: "127.0.0.1:6666",
					},
				},
			},
		},
	}
}
