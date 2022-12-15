// SiGG-Satellite-Network-SII  //

package nativeevent

import (
	"context"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc"

	nativeevent "skywalking.apache.org/repo/goapi/collect/event/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	receiver_grpc "github.com/apache/skywalking-satellite/plugins/receiver/grpc"
)

func TestReceiver_RegisterHandler(t *testing.T) {
	receiver_grpc.TestReceiver(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, ctx context.Context) string {
		client := nativeevent.NewEventServiceClient(conn)
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
		return data.GetEvent().String()
	}, t)
}

func initData(sequence int) *nativeevent.Event {
	seq := strconv.Itoa(sequence)
	return &nativeevent.Event{
		StartTime: time.Now().Unix() / 1e6,
		EndTime:   time.Now().Unix() / 1e6,
		Uuid:      "12345" + seq,
		Source: &nativeevent.Source{
			Service:         "demo-service" + seq,
			ServiceInstance: "demo-instance" + seq,
			Endpoint:        "demo-endpoint" + seq,
		},
		Name:    "test-name" + seq,
		Type:    nativeevent.Type_Error,
		Message: "test message" + seq,
	}
}
