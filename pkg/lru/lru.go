// Package lru implements a LRU cache.
package lru

import (
	"container/list"
	"sync"
)

// Cache is a LRU cache.
type Cache struct {
	// MaxEntries is the maximum number of cache entries
	// before an item is purged. Zero means no limit.
	MaxEntries int

	// OnPurged specifies a function to be executed
	// when an entry is purged from the cache.
	OnPurged func(key interface{}, value interface{})

	ll    *list.List
	cache map[interface{}]*list.Element
	mu    sync.RWMutex
}

type entry struct {
	key   interface{}
	value interface{}
}

// New creates a new cache, if maxEntries is zero, the cache has no limit.
func New(maxEntries int) *Cache {
	if maxEntries < 0 {
		panic("maxEntries can not be less than zero")
	}
	return &Cache{
		MaxEntries: maxEntries,
		ll:         list.New(),
		cache:      make(map[interface{}]*list.Element),
	}
}

// Add adds value to the cache.
func (c *Cache) Add(key interface{}, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		c.cache = make(map[interface{}]*list.Element)
		c.ll = list.New()
	}
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		e.Value.(*entry).value = value
		return
	}
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele
	if c.MaxEntries > 0 && c.ll.Len() > c.MaxEntries {
		c.removeOldest(false)
	}
}

// Get looks up value by key from the cache.
func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *Cache) Remove(key interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *Cache) RemoveOldest() {
	c.removeOldest(true)
}

func (c *Cache) removeOldest(toLock bool) {
	if toLock {
		c.mu.Lock()
		defer c.mu.Unlock()
	}

	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
	if c.OnPurged != nil {
		c.OnPurged(kv.key, kv.value)
	}
}

// Len returns the number of items in the cache.
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Clear purges all items from the cache.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.OnPurged != nil {
		for _, e := range c.cache {
			kv := e.Value.(*entry)
			c.OnPurged(kv.key, kv.value)
		}
	}
	c.ll = nil
	c.cache = nil
}
