package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var (
	k = "test_" // test key prefix
)

func TestCache(t *testing.T) {
	t1 := Cache("t1")
	t1.Add("k1", "v1", 0)
	t1.Add("k2", "v2", 0)
	t1.Add("k3", "v3", 10)
	fmt.Println(t1.Get("k2"))
	fmt.Println(t1.Get("k3"))
	fmt.Println(t1.Get("k4"))

	fmt.Println(t1.items)

	time.Sleep(15 * time.Second)
	fmt.Println(t1.Get("k3"))

}

func TestExpireKeys(t *testing.T) {
	t1 := Cache("expire_t1")
	//expireTime :=  5
	for i := 0; i < 10; i++ {
		key := k + strconv.Itoa(i)
		t1.Add(key, i, int64(i))
	}
	//fmt.Println(t1.items)

	fmt.Println("before..........")
	fmt.Println(t1.items)
	fmt.Println(t1.expireMap)
	time.Sleep(15 * time.Second)

	fmt.Println("after.....")
	fmt.Println(t1.items)
	fmt.Println(t1.expireMap)
}

func TestAccessCount(t *testing.T) {
	t1 := Cache("access_t1")
	t1.Add("k1", "v1", 0)
	fmt.Println(t1.GetAccessCount("k1"))
	t1.Get("k1")
	t1.Get("k1")
	fmt.Println(t1.GetAccessCount("k1"))
}
