// SiGG-Satellite-Network-SII  //

package nativeprofile

import (
	"github.com/apache/skywalking-satellite/internal/pkg/config"
	module "github.com/apache/skywalking-satellite/internal/satellite/module/api"
	forwarder "github.com/apache/skywalking-satellite/plugins/forwarder/api"
	frowarder_nativeprofile "github.com/apache/skywalking-satellite/plugins/forwarder/grpc/nativeprofile"
	"github.com/apache/skywalking-satellite/plugins/receiver/grpc"

	v3 "skywalking.apache.org/repo/goapi/collect/language/profile/v3"
	v3_compat "skywalking.apache.org/repo/goapi/collect/language/profile/v3/compat"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const (
	Name     = "grpc-native-profile-receiver"
	ShowName = "GRPC Native Profile Receiver"
)

type Receiver struct {
	config.CommonFields
	grpc.CommonGRPCReceiverFields
	service *ProfileService
}

func (r *Receiver) Name() string {
	return Name
}

func (r *Receiver) ShowName() string {
	return ShowName
}

func (r *Receiver) Description() string {
	return "This is a receiver for SkyWalking native profile format, " +
		"which is defined at https://github.com/apache/skywalking-data-collect-protocol/blob/master/profile/Profile.proto."
}

func (r *Receiver) DefaultConfig() string {
	return ""
}

func (r *Receiver) RegisterHandler(server interface{}) {
	r.CommonGRPCReceiverFields = *grpc.InitCommonGRPCReceiverFields(server)
	r.service = &ProfileService{receiveChannel: r.OutputChannel}
	v3.RegisterProfileTaskServer(r.Server, r.service)
	v3_compat.RegisterProfileTaskServer(r.Server, &ProfileServiceCompat{reportService: r.service})
}

func (r *Receiver) RegisterSyncInvoker(invoker module.SyncInvoker) {
	r.service.SyncInvoker = invoker
}

func (r *Receiver) Channel() <-chan *v1.SniffData {
	return r.OutputChannel
}

func (r *Receiver) SupportForwarders() []forwarder.Forwarder {
	return []forwarder.Forwarder{
		new(frowarder_nativeprofile.Forwarder),
	}
}
