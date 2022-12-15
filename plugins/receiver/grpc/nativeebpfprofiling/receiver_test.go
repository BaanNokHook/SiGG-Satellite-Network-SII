// SiGG-Satellite-Network-SII  //

package nativeebpfprofiling

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"

	agent "skywalking.apache.org/repo/goapi/collect/ebpf/profiling/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	receiver_grpc "github.com/apache/skywalking-satellite/plugins/receiver/grpc"
)

func TestReceiver_RegisterHandler(t *testing.T) {
	receiver_grpc.TestReceiver(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, ctx context.Context) string {
		client := agent.NewEBPFProfilingServiceClient(conn)
		data := initData()
		collect, err := client.CollectProfilingData(ctx)
		if err != nil {
			t.Fatalf("cannot send data: %v", err)
		}
		if err := collect.Send(data); err != nil {
			t.Fatalf("cannot send the data to the server: %v", err)
		}
		if err := collect.CloseSend(); err != nil {
			t.Fatalf("cannot close the stream: %v", err)
		}
		return data.String()
	}, func(data *v1.SniffData) string {
		return data.GetEBPFProfilingDataList().DataList[0].String()
	}, t)
}

func initData() *agent.EBPFProfilingData {
	return &agent.EBPFProfilingData{
		Task: &agent.EBPFProfilingTaskMetadata{
			TaskId:             "test-task",
			ProcessId:          "test-process",
			ProfilingStartTime: time.Now().Unix() / 1e6,
			CurrentTime:        time.Now().Unix() / 1e6,
		},
		Profiling: &agent.EBPFProfilingData_OnCPU{
			OnCPU: &agent.EBPFOnCPUProfiling{
				DumpCount: 1,
				Stacks: []*agent.EBPFProfilingStackMetadata{
					{
						StackType: agent.EBPFProfilingStackType_PROCESS_USER_SPACE,
						StackId:   123,
						StackSymbols: []string{
							"stack1", "stack2",
						},
					},
				},
			},
		},
	}
}
