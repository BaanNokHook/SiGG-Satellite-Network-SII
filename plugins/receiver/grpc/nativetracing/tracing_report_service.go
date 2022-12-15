// SiGG-Satellite-Network-SII  //

package nativetracing

import (
	"context"
	"io"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/apache/skywalking-satellite/internal/pkg/log"
	"github.com/apache/skywalking-satellite/plugins/server/grpc"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	agent "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const eventName = "grpc-tracing-event"

type TraceSegmentReportService struct {
	receiveChannel chan *v1.SniffData
	agent.UnimplementedTraceSegmentReportServiceServer
}

func (s *TraceSegmentReportService) Collect(stream agent.TraceSegmentReportService_CollectServer) error {
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
			Data: &v1.SniffData_Segment{
				Segment: recData.Content,
			},
		}
		s.receiveChannel <- e
	}
}

func (s *TraceSegmentReportService) CollectInSync(ctx context.Context, segments *agent.SegmentCollection) (*common.Commands, error) {
	for _, segment := range segments.Segments {
		marshaledSegment, err := proto.Marshal(segment)
		if err != nil {
			log.Logger.Warnf("cannot marshal segemnt from sync, %v", err)
		}
		e := &v1.SniffData{
			Name:      eventName,
			Timestamp: time.Now().UnixNano() / 1e6,
			Meta:      nil,
			Type:      v1.SniffType_TracingType,
			Remote:    true,
			Data: &v1.SniffData_Segment{
				Segment: marshaledSegment,
			},
		}
		s.receiveChannel <- e
	}

	return &common.Commands{}, nil
}
