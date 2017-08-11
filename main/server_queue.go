package main

import "github.com/goph/serverz"

// newServerQueue returns a new server queue with all the registered servers.
func newServerQueue(appCtx *application) *serverz.Queue {
	queue := serverz.NewQueue()
	queue.Manager.Logger = appCtx.logger

	server := newHTTPServer(appCtx)
	queue.Append(server)

	healthServer := newHealthServer(appCtx)
	queue.Prepend(healthServer)

	return queue
}
