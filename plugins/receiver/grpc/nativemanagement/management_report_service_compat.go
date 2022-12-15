// SiGG-Satellite-Network-SII  //

package nativemanagement

import (
	"context"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	management "skywalking.apache.org/repo/goapi/collect/management/v3"
	management_compat "skywalking.apache.org/repo/goapi/collect/management/v3/compat"
)

type ManagementReportServiceCompat struct {
	reportService *ManagementReportService
	management_compat.UnimplementedManagementServiceServer
}

func (m *ManagementReportServiceCompat) ReportInstanceProperties(ctx context.Context, in *management.InstanceProperties) (*common.Commands, error) {
	return m.reportService.ReportInstanceProperties(ctx, in)
}

func (m *ManagementReportServiceCompat) KeepAlive(ctx context.Context, in *management.InstancePingPkg) (*common.Commands, error) {
	return m.reportService.KeepAlive(ctx, in)
}
