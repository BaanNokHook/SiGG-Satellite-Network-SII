// SiGG-Satellite-Network-SII  //

package nativeclr

import (
	"context"
	"time"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	agent "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const eventName = "grpc-clr-event"

type CLRReportService struct {
	receiveChannel chan *v1.SniffData
	agent.UnimplementedCLRMetricReportServiceServer
}

func (j *CLRReportService) Collect(_ context.Context, clr *agent.CLRMetricCollection) (*common.Commands, error) {
	e := &v1.SniffData{
		Name:      eventName,
		Timestamp: time.Now().UnixNano() / 1e6,
		Meta:      nil,
		Type:      v1.SniffType_CLRMetricType,
		Remote:    true,
		Data: &v1.SniffData_Clr{
			Clr: clr,
		},
	}
	j.receiveChannel <- e
	return &common.Commands{}, nil
}
