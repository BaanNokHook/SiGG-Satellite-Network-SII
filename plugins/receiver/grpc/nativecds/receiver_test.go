// SiGG-Satellite-Network-SII  //

package nativecds

import (
	"context"
	"fmt"
	"testing"

	v3 "skywalking.apache.org/repo/goapi/collect/agent/configuration/v3"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"

	"google.golang.org/grpc"

	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	receiver_grpc "github.com/apache/skywalking-satellite/plugins/receiver/grpc"
)

func TestReceiver_RegisterHandler(t *testing.T) {
	receiver_grpc.TestReceiverWithSync(new(Receiver), func(t *testing.T, sequence int, conn *grpc.ClientConn, sendData *string, ctx context.Context) {
		client := v3.NewConfigurationDiscoveryServiceClient(conn)
		data := &v3.ConfigurationSyncRequest{
			Service: fmt.Sprintf("service-%d", sequence),
			Uuid:    "",
		}
		*sendData = data.String()
		_, err := client.FetchConfigurations(ctx, data)
		if err != nil {
			t.Fatalf("cannot send data: %v", err)
		}
	}, func(data *v1.SniffData) string {
		return data.GetConfigurationSyncRequest().String()
	}, &v1.SniffData{
		Data: &v1.SniffData_Commands{
			Commands: &common.Commands{},
		},
	}, t)
}
