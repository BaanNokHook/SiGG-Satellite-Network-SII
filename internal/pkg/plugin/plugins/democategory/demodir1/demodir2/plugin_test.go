// SiGG-Satellite-Network-SII  //

package demodir2

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/pkg/log"
	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
)

const pluginName = "demo-plugin"

type DemoCategory interface {
	plugin.Plugin
	Say() string
}

type DemoPlugin struct {
	config.CommonFields
	Organization string `mapstructure:"organization"`
	Project      string `mapstructure:"project"`
}

func (d *DemoPlugin) Say() string {
	return d.Organization + ":" + d.Project
}

func (d *DemoPlugin) Name() string {
	return "demo-plugin"
}

func (d *DemoPlugin) ShowName() string {
	return "Demo Plugin"
}

func (d *DemoPlugin) Description() string {
	return "this is just a demo"
}

func (d *DemoPlugin) DefaultConfig() string {
	return `
organization: "ASF"
project: "skywalking-satellite"
`
}

func TestPlugin(t *testing.T) {
	tests := []struct {
		name string
		args plugin.Config
		want *DemoPlugin
	}{
		{
			name: "test1",
			args: plugin.Config{
				"plugin_name":            pluginName,
				"organization":           "CNCF",
				"project":                "Fluentd",
				"commonfields_pipe_name": "b",
			},
			want: &DemoPlugin{
				CommonFields: config.CommonFields{
					PipeName: "b",
				},
				Organization: "CNCF",
				Project:      "Fluentd",
			},
		},
		{
			name: "demoplugin",
			args: plugin.Config{
				"plugin_name": pluginName,
			},
			want: &DemoPlugin{
				Organization: "ASF",
				Project:      "skywalking-satellite",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if i := recover(); i != nil {
					t.Errorf("the plugin initialized err: %v", i)
				}
			}()
			p := plugin.Get(reflect.TypeOf((*DemoCategory)(nil)).Elem(), tt.args)
			if !cmp.Equal(p, tt.want) {
				t.Errorf("Format() got = %v, want %v", p, tt.want)
			}
		})
	}
}

func init() {
	log.Init(new(log.LoggerConfig))
	plugin.RegisterPluginCategory(reflect.TypeOf((*DemoCategory)(nil)).Elem())
	plugin.RegisterPlugin(new(DemoPlugin))
}
