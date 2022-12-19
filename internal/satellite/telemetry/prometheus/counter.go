package prometheus

import (
	"github.com/apache/skywalking-satellite/internal/satellite/telemetry"

	"github.com/prometheus/client_golang/prometheus"
)

// The counter metric.
type Counter struct {
	Collector
	name    string // The name of counter.
	counter *prometheus.CounterVec
}

// NewCounter create a new counter if no metric with the same name exists.
func (s *Server) NewCounter(name, help string, labels ...string) telemetry.Counter {
	s.lock.Lock()
	defer s.lock.Unlock()
	collector, ok := s.collectorContainer[name]
	if !ok {
		counter := &Counter{
			name: name,
			counter: prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: name,
				Help: help,
			}, labels),
		}
		s.Register(s.WithMeta(name, counter.counter))
		s.collectorContainer[name] = counter
		collector = counter
	}
	return collector.(telemetry.Counter)
}

// Add one.
func (c *Counter) Inc(labelValues ...string) {
	c.counter.WithLabelValues(labelValues...).Inc()
}

// Add float value.
func (c *Counter) Add(val float64, labelValues ...string) {
	c.counter.WithLabelValues(labelValues...).Add(val)
}
