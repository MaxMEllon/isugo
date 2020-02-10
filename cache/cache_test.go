package cache

import (
	"testing"
)

type Object struct {
	Key   string
	Value string
}

func TestCacheAsSyncMap(t *testing.T) {
	cache, err := NewCache("sync.Map", Option{})

	if err != nil {
		t.Errorf("failure create cache instance %+v", err)
		return
	}

	cache.Store("example", Object{
		Key:   "case1",
		Value: "any string",
	})

	var result Object
	err = cache.Load("example", &result)

	if err != nil {
		t.Errorf("failure fetch data from cache %+v", err)
		return
	}

	if result.Key != "case1" || result.Value != "any string" {
		t.Errorf("assertion failure %+v", result)
		return
	}
}

func TestCacheAsRedis(t *testing.T) {
	cache, err := NewCache("redis", Option{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	if err != nil {
		t.Errorf("failure create cache instance %+v", err)
		return
	}

	cache.Store("example", Object{
		Key:   "case2",
		Value: "any string",
	})

	var result Object
	err = cache.Load("example", &result)

	if err != nil {
		t.Errorf("failure fetch data from cache %+v", err)
		return
	}

	if result.Key != "case2" || result.Value != "any string" {
		t.Errorf("assertion failure %+v", result)
		return
	}
}
