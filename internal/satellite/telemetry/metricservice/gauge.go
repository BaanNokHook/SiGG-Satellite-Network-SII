package metricservice

import (
	"github.com/apache/skywalking-satellite/internal/satellite/telemetry"

	v3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

type Gauge struct {
	BaseMetric
	labels map[string]string
	getter func() float64
}

type DynamicGauge struct {
	BaseMetric
}

type subDynamicGauge struct {
	val float64
}

func (s *Server) NewGauge(name, _ string, getter func() float64, labels ...string) telemetry.Gauge {
	s.lock.Lock()
	defer s.lock.Unlock()
	rebuildName := rebuildGaugeName(name, labels...)
	constLabels := make(map[string]string)
	for inx := 0; inx < len(labels); inx += 2 {
		constLabels[labels[inx]] = labels[inx+1]
	}
	metric, ok := s.metrics[rebuildName]
	if !ok {
		metric = &Gauge{
			*NewBaseMetric(name, nil, func(labelValues ...string) SubMetric {
				return nil
			}),
			constLabels,
			getter,
		}
		s.Register(rebuildName, metric)
	}
	return metric
}

func (s *Server) NewDynamicGauge(name, _ string, labels ...string) telemetry.DynamicGauge {
	s.lock.Lock()
	defer s.lock.Unlock()
	rebuildName := rebuildGaugeName(name, labels...)
	metric, ok := s.metrics[rebuildName]
	if !ok {
		metric = &DynamicGauge{
			*NewBaseMetric(name, labels, func(labelValues ...string) SubMetric {
				return &subDynamicGauge{0}
			}),
		}
		s.Register(rebuildName, metric)
	}
	return metric.(telemetry.DynamicGauge)
}

func (g *Gauge) WriteMetric(appender *MetricsAppender) {
	labels := make([]*v3.Label, 0)
	for k, v := range g.labels {
		labels = append(labels, &v3.Label{
			Name:  k,
			Value: v,
		})
	}
	appender.appendSingleValue(g.Name, labels, g.getter())
}

func rebuildGaugeName(name string, labels ...string) string {
	resultName := name
	for inx := 0; inx < len(labels); inx++ {
		resultName += "_" + labels[inx]
	}

	return resultName
}

func (d *DynamicGauge) Inc(labels ...string) {
	if counter, err := d.GetMetricWithLabelValues(labels...); err != nil {
		panic(err)
	} else {
		addFloat64(&(counter.(*subDynamicGauge).val), 1)
	}
}

func (d *DynamicGauge) Dec(labels ...string) {
	if counter, err := d.GetMetricWithLabelValues(labels...); err != nil {
		panic(err)
	} else {
		addFloat64(&(counter.(*subDynamicGauge).val), -1)
	}
}

func (c *subDynamicGauge) WriteMetric(base *BaseMetric, labels []*v3.Label, appender *MetricsAppender) {
	appender.appendSingleValue(base.Name, labels, c.val)
}
