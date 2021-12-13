package cache

import (
	"container/list"
	"sync"
)

// Cache defines the behaviors of our cache
type Cache interface {
	Get(key string) interface{}
	Put(key string, value interface{})
}

// cache implements Cache interface
type cache struct {
	size  int                    // size is the capacity of the cache
	keys  *list.List             // keys holds the keys of the cached values
	items map[string]interface{} // items holds the cached key value pairs
	mutex sync.Mutex
}

// NewCache creates and returns a new cache
func NewCache(size int) Cache {
	if size < 1 {
		panic("invalid size")
	}

	return &cache{
		size:  size,
		keys:  list.New(),
		items: map[string]interface{}{},
		mutex: sync.Mutex{},
	}
}

// Get returns the value of a key from the cache.
// Returns nil if the key doesn't exist.
// Moves the key to the front to indicate that the key is recently used.
func (c *cache) Get(key string) interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	value, ok := c.items[key]
	if !ok {
		return nil
	}

	c.moveKeyToFront(key)

	return value
}

// Put puts a key value pair to the cache.
// Moves the existing pair's key to the front or pushes a new key to the front to indicate that the key is recently used.
// Removes the least recently used key value pair if cache is full.
func (c *cache) Put(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, ok := c.items[key]
	if ok {
		c.moveKeyToFront(key)
		return
	}

	if len(c.items) == c.size {
		lruKey := c.keys.Back().Value.(string)
		delete(c.items, lruKey)
		c.removeKey(lruKey)
	}

	c.pushKeyToFront(key)
	c.items[key] = value
}

// pushKeyToFront pushes a new key to the front
func (c *cache) pushKeyToFront(key string) {
	c.keys.PushFront(key)
}

// moveKeyToFront moves an existing key to the front
func (c *cache) moveKeyToFront(key string) {
	if e := c.getElement(key); e != nil {
		c.keys.MoveToFront(e)
	}
}

// removeKey removes a key
func (c *cache) removeKey(key string) {
	if e := c.getElement(key); e != nil {
		c.keys.Remove(e)
	}
}

// getElement returns the key's list element
func (c *cache) getElement(key string) *list.Element {
	for e := c.keys.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == key {
			return e
		}
	}

	return nil
}
