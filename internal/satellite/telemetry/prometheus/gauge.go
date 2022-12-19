package prometheus

import (
	"github.com/apache/skywalking-satellite/internal/satellite/telemetry"

	"github.com/prometheus/client_golang/prometheus"
)

type Gauge struct {
	Collector
	name  string
	gauge prometheus.GaugeFunc
}

type DynamicGauge struct {
	Collector
	name  string
	gauge *prometheus.GaugeVec
}

func (s *Server) NewGauge(name, help string, getter func() float64, labels ...string) telemetry.Gauge {
	s.lock.Lock()
	defer s.lock.Unlock()
	rebuildName := rebuildGaugeName(name, labels...)
	constLabels := make(map[string]string)
	for inx := 0; inx < len(labels); inx += 2 {
		constLabels[labels[inx]] = labels[inx+1]
	}
	collector, ok := s.collectorContainer[rebuildName]
	if !ok {
		gauge := &Gauge{
			name: name,
			gauge: prometheus.NewGaugeFunc(prometheus.GaugeOpts{
				Name:        name,
				Help:        help,
				ConstLabels: constLabels,
			}, getter),
		}
		s.Register(s.WithMeta(rebuildName, gauge.gauge))
		s.collectorContainer[rebuildName] = gauge
		collector = gauge
	}
	return collector
}

func (s *Server) NewDynamicGauge(name, help string, labels ...string) telemetry.DynamicGauge {
	s.lock.Lock()
	defer s.lock.Unlock()
	rebuildName := rebuildGaugeName(name, labels...)
	collector, ok := s.collectorContainer[rebuildName]
	if !ok {
		gauge := &DynamicGauge{
			name: name,
			gauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Name: name,
				Help: help,
			}, labels),
		}
		s.Register(s.WithMeta(rebuildName, gauge.gauge))
		s.collectorContainer[rebuildName] = gauge
		collector = gauge
	}
	return collector.(telemetry.DynamicGauge)
}

func (i *DynamicGauge) Inc(labels ...string) {
	i.gauge.WithLabelValues(labels...).Inc()
}

func (i *DynamicGauge) Dec(labels ...string) {
	i.gauge.WithLabelValues(labels...).Dec()
}

func rebuildGaugeName(name string, labels ...string) string {
	resultName := name
	for inx := 0; inx < len(labels); inx++ {
		resultName += "_" + labels[inx]
	}

	return resultName
}
