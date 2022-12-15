// SiGG-Satellite-Network-SII  //

package nativecds

import (
	"context"

	v3 "skywalking.apache.org/repo/goapi/collect/agent/configuration/v3"

	module "github.com/apache/skywalking-satellite/internal/satellite/module/api"

	sniffer "skywalking.apache.org/repo/goapi/satellite/data/v1"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
)

type CDSService struct {
	receiveChannel chan *sniffer.SniffData

	module.SyncInvoker
	v3.UnimplementedConfigurationDiscoveryServiceServer
}

func (p *CDSService) FetchConfigurations(_ context.Context, req *v3.ConfigurationSyncRequest) (*common.Commands, error) {
	event := &sniffer.SniffData{
		Data: &sniffer.SniffData_ConfigurationSyncRequest{
			ConfigurationSyncRequest: req,
		},
	}
	data, err := p.SyncInvoker.SyncInvoke(event)
	if err != nil {
		return nil, err
	}
	return data.GetCommands(), nil
}
