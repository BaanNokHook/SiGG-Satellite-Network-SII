// SiGG-Satellite-Network-SII  //

package nativeprocess

import (
	"context"
	"time"

	module "github.com/apache/skywalking-satellite/internal/satellite/module/api"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	v3 "skywalking.apache.org/repo/goapi/collect/ebpf/profiling/process/v3"
	sniffer "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const eventName = "grpc-ebpfprocess-event"

type ProcessReportService struct {
	receiveChannel chan *sniffer.SniffData

	module.SyncInvoker
	v3.UnimplementedEBPFProcessServiceServer
}

func (p *ProcessReportService) ReportProcesses(ctx context.Context, d *v3.EBPFProcessReportList) (*v3.EBPFReportProcessDownstream, error) {
	event := &sniffer.SniffData{
		Data: &sniffer.SniffData_EBPFProcessReportList{
			EBPFProcessReportList: d,
		},
	}
	data, err := p.SyncInvoker.SyncInvoke(event)
	if err != nil {
		return nil, err
	}
	return data.GetEBPFReportProcessDownstream(), nil
}

func (p *ProcessReportService) KeepAlive(ctx context.Context, d *v3.EBPFProcessPingPkgList) (*common.Commands, error) {
	e := &sniffer.SniffData{
		Name:      eventName,
		Timestamp: time.Now().UnixNano() / 1e6,
		Meta:      nil,
		Type:      sniffer.SniffType_EBPFProcessType,
		Remote:    true,
		Data: &sniffer.SniffData_EBPFProcessPingPkgList{
			EBPFProcessPingPkgList: d,
		},
	}

	p.receiveChannel <- e
	return &common.Commands{}, nil
}
