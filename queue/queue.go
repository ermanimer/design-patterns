package queue

import "context"

// Item represents queue item
type Item struct{}

// Queue defines the basic behaviours of a queue
type Queue interface {
	Enqueue(*Item)
	Dequeue(context.Context) (*Item, error)
}

type queue struct {
	items   chan []*Item
	isEmpty chan bool
}

// compile time proof of interface implementation
var _ Queue = (*queue)(nil)

// NewQueue creates and returns a new, empty queue
func NewQueue() Queue {
	// create items and isEmpty channels
	items := make(chan []*Item, 1)
	isEmpty := make(chan bool, 1)
	// mark queue as empty
	isEmpty <- true
	// return queue
	return &queue{items, isEmpty}
}

// Enqueue enqueues an item to the queue
func (q *queue) Enqueue(i *Item) {
	// create items if queue is empty or get items
	var items []*Item
	select {
	case items = <-q.items:
	case <-q.isEmpty:
	}
	// append item to items
	items = append(items, i)
	// update items
	q.items <- items
}

// Dequeue dequeues and returns an item from the queue
func (q *queue) Dequeue(c context.Context) (*Item, error) {
	// get items
	var items []*Item
	select {
	case items = <-q.items:
	case <-c.Done():
		return nil, c.Err()
	}
	// get item
	item := items[0]
	// mark queue as empty if last item is dequeued or update items
	if len(items) == 1 {
		q.isEmpty <- true
	} else {
		q.items <- items[1:]
	}
	//return item
	return item, nil
}
