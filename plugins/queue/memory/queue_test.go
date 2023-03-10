// SiGG-Satellite-Network-SII  //

package memory

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	_ "github.com/apache/skywalking-satellite/internal/satellite/test"
	"github.com/apache/skywalking-satellite/plugins/queue/api"

	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

func initMemoryQueue(cfg plugin.Config) (*Queue, error) {
	plugin.RegisterPluginCategory(reflect.TypeOf((*api.Queue)(nil)).Elem())
	plugin.RegisterPlugin(&Queue{})
	var config plugin.Config = map[string]interface{}{
		plugin.NameField: Name,
	}
	for k, v := range cfg {
		config[k] = v
	}
	q := api.GetQueue(config)
	if q == nil {
		return nil, fmt.Errorf("cannot get a memoory queue from the registry")
	}
	if err := q.Initialize(); err != nil {
		return nil, fmt.Errorf("queue cannot initialize: %v", err)
	}
	return q.(*Queue), nil
}

func TestQueue_Enqueue(t *testing.T) {
	const num = 100000
	q, err := initMemoryQueue(map[string]interface{}{
		"event_buffer_size": num,
	})
	if err != nil {
		t.Fatalf("cannot init the memory queue: %v", err)
	}

	if _, err := q.Dequeue(); err == nil {
		t.Fatal("the dequeue want failure but success")
	}

	// enqueue
	for i := 0; i <= num; i++ {
		e := &v1.SniffData{
			Name: strconv.Itoa(i),
		}
		if i < num {
			if err := q.Enqueue(e); err != nil {
				t.Fatalf("the enqueue want seuccess but failure: %v", err)
			}
		} else {
			if err := q.Enqueue(e); err == nil {
				t.Fatal("the enqueue want failure but success when facing full")
			}
		}
	}

	// dequeue
	for i := 0; i < num; i++ {
		if e, err := q.Dequeue(); err != nil {
			t.Fatalf("the dequeue want seuccess but failure: %v", err)
		} else if e.Event.Name != strconv.Itoa(i) {
			t.Fatalf("want got %s, but got %s", strconv.Itoa(i), e.Event.Name)
		}
	}
	if _, err := q.Dequeue(); err == nil {
		t.Fatal("the dequeue want failure but success")
	}
}
