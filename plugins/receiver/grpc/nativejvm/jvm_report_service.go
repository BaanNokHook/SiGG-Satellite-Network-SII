// SiGG-Satellite-Network-SII  //

package nativelog

import (
	"context"
	"time"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	agent "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const eventName = "grpc-jvm-event"

type JVMReportService struct {
	receiveChannel chan *v1.SniffData
	agent.UnimplementedJVMMetricReportServiceServer
}

func (j *JVMReportService) Collect(_ context.Context, jvm *agent.JVMMetricCollection) (*common.Commands, error) {
	e := &v1.SniffData{
		Name:      eventName,
		Timestamp: time.Now().UnixNano() / 1e6,
		Meta:      nil,
		Type:      v1.SniffType_JVMMetricType,
		Remote:    true,
		Data: &v1.SniffData_Jvm{
			Jvm: jvm,
		},
	}
	j.receiveChannel <- e
	return &common.Commands{}, nil
}
