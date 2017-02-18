package serverz

import "context"

// Server in this context is an abstraction over anything that can be started and stopped (either gracefully or forcefully)
// Typically they accept connections and serve over network, like HTTP or RPC servers
type Server interface {
	// Serve is a blocking operation and is usually called in a separate goroutine
	// The returned error is always non-nil, which represents runtime errors as well as errors caused by shutdown operations
	// This is usually not a problem as shutdowns break the main loop after which the program exits
	Serve() error

	// Shutdown causes the server tp exit gracefully
	// Normally this could mean we wait indefinitely for connections to get terminated
	// Therefore a context can be passed which can be cancelled to instruct the shutdown handler to give control back
	// and optionally proceed with a forceful exit
	Shutdown(ctx context.Context) error

	// Close forcefully terminates all processes and connections
	// It also terminates any graceful shutdown attempts
	Close() error
}
