package mapcache

import (
	"container/heap"
	"encoding/json"
	"errors"
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
	h 	 SmallHeap // 小顶堆
}

type Cache interface {
	Get(key string) (string, error)
	Set(key, value string) error
	SetExpire(key string, expire time.Duration) error
	Del(key string) (string, error)
}

func (c *goMapCache) SetExpire(key string, expire time.Duration) error {
	// 1、判断在map中是否存在key，不存在返回错误
	if _, ok := c.m.Load(key); !ok {
		return errors.New("this key is not exists! ")
	}
	// 2、存在的话将过期时间插入到堆中
	heap.Push(&c.h, &ExpireDict{Key: key, Expire: expire})

	// 3、取堆顶元素，设置定时器
	go func() {
		pop := heap.Pop(&c.h).(*ExpireDict)
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


type ExpireDict struct {
	Key    string
	Expire time.Duration
	Index  int
}

func (e *ExpireDict) String() string {
	marshal, err := json.Marshal(e)
	if err != nil {
		return ""
	}
	return string(marshal)
}

type SmallHeap []*ExpireDict

func (s SmallHeap) Len() int {
	return len(s)
}
func (s SmallHeap) Less(i, j int) bool {
	return s[i].Expire < s[j].Expire
}

func (s *SmallHeap) Pop() interface{} {
	old := *s
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*s = old[0:n-1]
	return item
}

func (s *SmallHeap) Push(val interface{}) {
	n := len(*s)
	item := val.(*ExpireDict)
	item.Index = n
	*s = append(*s, item)
}

func (s SmallHeap) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
	s[i].Index = j
	s[j].Index = i
}

func (s *SmallHeap) Update(item *ExpireDict, expire int)  {
	item.Expire = time.Duration(expire)
	heap.Fix(s, item.Index)
}