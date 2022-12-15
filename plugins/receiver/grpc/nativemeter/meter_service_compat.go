// SiGG-Satellite-Network-SII  //

package nativemeter

import (
	agent_compat "skywalking.apache.org/repo/goapi/collect/language/agent/v3/compat"
)

type MeterServiceCompat struct {
	reportService *MeterService
	agent_compat.UnimplementedMeterReportServiceServer
}

func (m *MeterServiceCompat) Collect(stream agent_compat.MeterReportService_CollectServer) error {
	return m.reportService.Collect(stream)
}
