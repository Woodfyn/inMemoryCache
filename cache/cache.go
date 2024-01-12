package cache

import (
	"errors"
	"fmt"
	"time"
)

type Cache struct {
	items map[int]struct {
		value      interface{}
		expiration time.Duration
	}
	jobChannel chan int
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
	ttl := c.items[key].expiration

	<-time.After(ttl)
	c.Delete(key)
}

func (c *Cache) Get(key int) (interface{}, error) {
	item, ok := c.items[key]
	if !ok {
		err := errors.New("KEY_NOT_FOUND")
		return nil, err
	}
	return item.value, nil
}

func (c *Cache) Update(key int, newValue interface{}) error {
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

func (c *Cache) Delete(key int) error {
	_, ok := c.items[key]
	if !ok {
		err := errors.New("KEY_NOT_FOUND")
		fmt.Println(err)
		return err
	}
	delete(c.items, key)

	return nil
}
