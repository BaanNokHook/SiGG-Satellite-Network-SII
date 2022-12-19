package test

import (
	"github.com/apache/skywalking-satellite/internal/pkg/log"
	"github.com/apache/skywalking-satellite/internal/satellite/telemetry"

	// import default telemetry
	_ "github.com/apache/skywalking-satellite/internal/satellite/telemetry/none"
)

// init the dependency components.
func init() {
	log.Init(new(log.LoggerConfig))
	c := new(telemetry.Config)
	c.ExportType = "none"
	if err := telemetry.Init(c); err != nil {
		panic(err)
	}
}
