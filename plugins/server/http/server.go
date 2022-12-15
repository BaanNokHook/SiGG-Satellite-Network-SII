// SiGG-Satellite-Network-SII  //

package http

import (
	"net/http"

	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/pkg/log"
)

const (
	Name     = "http-server"
	ShowName = "HTTP Server"
)

type Server struct {
	config.CommonFields
	Address string         `mapstructure:"address"`
	Server  *http.ServeMux // The http server.
}

func (s *Server) Name() string {
	return Name
}

func (s *Server) ShowName() string {
	return ShowName
}

func (s *Server) Description() string {
	return "This is a sharing plugin, which would start a http server."
}

func (s *Server) DefaultConfig() string {
	return `
# The http server address.
address: ":12800"
`
}

func (s *Server) Prepare() error {
	s.Server = http.NewServeMux()
	return nil
}

func (s *Server) Start() error {
	log.Logger.WithField("address", s.Address).Info("http server is starting...")
	go func() {
		err := http.ListenAndServe(s.Address, s.Server)
		if err != nil {
			log.Logger.WithField("address", s.Address).Infof("http server has failure when starting: %v", err)
		}
	}()
	return nil
}

func (s *Server) Close() error {
	log.Logger.Info("http server is closed")
	return nil
}

func (s *Server) GetServer() interface{} {
	return s
}
