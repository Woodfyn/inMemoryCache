package cache

import (
	"errors"
	"sync"
	"time"
)

type Cache struct {
	items map[int]struct {
		value      interface{}
		expiration time.Duration
	}
	jobChannel chan int
	mutex      sync.RWMutex
}

func NewCache(workerCount int) *Cache {
	c := &Cache{
		items: make(map[int]struct {
			value      interface{}
			expiration time.Duration
		}),
		jobChannel: make(chan int),
	}

	for i := 0; i < workerCount; i++ {
		go c.worker()
	}

	return c
}

func (c *Cache) Set(key int, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.items[key] = struct {
		value      interface{}
		expiration time.Duration
	}{value, ttl}

	c.jobChannel <- key
}

func (c *Cache) worker() {
	for key := range c.jobChannel {
		go c.expireAfter(key)
	}
}

func (c *Cache) expireAfter(key int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ttl := c.items[key].expiration

	<-time.After(ttl)
	c.Delete(key)
}

func (c *Cache) Get(key int) (interface{}, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, errors.New("KEY_NOT_FOUND")
	}

	return item.value, nil
}

func (c *Cache) Update(key int, newValue interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, ok := c.items[key]
	if !ok {
		return errors.New("KEY_NOT_FOUND")
	}

	c.items[key] = struct {
		value      interface{}
		expiration time.Duration
	}{newValue, c.items[key].expiration}

	return nil
}

func (c *Cache) Delete(key int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, ok := c.items[key]
	if !ok {
		return errors.New("KEY_NOT_FOUND")
	}

	delete(c.items, key)
	return nil
}

func (c *Cache) Close() {
	close(c.jobChannel)
}
