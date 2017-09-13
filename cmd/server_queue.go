package main

import "github.com/goph/serverz"

// newServerQueue returns a new server queue with all the registered servers.
func newServerQueue(a *application) *serverz.Queue {
	queue := serverz.NewQueue()

	debugServer := newDebugServer(a)
	queue.Prepend(debugServer, nil)

	server := newSSHServer(a)
	queue.Append(server, nil)

	return queue
}
