// SiGG-Satellite-Network-SII  //

package timer

import (
	"errors"
	"reflect"
	"testing"

	"github.com/apache/skywalking-satellite/internal/pkg/log"
	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	"github.com/apache/skywalking-satellite/plugins/fallbacker/api"
)

func initFallbacker(cfg plugin.Config) *Fallbacker {
	plugin.RegisterPluginCategory(reflect.TypeOf((*api.Fallbacker)(nil)).Elem())
	plugin.RegisterPlugin(new(Fallbacker))
	cfg[plugin.NameField] = Name
	q := api.GetFallbacker(cfg)
	if q == nil {
		log.Logger.Errorf("cannot get a default config fallbacker from the registry")
		return nil
	}
	return q.(*Fallbacker)
}

func TestFallbacker_FallBack1(t1 *testing.T) {
	count := 0
	mockForwarderFunc := func(_ event.BatchEvents) error {
		count++
		if count < 4 {
			return errors.New("mock error")
		}
		return nil
	}
	tests := []struct {
		name      string
		args      plugin.Config
		want      bool
		wantCount int
	}{
		{
			name:      "default-fallbacker",
			args:      plugin.Config{},
			want:      false,
			wantCount: 2,
		},
		{
			name: "test-reach-max_attempts",
			args: plugin.Config{
				"max_attempts":        5,
				"exponential_backoff": 200,
				"max_backoff":         3000,
			},
			want:      true,
			wantCount: 4,
		},
		{
			name: "test-unreach-max_attempts",
			args: plugin.Config{
				"max_attempts":        10,
				"exponential_backoff": 20,
				"max_backoff":         30000000,
			},
			want:      true,
			wantCount: 4,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			f := initFallbacker(tt.args)
			count = 0
			if got := f.FallBack(make(event.BatchEvents, 0), mockForwarderFunc); got != tt.want {
				t1.Errorf("FallBack() = %v, want %v", got, tt.want)
			}
			if count != tt.wantCount {
				t1.Errorf("Fallback count = %v, want %v", count, tt.wantCount)
			}
		})
	}
}
