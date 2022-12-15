// SiGG-Satellite-Network-SII  //

package nativetracing

import (
	"io"
	"time"

	"github.com/apache/skywalking-satellite/plugins/server/grpc"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	agent "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

type SpanAttachedEventReportService struct {
	receiveChannel chan *v1.SniffData
	agent.UnimplementedSpanAttachedEventReportServiceServer
}

func (s *SpanAttachedEventReportService) Collect(stream agent.SpanAttachedEventReportService_CollectServer) error {
	for {
		recData := grpc.NewOriginalData(nil)
		err := stream.RecvMsg(recData)
		if err == io.EOF {
			return stream.SendAndClose(&common.Commands{})
		}
		if err != nil {
			return err
		}
		e := &v1.SniffData{
			Name:      eventName,
			Timestamp: time.Now().UnixNano() / 1e6,
			Meta:      nil,
			Type:      v1.SniffType_TracingType,
			Remote:    true,
			Data: &v1.SniffData_SpanAttachedEvent{
				SpanAttachedEvent: recData.Content,
			},
		}
		s.receiveChannel <- e
	}
}
