package processor

import (
	"github.com/apache/skywalking-satellite/internal/satellite/module/processor/api"
	filter "github.com/apache/skywalking-satellite/plugins/filter/api"
)

// Init Processor and dependency plugins
func NewProcessor(cfg *api.ProcessorConfig) api.Processor {
	p := &Processor{
		config:         cfg,
		runningFilters: []filter.Filter{},
	}
	for _, c := range p.config.FilterConfig {
		p.runningFilters = append(p.runningFilters, filter.GetFilter(c))
	}
	return p
}
