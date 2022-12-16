package buffer

import (
	"context"
	"testing"
	"time"
)

func TestLimitor(t *testing.T) {
	tests := []struct {
		name       string
		limitCount int
		flushTime  int
		sendCount  int
		isFlush    bool
	}{
		{
			name:       "reach limit to flush",
			limitCount: 5,
			flushTime:  5000,
			sendCount:  5,
			isFlush:    true,
		},
		{
			name:       "not reach limit",
			limitCount: 5,
			flushTime:  5000,
			sendCount:  2,
			isFlush:    false,
		},
		{
			name:       "using flush time to flush",
			limitCount: 5,
			flushTime:  100,
			sendCount:  1,
			isFlush:    true,
		},
	}
	for _, ts := range tests {
		t.Run(ts.name, func(t *testing.T) {
			conf := LimiterConfig{LimitCount: ts.limitCount, FlushTime: ts.flushTime}
			limiter := NewLimiter(conf, func() int {
				return ts.sendCount
			})

			flushChannel := make(chan struct{}, 1)
			defer limiter.Stop()
			limiter.Start(context.Background(), func() {
				flushChannel <- struct{}{}
			})
			limiter.Check()

			beenFlushed := false
			select {
			case <-flushChannel:
				beenFlushed = true
			case <-time.After(time.Second * 1):
			}

			if ts.isFlush != beenFlushed {
				t.Fail()
			}
		})
	}
}
