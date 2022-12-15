// SiGG-Satellite-Network-SII  //

package nativetracing

import (
	"context"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	agent "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	agent_compat "skywalking.apache.org/repo/goapi/collect/language/agent/v3/compat"
)

type TraceSegmentReportServiceCompat struct {
	reportService *TraceSegmentReportService
	agent_compat.UnimplementedTraceSegmentReportServiceServer
}

func (s *TraceSegmentReportServiceCompat) Collect(stream agent_compat.TraceSegmentReportService_CollectServer) error {
	return s.reportService.Collect(stream)
}

func (s *TraceSegmentReportServiceCompat) CollectInSync(ctx context.Context, segments *agent.SegmentCollection) (*common.Commands, error) {
	return s.reportService.CollectInSync(ctx, segments)
}
