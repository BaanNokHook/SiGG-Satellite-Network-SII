// SiGG-Satellite-Network-SII  //

package resolvers

import (
	"fmt"
	"strings"

	"github.com/apache/skywalking-satellite/internal/pkg/log"

	"google.golang.org/grpc/resolver"
)

var staticServerSchema = "static"

type staticServerResolver struct {
}

func (s *staticServerResolver) Type() string {
	return staticServerSchema
}

func (s *staticServerResolver) BuildTarget(c *ServerFinderConfig) (string, error) {
	// build target using uri endpoint
	return fmt.Sprintf("%s:///%s", staticServerSchema, c.ServerAddr), nil
}

//nolint:gocritic // Implement for resolver.Target
func (*staticServerResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &staticResolver{
		target: target,
		cc:     cc,
	}
	r.analyzeClients()
	return r, nil
}

func (*staticServerResolver) Scheme() string {
	return staticServerSchema
}

type staticResolver struct {
	target resolver.Target
	cc     resolver.ClientConn
}

func (r *staticResolver) ResolveNow(o resolver.ResolveNowOptions) {
	r.analyzeClients()
}

func (*staticResolver) Close() {
}

func (r *staticResolver) analyzeClients() {
	addresses := strings.Split(strings.TrimLeft(r.target.URL.Path, "/"), ",")
	addrs := make([]resolver.Address, len(addresses))
	for i, s := range addresses {
		addrs[i] = resolver.Address{Addr: s}
	}
	if err := r.cc.UpdateState(resolver.State{Addresses: addrs}); err != nil {
		log.Logger.Warnf("error update static grpc client list: %v", err)
	}
}
