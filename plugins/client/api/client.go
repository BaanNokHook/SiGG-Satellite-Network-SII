package api

import (
	"reflect"

	"github.com/apache/skywalking-satellite/internal/pkg/plugin"
)

// The client statuses.
const (
	_ ClientStatus = iota
	Connected
	Disconnect
)

// ClientStatus represents the status of the client.
type ClientStatus int8

// Client is a plugin interface, that defines new clients, such as gRPC client and Kafka client.
type Client interface {
	plugin.SharingPlugin

	// GetConnectedClient returns the connected client to publish events.
	GetConnectedClient() interface{}
	// RegisterListener register a listener to listen the client status.
	RegisterListener(chan<- ClientStatus)
}

// GetClient gets an initialized client plugin.
func GetClient(config plugin.Config) Client {
	return plugin.Get(reflect.TypeOf((*Client)(nil)).Elem(), config).(Client)
}
