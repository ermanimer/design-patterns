package pool

import (
	"context"
	"net"
)

// Pool defines the basic behaviors of a connection pool
type Pool interface {
	Acquire(context.Context) (net.Conn, error)
	Release(net.Conn)
}

type token struct{}

type pool struct {
	address         string
	semaphore       chan *token
	idleConnections chan net.Conn
}

// compile time proof of interface implementation
var _ Pool = (*pool)(nil)

// NewPool creates and returns a new connection pool
func NewPool(address string, limit int) Pool {
	semaphore := make(chan *token, limit)
	idleConnections := make(chan net.Conn, limit)
	return &pool{address, semaphore, idleConnections}
}

// Acquire acquires an idle connection from the pool
func (p *pool) Acquire(c context.Context) (net.Conn, error) {
	// get an idle connection or create a new connection if semaphore is empty
	select {
	case connection := <-p.idleConnections:
		return connection, nil
	case p.semaphore <- &token{}:
		conn, err := net.Dial("tcp", p.address)
		if err != nil {
			<-p.semaphore
			return nil, err
		}
		return conn, err
	case <-c.Done():
		return nil, c.Err()
	}
}

// Release releases the connection back to the pool
func (p *pool) Release(c net.Conn) {
	p.idleConnections <- c
}
