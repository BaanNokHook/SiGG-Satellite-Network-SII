// SiGG-Satellite-Network-SII  //

package resolvers

import (
	"strconv"

	"github.com/prometheus/prometheus/discovery/kubernetes"

	"github.com/prometheus/prometheus/discovery/targetgroup"

	"github.com/prometheus/common/model"
)

type PodAnalyzer struct {
}

func (p *PodAnalyzer) KindType() string {
	return string(kubernetes.RolePod)
}

func (p *PodAnalyzer) GetAddresses(cache map[string]*targetgroup.Group, config *KubernetesConfig) []string {
	result := make([]string, 0)
	for _, group := range cache {
		for _, target := range group.Targets {
			val, exists := target[model.LabelName("__meta_kubernetes_pod_container_port_number")]
			if exists && string(val) == strconv.Itoa(config.ExtraPort.Port) {
				result = append(result, string(target[model.LabelName("__address__")]))
			}
		}
	}
	return result
}
