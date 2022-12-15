// SiGG-Satellite-Network-SII  //

package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc/peer"
)

func GetPeerHostFromStreamContext(ctx context.Context) string {
	peerAddr := GetPeerAddressFromStreamContext(ctx)
	if inx := strings.IndexByte(peerAddr, ':'); inx > 0 {
		peerAddr = peerAddr[:strings.IndexByte(peerAddr, ':')]
	}
	return peerAddr
}

func GetPeerAddressFromStreamContext(ctx context.Context) string {
	if peerAddr, ok := peer.FromContext(ctx); ok {
		return peerAddr.Addr.String()
	}
	return ""
}
