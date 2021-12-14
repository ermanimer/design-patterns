package cache

import (
	"testing"
)

func TestGetPut(t *testing.T) {
	c := NewCache(2).(*cache)

	item1 := Item{
		Key:   "key1",
		Value: "value1",
	}

	item2 := Item{
		Key:   "key2",
		Value: "value2",
	}

	item3 := Item{
		Key:   "key3",
		Value: "value3",
	}

	c.Put(item1)
	c.Put(item2)
	if c.items.Len() != 2 {
		t.Error("item count should be 2 after putting 2 keys")
	}

	c.Put(item1)
	if c.items.Len() != 2 {
		t.Error("item count should stay the same after putting an existing key")
	}

	c.Put(item3)
	if c.items.Len() != 2 {
		t.Error("item count should stay the same after putting a new key that overflows the cache")
	}
}

func TestGet(t *testing.T) {
	c := NewCache(3).(*cache)

	item1 := Item{
		Key:   "key1",
		Value: "value1",
	}

	item2 := Item{
		Key:   "key2",
		Value: "value2",
	}

	item3 := Item{
		Key:   "key3",
		Value: "value3",
	}

	c.Put(item1)
	c.Put(item2)

	item4 := c.Get(item3.Key)
	if item4 != nil {
		t.Error("a non-existing item's key should return a nil item")
	}

	item5 := c.Get(item1.Key)
	if item5.Value.(string) != item1.Value.(string) {
		t.Error("an existing items key's should return the existing item")
	}
}
