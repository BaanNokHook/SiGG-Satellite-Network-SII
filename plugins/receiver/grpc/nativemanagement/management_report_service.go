// SiGG-Satellite-Network-SII  //

package nativemanagement

import (
	"context"
	"time"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	management "skywalking.apache.org/repo/goapi/collect/management/v3"
	sniffer "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const eventName = "grpc-management-event"

type ManagementReportService struct {
	receiveChannel chan *sniffer.SniffData
	management.UnimplementedManagementServiceServer
}

func (m *ManagementReportService) ReportInstanceProperties(ctx context.Context, in *management.InstanceProperties) (*common.Commands, error) {
	e := &sniffer.SniffData{
		Name:      eventName,
		Timestamp: time.Now().UnixNano() / 1e6,
		Meta:      nil,
		Type:      sniffer.SniffType_ManagementType,
		Remote:    true,
		Data: &sniffer.SniffData_Instance{
			Instance: in,
		},
	}

	m.receiveChannel <- e
	return &common.Commands{}, nil
}

func (m *ManagementReportService) KeepAlive(ctx context.Context, in *management.InstancePingPkg) (*common.Commands, error) {
	e := &sniffer.SniffData{
		Name:      eventName,
		Timestamp: time.Now().UnixNano() / 1e6,
		Meta:      nil,
		Type:      sniffer.SniffType_ManagementType,
		Remote:    true,
		Data: &sniffer.SniffData_InstancePing{
			InstancePing: in,
		},
	}

	m.receiveChannel <- e
	return &common.Commands{}, nil
}
