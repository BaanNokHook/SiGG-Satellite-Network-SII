package prometheus

import "github.com/apache/skywalking-satellite/internal/satellite/telemetry"

// The Self-telemetry data collection interface.
type Collector interface {
	telemetry.Metric
}
