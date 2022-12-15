// SiGG-Satellite-Network-SII  //

package kafka

import (
	"context"
	"time"

	"github.com/apache/skywalking-satellite/plugins/client/api"
)

// snifferBrokerStatus would sniffer the broker status to notify the listeners.
func (c *Client) snifferBrokerStatus() {
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()
	timeTicker := time.NewTicker(time.Duration(c.RefreshPeriod) * time.Minute)
	for {
		select {
		case <-timeTicker.C:
			brokers := c.client.Brokers()
			if len(brokers) == 0 && c.status == api.Connected {
				c.status = api.Disconnect
				c.notify()
			} else if len(brokers) > 0 && c.status == api.Disconnect {
				c.status = api.Connected
				c.notify()
			}
		case <-ctx.Done():
			timeTicker.Stop()
			return
		}
	}
}

// notify the current status to the listeners.
func (c *Client) notify() {
	for _, listener := range c.listeners {
		listener <- c.status
	}
}
