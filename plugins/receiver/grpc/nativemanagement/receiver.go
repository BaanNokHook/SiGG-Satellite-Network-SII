// SiGG-Satellite-Network-SII  //

package nativemanagement

import (
	"github.com/apache/skywalking-satellite/internal/pkg/config"
	module "github.com/apache/skywalking-satellite/internal/satellite/module/api"
	forwarder "github.com/apache/skywalking-satellite/plugins/forwarder/api"
	forwarder_nativemanagement "github.com/apache/skywalking-satellite/plugins/forwarder/grpc/nativemanagement"
	"github.com/apache/skywalking-satellite/plugins/receiver/grpc"

	v3 "skywalking.apache.org/repo/goapi/collect/management/v3"
	v3_compat "skywalking.apache.org/repo/goapi/collect/management/v3/compat"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const (
	Name     = "grpc-native-management-receiver"
	ShowName = "GRPC Native Management Receiver"
)

type Receiver struct {
	config.CommonFields
	grpc.CommonGRPCReceiverFields
	service *ManagementReportService
}

func (r *Receiver) Name() string {
	return Name
}

func (r *Receiver) ShowName() string {
	return ShowName
}

func (r *Receiver) Description() string {
	return "This is a receiver for SkyWalking native management format, " +
		"which is defined at https://github.com/apache/skywalking-data-collect-protocol/blob/master/management/Management.proto."
}

func (r *Receiver) DefaultConfig() string {
	return ""
}

func (r *Receiver) RegisterHandler(server interface{}) {
	r.CommonGRPCReceiverFields = *grpc.InitCommonGRPCReceiverFields(server)
	r.service = &ManagementReportService{receiveChannel: r.OutputChannel}
	v3.RegisterManagementServiceServer(r.Server, r.service)
	v3_compat.RegisterManagementServiceServer(r.Server, &ManagementReportServiceCompat{reportService: r.service})
}

func (r *Receiver) RegisterSyncInvoker(_ module.SyncInvoker) {
}

func (r *Receiver) Channel() <-chan *v1.SniffData {
	return r.OutputChannel
}

func (r *Receiver) SupportForwarders() []forwarder.Forwarder {
	return []forwarder.Forwarder{
		new(forwarder_nativemanagement.Forwarder),
	}
}
