// SiGG-Satellite-Network-SII  //

package lb

import "context"

type LoadBalancerConfig struct {
	appointAddr string
	routeKey    string
}

type ctxKey struct{}

var ctxKeyInstance = ctxKey{}

func WithLoadBalanceConfig(ctx context.Context, routeKey, appointAddr string) context.Context {
	return context.WithValue(ctx, ctxKeyInstance, &LoadBalancerConfig{
		routeKey:    routeKey,
		appointAddr: appointAddr,
	})
}

func GetAddress(ctx context.Context) string {
	if config := queryConfig(ctx); config != nil {
		return config.appointAddr
	}
	return ""
}

func queryConfig(ctx context.Context) *LoadBalancerConfig {
	value := ctx.Value(ctxKeyInstance)
	if value == nil {
		return nil
	}
	return value.(*LoadBalancerConfig)
}
