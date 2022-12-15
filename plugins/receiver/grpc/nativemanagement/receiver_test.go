// SiGG-Satellite-Network-SII  //

package nativemanagement

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	management "skywalking.apache.org/repo/goapi/collect/management/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	receiver_grpc "github.com/apache/skywalking-satellite/plugins/receiver/grpc"
)

func TestReceiver_RegisterHandler_ReportInstance(t *testing.T) {
	receiver_grpc.TestReceiver(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, ctx context.Context) string {
		client := management.NewManagementServiceClient(conn)
		properties := &management.InstanceProperties{
			Service:         fmt.Sprintf("service_%d", sequence),
			ServiceInstance: fmt.Sprintf("instance_%d", sequence),
			Properties:      []*common.KeyStringValuePair{},
		}
		commands, err := client.ReportInstanceProperties(ctx, properties)
		if err != nil {
			t.Fatalf("cannot send the data to the server: %v", err)
		}
		if commands == nil {
			t.Fatalf("report instance result is nil")
		}
		return properties.String()
	}, func(data *v1.SniffData) string {
		return data.GetInstance().String()
	}, t)
}

func TestReceiver_RegisterHandler_InstancePing(t *testing.T) {
	receiver_grpc.TestReceiver(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, ctx context.Context) string {
		client := management.NewManagementServiceClient(conn)
		instancePing := &management.InstancePingPkg{
			Service:         fmt.Sprintf("service_%d", sequence),
			ServiceInstance: fmt.Sprintf("instance_%d", sequence),
		}
		commands, err := client.KeepAlive(ctx, instancePing)
		if err != nil {
			t.Fatalf("cannot send the data to the server: %v", err)
		}
		if commands == nil {
			t.Fatalf("instance ping result is nil")
		}
		return instancePing.String()
	}, func(data *v1.SniffData) string {
		return data.GetInstancePing().String()
	}, t)
}
