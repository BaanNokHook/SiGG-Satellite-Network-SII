// SiGG-Satellite-Network-SII  //

package nativelog

import (
	"context"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	agent "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	agent_compat "skywalking.apache.org/repo/goapi/collect/language/agent/v3/compat"
)

type JVMReportServiceCompat struct {
	reportService *JVMReportService
	agent_compat.UnimplementedJVMMetricReportServiceServer
}

func (j *JVMReportServiceCompat) Collect(ctx context.Context, jvm *agent.JVMMetricCollection) (*common.Commands, error) {
	return j.reportService.Collect(ctx, jvm)
}
