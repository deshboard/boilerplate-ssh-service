package serverz

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// Manager holds the server objects and handles startup/shutdown
type Manager struct {
	servers []Server
	Logger  *log.Logger
}

type OperationIncompleteError struct {
	operation string
}

func (e OperationIncompleteError) Error() string {
	return fmt.Sprintf("Operation %s could not be completed for all servers", e.operation)
}

// NewManager returns a new Manager
func NewManager(servers ...Server) *Manager {
	return &Manager{
		servers: servers,
	}
}

// Serve starts all servers, but it doesn't wait for them to properly start and listen
// It returns a channel which receives errors from servers
func (m *Manager) Serve() <-chan error {
	ch := make(chan error, len(m.servers))

	for _, server := range m.servers {
		go func(server Server, ch chan<- error) {
			ch <- server.Serve()
		}(server, ch)
	}

	return ch
}

// Shutdown initiates graceful server shutdowns and waits for them to complete
func (m *Manager) Shutdown(ctx context.Context) error {
	var incomplete bool
	var wg sync.WaitGroup
	wg.Add(len(m.servers))

	for _, server := range m.servers {
		go func(server Server, ctx context.Context) {
			err := server.Shutdown(ctx)
			if err != nil {
				m.Logger.Println(err)
				incomplete = true
			}

			wg.Done()
		}(server, ctx)
	}

	wg.Wait()

	if incomplete {
		return OperationIncompleteError{"shutdown"}
	}

	return nil
}

// Close closes everything
func (m *Manager) Close() error {
	var incomplete bool
	var wg sync.WaitGroup
	wg.Add(len(m.servers))

	for _, server := range m.servers {
		err := server.Close()
		if err != nil {
			m.Logger.Println(err)
			incomplete = true
		}
	}

	if incomplete {
		return OperationIncompleteError{"close"}
	}

	return nil
}
