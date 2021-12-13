package cache

import (
	"testing"
)

func TestGetPut(t *testing.T) {
	c := NewCache(2).(*cache)

	key1 := "key1"
	value1 := "value1"
	key2 := "key2"
	value2 := "value2"
	key3 := "key3"
	value3 := "value3"

	c.Put(key1, value1)
	c.Put(key2, value2)
	if len(c.items) != 2 {
		t.Error("item count should be 2 after putting 2 keys")
	}

	c.Put(key1, value1)
	if len(c.items) != 2 {
		t.Error("item count should stay the same after putting an existing key")
	}

	c.Put(key3, value3)
	if len(c.items) != 2 {
		t.Error("item count should stay the same after putting a new key that overflows the cache")
	}
}

func TestGet(t *testing.T) {
	c := NewCache(3).(*cache)

	key1 := "key1"
	value1 := "value1"
	key2 := "key2"
	value2 := "value2"
	key3 := "key3"

	c.Put(key1, value1)
	c.Put(key2, value2)

	value3 := c.Get(key3)
	if value3 != nil {
		t.Error("a non-existing key's value should be nil")
	}

	value4 := c.Get(key1).(string)
	if value4 != value1 {
		t.Error("an existing key value pair don't match")
	}
}
