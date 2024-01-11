package cache

import (
	"errors"
	"fmt"
	"time"
)

type Cache struct {
	items map[string]struct {
		value      interface{}
		expiration time.Duration
	}
	jobChannel chan string
}

func NewCache(workerCount int) *Cache {
	c := &Cache{
		items: make(map[string]struct {
			value      interface{}
			expiration time.Duration
		}),
		jobChannel: make(chan string),
	}

	for i := 0; i < workerCount; i++ {
		go c.worker()
	}

	return c
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
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

func (c *Cache) expireAfter(key string) {
	ttl := c.items[key].expiration

	<-time.After(ttl)
	c.Delete(key)
}

func (c *Cache) Get(key string) error {
	item, ok := c.items[key]
	if !ok {
		err := errors.New("KEY_NOT_FOUND")
		fmt.Println(err)
		return err
	}
	fmt.Println(item.value)
	return nil
}

func (c *Cache) Update(key string, newValue interface{}) error {
	_, ok := c.items[key]
	if !ok {
		err := errors.New("KEY_NOT_FOUND")
		fmt.Println(err)
		return err
	}

	c.items[key] = struct {
		value      interface{}
		expiration time.Duration
	}{newValue, c.items[key].expiration}

	return nil
}

func (c *Cache) Delete(key string) error {
	_, ok := c.items[key]
	if !ok {
		err := errors.New("KEY_NOT_FOUND")
		fmt.Println(err)
		return err
	}
	delete(c.items, key)

	return nil
}
