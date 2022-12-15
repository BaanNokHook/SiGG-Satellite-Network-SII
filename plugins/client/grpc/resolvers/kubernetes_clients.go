// SiGG-Satellite-Network-SII  //

package resolvers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/grpc/resolver"

	"github.com/apache/skywalking-satellite/internal/pkg/log"
)

var kubernetesServerSchema = "kubernetes"

type kubernetesServerResolver struct {
}

func (k *kubernetesServerResolver) Type() string {
	return kubernetesServerSchema
}

func (k *kubernetesServerResolver) BuildTarget(c *ServerFinderConfig) (string, error) {
	marshal, err := json.Marshal(c.KubernetesConfig)
	if err != nil {
		return "", fmt.Errorf("convert kubernetes config error: %v", err)
	}
	return fmt.Sprintf("%s:///%s", kubernetesServerSchema, string(marshal)), nil
}

//nolint:gocritic // Implement for resolver.Target
func (*kubernetesServerResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	// convert data
	kubernetesConfig := &KubernetesConfig{}
	if err := json.Unmarshal([]byte(strings.TrimLeft(target.URL.Path, "/")), kubernetesConfig); err != nil {
		return nil, fmt.Errorf("could not analyze the address: %v", err)
	}

	// validate http config
	if kubernetesConfig.APIServer != "" {
		httpConfig, err := kubernetesConfig.HTTPClientConfig.convertHTTPConfig()
		if err != nil {
			return nil, err
		}
		if err = httpConfig.Validate(); err != nil {
			return nil, fmt.Errorf("http config validate error: %v", err)
		}
	}

	// init cache
	ctx, cancel := context.WithCancel(context.Background())
	cache, err := NewKindCache(ctx, kubernetesConfig, cc)
	if err != nil {
		cancel()
		return nil, err
	}

	// build resolver
	r := &kubernetesResolver{
		cache:  cache,
		cancel: cancel,
	}
	return r, nil
}

func (*kubernetesServerResolver) Scheme() string {
	return kubernetesServerSchema
}

type kubernetesResolver struct {
	cache  *KindCache
	cancel context.CancelFunc
}

func (k *kubernetesResolver) ResolveNow(o resolver.ResolveNowOptions) {
	if err := k.cache.UpdateAddresses(); err != nil {
		log.Logger.Warnf("error update static grpc client list: %v", err)
	}
}

func (k *kubernetesResolver) Close() {
	k.cancel()
}
