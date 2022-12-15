// SiGG-Satellite-Network-SII  //

package nativemeter

import (
	"context"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc"

	meter "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	receiver_grpc "github.com/apache/skywalking-satellite/plugins/receiver/grpc"
)

func TestReceiver_RegisterHandler(t *testing.T) {
	receiver_grpc.TestReceiver(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, ctx context.Context) string {
		client := meter.NewMeterReportServiceClient(conn)
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
		return data.GetMeterCollection().MeterData[0].String()
	}, t)
}

func initData(sequence int) *meter.MeterData {
	seq := strconv.Itoa(sequence)
	return &meter.MeterData{
		Service:         "demo-service" + seq,
		ServiceInstance: "demo-instance" + seq,
		Timestamp:       time.Now().Unix() / 1e6,
		Metric: &meter.MeterData_SingleValue{
			SingleValue: &meter.MeterSingleValue{
				Name:  "name" + seq,
				Value: float64(sequence),
				Labels: []*meter.Label{
					{
						Name:  "label-name" + seq,
						Value: "label-value" + seq,
					},
				},
			},
		},
	}
}
