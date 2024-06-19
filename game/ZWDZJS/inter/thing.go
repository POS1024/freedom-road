package inter

import (
	"container/list"
	"github.com/hajimehoshi/ebiten/v2"
	"sync"
)

type Thing interface {
	GetKey() int64
	SetKey(int64)
	Update() error
	Draw(screen *ebiten.Image)
	GetLayers() int
}

type LFUCache struct {
	keyCount int64
	entries  map[int64]*list.Element
	freqs    map[int64]*list.List
	mu       sync.RWMutex
	gomu     sync.RWMutex
}

type entry struct {
	key   int64
	value Thing
	freq  int64
}

func NewLFUCache() *LFUCache {
	return &LFUCache{
		entries: make(map[int64]*list.Element),
		freqs:   make(map[int64]*list.List),
	}
}

func (c *LFUCache) Get(key int64) Thing {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if elem, ok := c.entries[key]; ok {
		c.increment(elem)
		return elem.Value.(*entry).value
	}
	return nil
}

func (c *LFUCache) Put(value Thing) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.keyCount++
	key := c.keyCount
	value.SetKey(key)
	if elem, ok := c.entries[key]; ok {
		c.increment(elem)
		elem.Value.(*entry).value = value
	} else {
		ent := &entry{key: key, value: value, freq: 1}
		if c.freqs[1] == nil {
			c.freqs[1] = list.New()
		}
		c.entries[key] = c.freqs[1].PushFront(ent)
	}
	return key
}

func (c *LFUCache) GoPut(value Thing) int64 {
	c.gomu.Lock()
	defer c.gomu.Unlock()
	c.keyCount++
	key := c.keyCount
	value.SetKey(key)
	if elem, ok := c.entries[key]; ok {
		c.increment(elem)
		elem.Value.(*entry).value = value
	} else {
		ent := &entry{key: key, value: value, freq: 1}
		if c.freqs[1] == nil {
			c.freqs[1] = list.New()
		}
		c.entries[key] = c.freqs[1].PushFront(ent)
	}
	return key
}

func (c *LFUCache) increment(elem *list.Element) {
	ent := elem.Value.(*entry)
	c.freqs[ent.freq].Remove(elem)
	if c.freqs[ent.freq].Len() == 0 {
		delete(c.freqs, ent.freq)
	}

	ent.freq++
	if c.freqs[ent.freq] == nil {
		c.freqs[ent.freq] = list.New()
	}
	c.entries[ent.key] = c.freqs[ent.freq].PushFront(ent)
}

func (c *LFUCache) Delete(key int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.entries[key]
	if !ok {
		return
	}

	ent := elem.Value.(*entry)
	c.freqs[ent.freq].Remove(elem)
	if c.freqs[ent.freq].Len() == 0 {
		delete(c.freqs, ent.freq)
	}
	delete(c.entries, key)
}

func (c *LFUCache) Range(f func(key int64, value Thing) bool) {
	c.gomu.RLock()
	defer c.gomu.RUnlock()

	for _, elem := range c.entries {
		ent := elem.Value.(*entry)
		if !f(ent.key, ent.value) {
			break
		}
	}
}

var GameThings = NewLFUCache()
