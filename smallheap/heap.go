package smallheap

import (
	"container/heap"
	"encoding/json"
	"time"
)

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