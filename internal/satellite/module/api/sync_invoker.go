package api

import v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"

type SyncInvoker interface {
	// SyncInvoke means synchronized process event
	SyncInvoke(d *v1.SniffData) (*v1.SniffData, error)
}
