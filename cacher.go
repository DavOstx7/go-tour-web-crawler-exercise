package main

import "sync"

type Cacher struct {
	mu   sync.Mutex
	keys map[string]bool
}

func (c *Cacher) IsCached(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, exists := c.keys[key]
	return exists
}

func (c *Cacher) Cache(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.keys[key] = true
}
