package none

import (
	"time"

	"github.com/apache/skywalking-satellite/internal/satellite/telemetry"
)

func init() {
	telemetry.Register("none", &Server{}, true)
}

type Server struct {
}

func (s *Server) Start(config *telemetry.Config) error {
	return nil
}

func (s *Server) AfterSharingStart() error {
	return nil
}

func (s *Server) Close() error {
	return nil
}

type Counter struct {
}

func (c *Counter) Inc(labelValues ...string) {
}
func (c *Counter) Add(val float64, labelValues ...string) {
}

func (s *Server) NewCounter(name, help string, labels ...string) telemetry.Counter {
	return &Counter{}
}

type Gauge struct {
}

func (s *Server) NewGauge(name, help string, getter func() float64, labels ...string) telemetry.Gauge {
	return &Gauge{}
}

type DynamicGauge struct {
}

func (d *DynamicGauge) Inc(labelValues ...string) {
}
func (d *DynamicGauge) Dec(labelValues ...string) {
}

func (s *Server) NewDynamicGauge(name, help string, labels ...string) telemetry.DynamicGauge {
	return &DynamicGauge{}
}

type Timer struct {
}
type TimeRecorder struct {
}

func (t *Timer) Start(labelValues ...string) telemetry.TimeRecorder {
	return &TimeRecorder{}
}
func (t *Timer) AddTime(d time.Duration, labelValues ...string) {
}
func (r *TimeRecorder) Stop() {
}

func (s *Server) NewTimer(name, help string, labels ...string) telemetry.Timer {
	return &Timer{}
}
