package pubsub

import (
	"errors"
	"sync"
)

// Pubsub defines the basic behaviors of a pubsub
type Pubsub interface {
	Subscribe(id string) chan string
	Unsubscribe(id string) error
	Publish(message string)
}

type pubsub struct {
	m           *sync.Mutex
	subscribers map[string]chan string
}

// compile time proof of interface implementation
var _ Pubsub = (*pubsub)(nil)

// NewPubsub creates and returns a new pubsub
func NewPubsub() Pubsub {
	return &pubsub{
		m:           &sync.Mutex{},
		subscribers: make(map[string]chan string),
	}
}

// Subscribe subscribes to pubsub
func (p *pubsub) Subscribe(id string) chan string {
	c := make(chan string)

	p.m.Lock()
	p.subscribers[id] = c
	p.m.Unlock()

	return c
}

// Unsubscribe unsubscribes from pubsub
func (p *pubsub) Unsubscribe(id string) error {
	c, ok := p.subscribers[id]
	if !ok {
		return errors.New("id not found")
	}

	p.m.Lock()
	delete(p.subscribers, id)
	close(c)
	p.m.Unlock()

	return nil
}

// Publish publishes message to the subscribers
func (p *pubsub) Publish(message string) {
	p.m.Lock()
	for _, c := range p.subscribers {
		select {
		case c <- message:
		default:
		}
	}
	p.m.Unlock()
}
