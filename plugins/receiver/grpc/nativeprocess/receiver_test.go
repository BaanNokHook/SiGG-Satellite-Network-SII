// SiGG-Satellite-Network-SII  //

package nativeprocess

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	agent "skywalking.apache.org/repo/goapi/collect/ebpf/profiling/process/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	receiver_grpc "github.com/apache/skywalking-satellite/plugins/receiver/grpc"
)

func TestReceiver_RegisterHandler(t *testing.T) {
	receiver_grpc.TestReceiver(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, ctx context.Context) string {
		client := agent.NewEBPFProcessServiceClient(conn)
		data := initData()
		_, err := client.KeepAlive(ctx, data)
		if err != nil {
			t.Fatalf("cannot send data: %v", err)
		}
		return data.String()
	}, func(data *v1.SniffData) string {
		return data.GetEBPFProcessPingPkgList().String()
	}, t)
}

func initData() *agent.EBPFProcessPingPkgList {
	return &agent.EBPFProcessPingPkgList{
		Processes: []*agent.EBPFProcessPingPkg{
			{
				EntityMetadata: &agent.EBPFProcessEntityMetadata{
					Layer:        "GENERAL",
					ServiceName:  "test-service",
					InstanceName: "test-instance",
					ProcessName:  "test-process",
					Labels: []string{
						"test-label",
					},
				},
			},
		},
	}
}
