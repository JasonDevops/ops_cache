package main

import (
	"sync"
	"time"
)

var (
	cache map[string]*CacheTable
	lock  sync.Mutex
)

// Cache Create a table
func Cache(name string) *CacheTable {
	lock.Lock()
	t, ok := cache[name] // 获取指定的缓存表
	lock.Unlock()

	if !ok {
		lock.Lock()
		defer lock.Unlock()
		t = &CacheTable{
			name:       name,
			items:      make(map[interface{}]*CacheItem),
			createTime: time.Now().Unix(),
			Mutex:      sync.Mutex{},
			expireMap:  make(map[int64][]interface{}),
		}
		go t.run(t.createTime)
	}
	return t
}
