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
	t1.Add("k3", "v3", 5)
	k1, _ := t1.Get("k1")
	k2, _ := t1.Get("k2")
	k3, _ := t1.Get("k3")

	fmt.Println(k1.Data())
	fmt.Println(k2.Data())
	fmt.Println(k3.Data())

	// Delete()
	t1.Delete("k1")

	time.Sleep(6 * time.Second)
	k1, _ = t1.Get("k1")
	k3, _ = t1.Get("k3")

	fmt.Println(k1)
	fmt.Println(k3)

	// Exist()
	fmt.Println(t1.Exist("k1"))
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
	k1, ok := t1.Get("k1")
	if ok {
		fmt.Println(k1.AccessCount())
	}

}

func TestTableCount(t *testing.T) {
	t1 := Cache("count_t1")
	count := 100
	for i := 0; i < count; i++ {
		key := k + strconv.Itoa(i)
		t1.Add(key, i, int64(i))
	}
	fmt.Println(t1.Count())
}

func TestTableFlush(t *testing.T) {
	t1 := Cache("flush_t1")
	count := 10000
	for i := 0; i < count; i++ {
		key := k + strconv.Itoa(i)
		t1.Add(key, i, int64(i))
	}
	fmt.Println(t1.Count())
	// Flush()
	t1.Flush()
	fmt.Println(t1.Count())

}
