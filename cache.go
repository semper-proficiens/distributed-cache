package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Cacher interface {
	Set([]byte, []byte, time.Duration) error
	Get([]byte) ([]byte, error)
	Has([]byte) bool
	Delete([]byte) error
}

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Delete(key []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	keyStr := string(key)
	delete(c.data, keyStr)
	return nil
}

func (c *Cache) Has(key []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	keyStr := string(key)
	_, ok := c.data[keyStr]
	return ok
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	keyStr := string(key)
	v, ok := c.data[keyStr]
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", key)
	}
	log.Printf("GET %s = %s\n", keyStr, string(v))
	return v, nil
}

func (c *Cache) Set(key, value []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[string(key)] = value

	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			delete(c.data, string(key))
		}()
	}

	c.data[string(key)] = value
	return nil
}
