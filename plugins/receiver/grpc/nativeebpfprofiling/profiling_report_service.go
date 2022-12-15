// SiGG-Satellite-Network-SII  //

package nativeebpfprofiling

import (
	"context"
	"io"
	"time"

	module "github.com/apache/skywalking-satellite/internal/satellite/module/api"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	v3 "skywalking.apache.org/repo/goapi/collect/ebpf/profiling/v3"
	sniffer "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const eventName = "grpc-ebpfebpfprofiling-event"

type ProfilingReportService struct {
	receiveChannel chan *sniffer.SniffData

	module.SyncInvoker
	v3.UnimplementedEBPFProfilingServiceServer
}

func (p *ProfilingReportService) QueryTasks(ctx context.Context, d *v3.EBPFProfilingTaskQuery) (*common.Commands, error) {
	event := &sniffer.SniffData{
		Data: &sniffer.SniffData_EBPFProfilingTaskQuery{
			EBPFProfilingTaskQuery: d,
		},
	}
	data, err := p.SyncInvoker.SyncInvoke(event)
	if err != nil {
		return nil, err
	}
	return data.GetCommands(), nil
}

func (p *ProfilingReportService) CollectProfilingData(stream v3.EBPFProfilingService_CollectProfilingDataServer) error {
	dataList := make([]*v3.EBPFProfilingData, 0)
	for {
		snapshot, err := stream.Recv()
		if err == io.EOF {
			return p.sendData(dataList, stream)
		}
		if err != nil {
			return err
		}
		dataList = append(dataList, snapshot)
	}
}

func (p *ProfilingReportService) sendData(dataList []*v3.EBPFProfilingData, stream v3.EBPFProfilingService_CollectProfilingDataServer) error {
	e := &sniffer.SniffData{
		Name:      eventName,
		Timestamp: time.Now().UnixNano() / 1e6,
		Meta:      nil,
		Type:      sniffer.SniffType_EBPFProfilingType,
		Remote:    true,
		Data: &sniffer.SniffData_EBPFProfilingDataList{
			EBPFProfilingDataList: &sniffer.EBPFProfilingDataList{
				DataList: dataList,
			},
		},
	}
	p.receiveChannel <- e
	return stream.SendAndClose(&common.Commands{})
}
