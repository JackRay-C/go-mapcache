package smallheap

import (
	"container/heap"
	"fmt"
	"testing"
	"time"
)

func TestInit0(t *testing.T) {
	items := map[string]int{
		"renhj:asdfasdfad": 1, "renhj:asdf1324ssdf": 2,
		"renhj:adsafdsfadsfa": 3, "renhj:asdfadsfadsfadfsd": 4,
		"admin:adsfadsfasdf": 1, "admin:asdfasdfadsfa": 2,
	}
	p := make(SmallHeap, len(items))
	i := 0
	for key, expire := range items {
		p[i] = &ExpireDict{
			Key:    key,
			Expire: time.Duration(expire),
			Index:  i,
		}
		i++
	}
	heap.Init(&p)

	e := &ExpireDict{
		Key:    "renhj:adsfa42ubfanj",
		Expire: 20,
	}
	heap.Push(&p, e)
	p.Update(e, 20)

	for p.Len() > 0 {
		dict := heap.Pop(&p).(*ExpireDict)
		fmt.Printf("%s\n", dict)
	}

	// Output:
	// admin:adsfadsfasdf:1
	// renhj:asdfasdfad:1
	// renhj:asdf1324ssdf:2
	// admin:asdfasdfadsfa:2
	// renhj:adsafdsfadsfa:3
	// renhj:adsfa42ubfanj:3
	// renhj:asdfadsfadsfadfsd:4

}
