// SiGG-Satellite-Network-SII  //

package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/status"

	"github.com/apache/skywalking-satellite/plugins/client/api"
)

// sniffer
func (c *Client) snifferChannelStatus() {
	ctx, cancel := context.WithCancel(c.ctx)
	defer cancel()
	timeTicker := time.NewTicker(time.Duration(c.CheckPeriod) * time.Second)
	for {
		select {
		case <-timeTicker.C:
			state := c.client.GetState()
			if state == connectivity.Shutdown || state == connectivity.TransientFailure {
				c.updateStatus(api.Disconnect)
			} else if state == connectivity.Ready || state == connectivity.Idle {
				c.updateStatus(api.Connected)
			}
		case <-ctx.Done():
			timeTicker.Stop()
			return
		}
	}
}

func (c *Client) reportError(err error) {
	if err == nil {
		return
	}
	fromError, ok := status.FromError(err)
	if ok {
		errCode := fromError.Code()
		if errCode == codes.Unavailable || errCode == codes.PermissionDenied ||
			errCode == codes.Unauthenticated || errCode == codes.ResourceExhausted || errCode == codes.Unknown {
			c.updateStatus(api.Disconnect)
		}
	}
}

func (c *Client) updateStatus(clientStatus api.ClientStatus) {
	if c.status != clientStatus {
		c.status = clientStatus
		for _, listener := range c.listeners {
			listener <- c.status
		}
	}
}
