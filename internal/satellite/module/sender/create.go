package sender

import (
	"github.com/apache/skywalking-satellite/internal/satellite/module/sender/api"
	"github.com/apache/skywalking-satellite/internal/satellite/sharing"
	client "github.com/apache/skywalking-satellite/plugins/client/api"
	fallbacker "github.com/apache/skywalking-satellite/plugins/fallbacker/api"
	forwarder "github.com/apache/skywalking-satellite/plugins/forwarder/api"
)

// NewSender crate a Sender.
func NewSender(cfg *api.SenderConfig) api.Sender {
	s := &Sender{
		config:            cfg,
		runningForwarders: []forwarder.Forwarder{},
		runningFallbacker: fallbacker.GetFallbacker(cfg.FallbackerConfig),
		runningClient:     sharing.Manager[cfg.ClientName].(client.Client),
		listener:          make(chan client.ClientStatus),
	}
	for _, c := range s.config.ForwardersConfig {
		s.runningForwarders = append(s.runningForwarders, forwarder.GetForwarder(c))
	}
	return s
}
