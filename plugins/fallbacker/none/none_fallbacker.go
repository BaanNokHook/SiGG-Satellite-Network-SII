// SiGG-Satellite-Network-SII  //

package none

import (
	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/satellite/event"
	"github.com/apache/skywalking-satellite/plugins/forwarder/api"
)

const (
	Name     = "none-fallbacker"
	ShowName = "None Fallbacker"
)

type Fallbacker struct {
	config.CommonFields
}

func (f *Fallbacker) Name() string {
	return Name
}

func (f *Fallbacker) ShowName() string {
	return ShowName
}

func (f *Fallbacker) Description() string {
	return "The fallbacker would do nothing when facing failure data."
}

func (f *Fallbacker) DefaultConfig() string {
	return ""
}

func (f *Fallbacker) FallBack(batch event.BatchEvents, forward api.ForwardFunc) bool {
	return true
}
