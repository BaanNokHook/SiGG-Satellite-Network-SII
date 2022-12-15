// SiGG-Satellite-Network-SII  //

package resolvers

import (
	"fmt"
	"strings"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/discovery/kubernetes"
	"github.com/prometheus/prometheus/discovery/targetgroup"
)

type ServiceAnalyzer struct {
}

func (p *ServiceAnalyzer) KindType() string {
	return string(kubernetes.RoleService)
}

func (p *ServiceAnalyzer) GetAddresses(cache map[string]*targetgroup.Group, config *KubernetesConfig) []string {
	result := make([]string, 0)
	for _, group := range cache {
		for _, target := range group.Targets {
			address := string(target[model.LabelName("__address__")])
			if strings.HasSuffix(address, fmt.Sprintf(":%d", config.ExtraPort.Port)) {
				result = append(result, address)
			}
		}
	}
	return result
}
