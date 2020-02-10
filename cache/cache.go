package cache

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/cache/v7"
	"github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack/v4"
)

type store func(key string, value interface{}) error
type load func(key string, value interface{}) error

type client struct {
	Store store
	Load  load
}

type Option struct {
	Addr     string
	Password string
	DB       int
}

var codec *cache.Codec
var once sync.Once

func NewCache(kind string, option Option) (*client, error) {
	c := &client{}
	switch kind {
	case "redis":
		once.Do(func() {
			r := redis.NewClient(&redis.Options{
				Addr:     option.Addr,
				Password: option.Password,
				DB:       option.DB,
			})
			codec = &cache.Codec{
				Redis: r,
				Marshal: func(v interface{}) ([]byte, error) {
					return msgpack.Marshal(v)
				},
				Unmarshal: func(b []byte, v interface{}) error {
					return msgpack.Unmarshal(b, v)
				},
			}
		})
		c.Store = func(key string, value interface{}) error {
			err := codec.Set(&cache.Item{
				Key:        key,
				Object:     value,
				Expiration: time.Hour * 4,
			})
			if err != nil {
				return err
			}
			return nil
		}
		c.Load = func(key string, value interface{}) error {
			err := codec.Get(key, value)
			if err != nil {
				return err
			}
			return nil
		}
		return c, nil
	case "sync.Map":
		var cache sync.Map
		c.Store = func(key string, value interface{}) error {
			v, err := msgpack.Marshal(value)
			if err != nil {
				return err
			}
			cache.Store(key, v)
			return nil
		}
		c.Load = func(key string, value interface{}) error {
			v, _ := cache.Load(key)
			if v == nil {
				return nil
			}
			b, ok := v.([]byte)
			if !ok {
				return fmt.Errorf("failure type assertion (key: %s)", key)
			}
			err := msgpack.Unmarshal(b, value)
			if err != nil {
				return err
			}
			return nil
		}
		return c, nil
	}
	return nil, errors.New(fmt.Sprintf("unexpected kind: %s", kind))
}
