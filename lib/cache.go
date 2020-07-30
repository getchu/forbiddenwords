package lib

import (
	"sync"
)

type Cache struct {
	dataList map[string]interface{}
	lock     *sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		dataList: make(map[string]interface{}),
		lock:     new(sync.RWMutex),
	}
}

func (c *Cache) Set(key string, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.dataList[key] = value
}

func (c *Cache) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if value, ok := c.dataList[key]; ok {
		return value
	}
	return nil
}
