package buffer

import (
	"testing"

	"github.com/apache/skywalking-satellite/internal/pkg/log"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
)

func TestNewBuffer(t *testing.T) {
	buffer := NewBatchBuffer(3)
	tests := []struct {
		name string
		args *event.OutputEventContext
		want int
	}{
		{
			name: "add-1",
			args: &event.OutputEventContext{Offset: &event.Offset{Position: "1"}},
			want: 1,
		},
		{
			name: "add-2",
			args: &event.OutputEventContext{Offset: &event.Offset{Position: "1"}},
			want: 2,
		},
		{
			name: "add-3",
			args: &event.OutputEventContext{Offset: &event.Offset{Position: "1"}},
			want: 3,
		},
		{
			name: "add-4",
			args: &event.OutputEventContext{Offset: &event.Offset{Position: "1"}},
			want: 3,
		},
		{
			name: "add-5",
			args: &event.OutputEventContext{Offset: &event.Offset{Position: "1"}},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer.Add(tt.args)
			if got := buffer.Len(); got != tt.want {
				t.Errorf("Buffer Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func init() {
	log.Init(&log.LoggerConfig{})
}
