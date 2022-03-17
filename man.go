package mapcache

import (
	"container/heap"
	"errors"
	"mapcache/smallheap"
	"sync"
	"time"
)

func NewCache() Cache {
	g := &goMapCache{}
	heap.Init(&g.h)
	return g
}

type goMapCache struct {
	m    sync.Map       // 键值对
	h 	 smallheap.SmallHeap // 小顶堆
}

type Cache interface {
	Get(key string) (string, error)
	Set(key, value string) error
	SetExpire(key string, expire int) error
	Del(key string) (string, error)
}

func (c *goMapCache) SetExpire(key string, expire int) error {
	// 1、判断在map中是否存在key，不存在返回错误
	if _, ok := c.m.Load(key); !ok {
		return errors.New("this key is not exists! ")
	}
	// 2、存在的话将过期时间插入到堆中
	heap.Push(&c.h, &smallheap.ExpireDict{Key: key, Expire: time.Duration(expire)})

	// 3、取堆顶元素，设置定时器
	go func() {
		pop := heap.Pop(&c.h).(*smallheap.ExpireDict)
		time.AfterFunc(pop.Expire, func() {
			c.m.LoadAndDelete(pop.Key)
		})
	}()
	return nil
}

func (c *goMapCache) Get(key string) (string, error) {
	load, ok := c.m.Load(key)
	if ok {
		return load.(string), nil
	}
	return "", errors.New("not exists! ")
}

func (c *goMapCache) Set(key, value string) error {
	store, loaded := c.m.LoadOrStore(key, value)
	if loaded && store != value {
		return errors.New("this key is stored！")
	}
	return nil
}

func (c *goMapCache) Del(key string) (string, error) {
	_, loaded := c.m.LoadAndDelete(key)
	if !loaded {
		return "", errors.New("this key isn't exists! ")
	}
	return "ok", nil
}
