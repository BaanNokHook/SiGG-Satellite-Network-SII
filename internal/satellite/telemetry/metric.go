package telemetry

import "time"

type Metric interface {
}

type Counter interface {
	Metric
	Inc(labelValues ...string)
	Add(val float64, labelValues ...string)
}

type Gauge interface {
	Metric
}

type DynamicGauge interface {
	Metric
	Inc(labelValues ...string)
	Dec(labelValues ...string)
}

type Timer interface {
	Metric
	Start(labelValues ...string) TimeRecorder
	AddTime(t time.Duration, labelValues ...string)
}

type TimeRecorder interface {
	Stop()
}
