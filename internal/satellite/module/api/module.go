package api

import (
	"context"
	"time"
)

// ShutdownHookTime is the global shutdown hook time.
var ShutdownHookTime = time.Second * 5

// Module id a custom plugin interface, which defines the processing.
type Module interface {
	// Prepare would do some preparing workers, such build connection with external services.
	Prepare() error
	// Boot would start the module and return error when started failed. When a stop signal received
	// or an exception occurs, the shutdown function would be called.
	Boot(ctx context.Context)

	// Shutdown could do some clean job to close Module.
	Shutdown()
}
