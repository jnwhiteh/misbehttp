package main

import "sync"

type Counter struct {
	value int
	mutex *sync.RWMutex
}

func NewCounter(value int) *Counter {
	return &Counter{
		value: 0,
		mutex: new(sync.RWMutex),
	}
}

func (c *Counter) Increment() {
	c.mutex.Lock()
	c.value += 1
	c.mutex.Unlock()
}

func (c *Counter) Decrement() {
	c.mutex.Lock()
	c.value -= 1
	c.mutex.Unlock()
}

func (c *Counter) Get() int {
	c.mutex.RLock()
	value := c.value
	c.mutex.RUnlock()
	return value
}
