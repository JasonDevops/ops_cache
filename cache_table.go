package main

import (
	"sync"
	"time"
)

var (
	delKeysChannelCap = 1000
)

type CacheTable struct {
	sync.Mutex // table leve lock

	name       string                     // table name
	items      map[interface{}]*CacheItem // store real data in the table
	createTime int64                      // write createTime when table was created
	expireMap  map[int64][]interface{}    // store set expire key
}

type delItem struct {
	keys []interface{}
	t    int64
}

// cache expire callback
// when create a new table, fun() was started
func (c *CacheTable) run(now int64) {
	t := time.NewTicker(time.Second * 1) // 获取定时器，每一秒定时器都会往管道里写数据
	defer t.Stop()

	delChan := make(chan *delItem, delKeysChannelCap) // 初始化管道

	// 监听管道，如果管道存在数据则表示有缓存过期，需要回收内存
	go func() {
		for v := range delChan {
			c.multiDelete(v.keys, v.t)
		}
	}()

	// 获取设置了缓存时间的key，如过期则将其加入到过期缓存管道中
	for {
		select {
		case <-t.C:
			now++
			c.Lock()
			if keys, ok := c.expireMap[now]; ok {
				c.Unlock()
				delChan <- &delItem{
					keys: keys,
					t:    now,
				}
			} else {
				c.Unlock()
			}
		}

	}
}

// Add add cache to in the table
func (c *CacheTable) Add(key interface{}, data interface{}, expire int64) {
	c.Lock()
	defer c.Unlock()

	// 判断key是否存在
	_, ok := c.items[key]
	if ok {
		return
	}

	// 设置过期时间，计算：createTime + expire = 缓存过期时间、0：不过期
	if expire > 0 {
		expire = time.Now().Unix() + expire
	}

	val := newCacheItem(data, expire)
	c.items[key] = val
	c.expireMap[expire] = append(c.expireMap[expire], key)
}

// Get get cache data from in the table
func (c *CacheTable) Get(key interface{}) (*CacheItem, bool) {
	c.Lock()
	defer c.Unlock()

	// check whether key was expired in the cache
	if c.checkExpireDeleteKey(key) {
		return nil, false
	}

	value, found := c.items[key]
	if !found {
		return nil, false
	}

	value.accountCount += 1
	return value, true

}

// Delete delete cache from in the table
func (c *CacheTable) Delete(key interface{}) {
	c.Lock()
	defer c.Unlock()
	delete(c.items, key)
}

// Remove remove cache from in the table
// unsafely，because this function don't start Lock() and ULock()，but that Remove() is faster than Delete()
func (c *CacheTable) Remove(key interface{}) {
	delete(c.items, key)
}

// if the key was expired, delete cache and return true
func (c *CacheTable) checkExpireDeleteKey(key interface{}) bool {

	if value, found := c.items[key]; found {
		// if expireTime !=0 and expireTime >= current time is true, delete this key from in the cache
		if value.expireTime != 0 && value.expireTime <= time.Now().Unix() {
			delete(c.items, key)
			return true
		}
	}
	return false

}

// multi delete many cache
func (c *CacheTable) multiDelete(keys []interface{}, now int64) {
	c.Lock()
	defer c.Unlock()

	for _, key := range keys {
		// delete key
		delete(c.items, key)

		// delete the current expire time
		delete(c.expireMap, now)
	}

}

// Exist return true, If a key exist in the table
func (c *CacheTable) Exist(key interface{}) bool {
	c.Lock()
	defer c.Unlock()
	_, found := c.items[key]
	return found
}

// Count count all keys number in the table
func (c *CacheTable) Count() (count int) {
	c.Lock()
	count = len(c.items)
	defer c.Unlock()
	return
}

// Flush clear all items and expireMap in the table
func (c *CacheTable) Flush() {
	c.Lock()
	defer c.Unlock()

	c.items = make(map[interface{}]*CacheItem)
	c.expireMap = make(map[int64][]interface{})
}
