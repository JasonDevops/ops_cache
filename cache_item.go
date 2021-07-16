package main

import (
	"sync"
	"time"
)

type CacheItem struct {
	data         interface{}
	expireTime   int64 // 0：no-expire
	accountCount int64 // every time cache was  accessed to add 1

	accessTime int64
	createTime int64
	sync.Mutex // item level lock
}

func newCacheItem(data interface{}, expireTime int64) *CacheItem {
	nowTime := time.Now().Unix()
	return &CacheItem{
		data:       data,
		expireTime: expireTime,
		createTime: nowTime,
		accessTime: nowTime,
	}
}

func (c *CacheItem) Data() interface{} {
	return c.data
}

func (c *CacheItem) AccessCount() int64 {
	c.Lock()
	defer c.Unlock()
	return c.accountCount
}
