package sharing

import (
	"fmt"
	"sync"

	"github.com/apache/skywalking-satellite/internal/pkg/log"
	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
	"github.com/apache/skywalking-satellite/internal/satellite/config"
	client "github.com/apache/skywalking-satellite/plugins/client/api"
	server "github.com/apache/skywalking-satellite/plugins/server/api"
)

// Manager contains the sharing plugins, only supports client and server plugins.
var Manager map[string]plugin.SharingPlugin
var once sync.Once

// Load loads the sharing config to the Manager.
func Load(cfg *config.SharingConfig) {
	once.Do(func() {
		Manager = make(map[string]plugin.SharingPlugin)
		for _, c := range cfg.Clients {
			p := client.GetClient(c)
			Manager[p.Name()] = p
		}
		for _, c := range cfg.Servers {
			p := server.GetServer(c)
			Manager[p.Name()] = p
		}
	},
	)
}

func Prepare() error {
	for _, sharingPlugin := range Manager {
		if err := sharingPlugin.Prepare(); err != nil {
			log.Logger.Errorf("error in closing the %s sharing plugin: %v", sharingPlugin.Name(), err)
			Close()
			return fmt.Errorf("cannot preare the sharing plugins named %s: %v", sharingPlugin.Name(), err)
		}
	}
	return nil
}

func Start() error {
	for _, sharingPlugin := range Manager {
		if err := sharingPlugin.Start(); err != nil {
			log.Logger.Errorf("error in closing the %s sharing plugin: %v", sharingPlugin.Name(), err)
			Close()
			return fmt.Errorf("cannot preare the sharing plugins named %s: %v", sharingPlugin.Name(), err)
		}
	}
	return nil
}

func Close() {
	for _, sharingPlugin := range Manager {
		if err := sharingPlugin.Close(); err != nil {
			log.Logger.Errorf("error in closing the %s sharing plugin: %v", sharingPlugin.Name(), err)
		}
	}
}
