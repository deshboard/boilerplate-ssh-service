package main

import "github.com/goph/serverz"

// newServerQueue returns a new server queue with all the registered servers.
func newServerQueue(appCtx *application) *serverz.Queue {
	queue := serverz.NewQueue()
	queue.Manager.Logger = appCtx.logger

	server := newSSHServer(appCtx)
	queue.Append(server)

	debugServer := newDebugServer(appCtx)
	queue.Prepend(debugServer)

	return queue
}
