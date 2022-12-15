// SiGG-Satellite-Network-SII  //

package nativeprofile

import (
	"context"
	"fmt"
	"testing"
	"time"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"

	"google.golang.org/grpc"

	profile "skywalking.apache.org/repo/goapi/collect/language/profile/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	receiver_grpc "github.com/apache/skywalking-satellite/plugins/receiver/grpc"
)

func TestReceiver_RegisterHandler_ThreadSnapshot(t *testing.T) {
	receiver_grpc.TestReceiver(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, ctx context.Context) string {
		client := profile.NewProfileTaskClient(conn)
		data := &profile.ThreadSnapshot{
			TaskId:         fmt.Sprintf("task-%d", sequence),
			TraceSegmentId: fmt.Sprintf("segment-%d", sequence),
			Time:           time.Now().Unix(),
			Sequence:       int32(sequence),
			Stack: &profile.ThreadStack{
				CodeSignatures: []string{
					"code",
				},
			},
		}
		collect, err := client.CollectSnapshot(ctx)
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
		return data.GetProfile().String()
	}, t)
}

func TestReceiver_RegisterHandler_TaskCommandQuery(t *testing.T) {
	receiver_grpc.TestReceiverWithSync(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, sendData *string, ctx context.Context) {
		client := profile.NewProfileTaskClient(conn)
		data := &profile.ProfileTaskCommandQuery{
			Service:         fmt.Sprintf("service-%d", sequence),
			ServiceInstance: fmt.Sprintf("instance-%d", sequence),
			LastCommandTime: time.Now().Unix(),
		}
		*sendData = data.String()
		_, err := client.GetProfileTaskCommands(ctx, data)
		if err != nil {
			t.Fatalf("cannot send data: %v", err)
		}
	}, func(data *v1.SniffData) string {
		return data.GetProfileTaskQuery().String()
	}, &v1.SniffData{
		Data: &v1.SniffData_Commands{
			Commands: &common.Commands{},
		},
	}, t)
}

func TestReceiver_RegisterHandler_ReportTaskFinish(t *testing.T) {
	receiver_grpc.TestReceiverWithSync(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, sendData *string, ctx context.Context) {
		client := profile.NewProfileTaskClient(conn)
		data := &profile.ProfileTaskFinishReport{
			Service:         fmt.Sprintf("service-%d", sequence),
			ServiceInstance: fmt.Sprintf("instance-%d", sequence),
			TaskId:          fmt.Sprintf("task-%d", sequence),
		}
		*sendData = data.String()
		_, err := client.ReportTaskFinish(ctx, data)
		if err != nil {
			t.Fatalf("cannot send data: %v", err)
		}
	}, func(data *v1.SniffData) string {
		return data.GetProfileTaskFinish().String()
	}, &v1.SniffData{
		Data: &v1.SniffData_Commands{
			Commands: &common.Commands{},
		},
	}, t)
}
