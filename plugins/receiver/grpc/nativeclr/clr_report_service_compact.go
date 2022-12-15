// SiGG-Satellite-Network-SII  //

package nativeclr

import (
	"context"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	agent "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	agent_compat "skywalking.apache.org/repo/goapi/collect/language/agent/v3/compat"
)

type CLRReportServiceCompat struct {
	reportService *CLRReportService
	agent_compat.UnimplementedCLRMetricReportServiceServer
}

func (j *CLRReportServiceCompat) Collect(ctx context.Context, clr *agent.CLRMetricCollection) (*common.Commands, error) {
	return j.reportService.Collect(ctx, clr)
}
