package prometheus

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/apache/skywalking-satellite/internal/satellite/telemetry"
	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
)

func initPrometheusServer() (*telemetry.PrometheusConfig, error) {
	c := new(telemetry.Config)
	c.ExportType = "prometheus"
	c.Prometheus.Endpoint = "/metrics"
	c.Prometheus.Address = ":1234"
	if err := telemetry.Init(c); err != nil {
		return nil, fmt.Errorf("telemetry cannot initialize: %v", err)
	}
	return &c.Prometheus, nil
}

func TestServer_Start(t *testing.T) {
	server, err := initPrometheusServer()
	if err != nil {
		t.Fatalf("cannot init the prometheus server: %v", err)
	}
	time.Sleep(time.Second)
	response, err := http.Get("http://127.0.0.1" + server.Address + server.Endpoint)
	defer func() {
		_ = response.Body.Close()
	}()
	if err != nil {
		t.Fatalf("cannot fetch telemetry data from prometheus server: %v", err)
	}
	if response.Status != "200 OK" {
		t.Fatalf("the response code is not 200:%s", response.Status)
	}
}
