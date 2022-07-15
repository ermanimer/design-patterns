package cache

import (
	"container/list"
	"sync"
)

// Item represents a key-value pair
type Item struct {
	Key   string
	Value interface{}
}

// Cache defines the behaviors of our cache
type Cache interface {
	Get(key string) *Item
	Put(Item)
}

// cache implements Cache interface
type cache struct {
	size  int
	items *list.List
	mutex sync.RWMutex
}

// NewCache creates and returns a new cache
func NewCache(size int) Cache {
	if size < 1 {
		panic("invalid size")
	}

	return &cache{
		size:  size,
		items: list.New(),
		mutex: sync.RWMutex{},
	}
}

// Get returns an existing item from the cache.
// Get also moves the existing item to the front of the items list to indicate that the existing item is recently used.
func (c *cache) Get(key string) *Item {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	e := c.getElement(key)
	if e == nil {
		return nil
	}

	c.items.MoveToFront(e)

	i := e.Value.(Item)

	return &i
}

// Put puts a new item into the cache.
// Put removes the least recently used item from the items list when the cache is full.
// Put pushes the new item to the front of the items list to indicate that the new item is recently used.
func (c *cache) Put(i Item) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	e := c.getElement(i.Key)
	if e != nil {
		c.items.MoveToFront(e)
		return
	}

	if c.items.Len() == c.size {
		c.items.Remove(c.items.Back())
	}

	c.items.PushFront(i)
}

// getElement returns list element of an existing item
func (c *cache) getElement(key string) *list.Element {
	for e := c.items.Front(); e != nil; e = e.Next() {
		i := e.Value.(Item)
		if i.Key == key {
			return e
		}
	}

	return nil
}
