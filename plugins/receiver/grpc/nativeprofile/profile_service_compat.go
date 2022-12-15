// SiGG-Satellite-Network-SII  //

package nativeprofile

import (
	"context"

	common "skywalking.apache.org/repo/goapi/collect/common/v3"
	profile "skywalking.apache.org/repo/goapi/collect/language/profile/v3"
	profile_compat "skywalking.apache.org/repo/goapi/collect/language/profile/v3/compat"
)

type ProfileServiceCompat struct {
	reportService *ProfileService
	profile_compat.UnimplementedProfileTaskServer
}

func (p *ProfileServiceCompat) GetProfileTaskCommands(ctx context.Context, q *profile.ProfileTaskCommandQuery) (*common.Commands, error) {
	return p.reportService.GetProfileTaskCommands(ctx, q)
}

func (p *ProfileServiceCompat) CollectSnapshot(stream profile_compat.ProfileTask_CollectSnapshotServer) error {
	return p.reportService.CollectSnapshot(stream)
}

func (p *ProfileServiceCompat) ReportTaskFinish(ctx context.Context, report *profile.ProfileTaskFinishReport) (*common.Commands, error) {
	return p.reportService.ReportTaskFinish(ctx, report)
}
