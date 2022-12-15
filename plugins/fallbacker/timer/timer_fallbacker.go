// SiGG-Satellite-Network-SII  //

package timer

import (
	"time"

	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
	"github.com/apache/skywalking-satellite/plugins/forwarder/api"
)

const (
	Name     = "timer-fallbacker"
	ShowName = "Timer Fallbacker"
)

// Fallbacker is a timer fallbacker when forward fails.
type Fallbacker struct {
	config.CommonFields
	MaxAttempts        int `mapstructure:"max_attempts"`
	ExponentialBackoff int `mapstructure:"exponential_backoff"`
	MaxBackoff         int `mapstructure:"max_backoff"`
}

func (t *Fallbacker) Name() string {
	return Name
}

func (t *Fallbacker) ShowName() string {
	return ShowName
}

func (t *Fallbacker) Description() string {
	return "This is a timer fallback trigger to process the forward failure data."
}

func (t *Fallbacker) DefaultConfig() string {
	return `
# The forwarder max attempt times.
max_attempts: 3
# The exponential_backoff is the standard retry duration, and the time for each retry is expanded
# by 2 times until the number of retries reaches the maximum.(Time unit is millisecond.)
exponential_backoff: 2000
# The max backoff time used in retrying, which would override the latency time when the latency time
# with exponential increasing larger than it.(Time unit is millisecond.)
max_backoff: 5000
`
}

func (t *Fallbacker) FallBack(batch event.BatchEvents, forward api.ForwardFunc) bool {
	currentLatency := t.ExponentialBackoff
	for i := 1; i < t.MaxAttempts; i++ {
		time.Sleep(time.Duration(currentLatency) * time.Millisecond)
		if err := forward(batch); err != nil {
			currentLatency *= 2
			if currentLatency > t.MaxBackoff {
				currentLatency = t.MaxBackoff
			}
		} else {
			return true
		}
	}
	return false
}
