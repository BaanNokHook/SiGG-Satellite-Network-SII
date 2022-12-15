// SiGG-Satellite-Network-SII  //

package resolvers

import (
	"fmt"

	"google.golang.org/grpc/resolver"
)

// all customized resolvers
var rs = []GrpcResolver{
	&staticServerResolver{},
	&kubernetesServerResolver{},
}

type ServerFinderConfig struct {
	FinderType string `mapstructure:"finder_type"` // The gRPC server address finder type, support "static" and "kubernetes"
	// The gRPC server address, only works for "static" address finder
	ServerAddr string `mapstructure:"server_addr"`
	// The kubernetes config to lookup addresses, only works for "kubernetes" address finder
	KubernetesConfig *KubernetesConfig `mapstructure:"kubernetes_config"`
}

type GrpcResolver interface {
	resolver.Builder

	// Type of resolver
	Type() string
	// BuildTarget address by client config
	BuildTarget(c *ServerFinderConfig) (string, error)
}

func RegisterAllGrpcResolvers() {
	for _, r := range rs {
		resolver.Register(r)
	}
}

func BuildTarget(client *ServerFinderConfig) (string, error) {
	for _, r := range rs {
		if client.FinderType == r.Type() {
			return r.BuildTarget(client)
		}
	}
	return "", fmt.Errorf("could not find client finder: %s", client.FinderType)
}
