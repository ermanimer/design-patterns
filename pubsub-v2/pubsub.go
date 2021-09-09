package pubsub

import (
	"errors"
	"sync"
)

// actions
const (
	subscribe = 1 + iota
	unsubscribe
	publish
)

type internalMessage struct {
	action  int
	id      string
	mc      chan string
	message string
}

// Pubsub defines behaviors of a pubsub
type Pubsub interface {
	Start()
	Stop()
	Subscribe(id string) chan string
	Unsubscribe(id string) error
	Publish(message string)
}

type pubsub struct {
	wg          *sync.WaitGroup
	imc         chan *internalMessage
	subscribers map[string]chan string
}

// compile time proof of interface implementation
var _ Pubsub = (*pubsub)(nil)

// NewPubsub creates and returns a new pubsub
func NewPubsub() Pubsub {
	return &pubsub{
		wg:          &sync.WaitGroup{},
		imc:         make(chan *internalMessage),
		subscribers: make(map[string]chan string),
	}
}

// Start starts pubsub
func (p *pubsub) Start() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		for im := range p.imc {
			switch im.action {
			case subscribe:
				p.subscribers[im.id] = im.mc
			case unsubscribe:
				delete(p.subscribers, im.id)
				close(im.mc)
			case publish:
				for id := range p.subscribers {
					p.wg.Add(1)
					go func(id string) {
						defer p.wg.Done()

						mc, ok := p.subscribers[id]
						if ok {
							mc <- im.message
						}
					}(id)
				}
			}
		}
	}()
}

// Stop stops pubsub
func (p *pubsub) Stop() {
	close(p.imc)
	p.wg.Wait()
}

// Subscribe subscribes to pubsub
func (p *pubsub) Subscribe(id string) chan string {
	mc := make(chan string)

	im := &internalMessage{
		action: subscribe,
		id:     id,
		mc:     mc,
	}

	p.imc <- im

	return mc
}

// Unsubscribe unsubscribes from pubsub
func (p *pubsub) Unsubscribe(id string) error {
	mc, ok := p.subscribers[id]
	if !ok {
		return errors.New("message channel not found")
	}

	im := &internalMessage{
		action: unsubscribe,
		id:     id,
		mc:     mc,
	}

	p.imc <- im

	return nil
}

// Publish publishes a message to pubsub
func (p *pubsub) Publish(message string) {
	im := &internalMessage{
		action:  publish,
		message: message,
	}

	p.imc <- im
}
