/**
 * @Title  lru
 * @description  LRU算法淘汰最近最近访问频率最低的记录
 * @Author  沈来
 * @Update  2020/8/17 16:29
 **/
package lru

import (
	"container/list"
	"reCaChe"
)

type lru struct {
	maxBytes int

	onEvicted func(key string, value interface{})

	usedBytes int

	ll    *list.List
	cache map[string]*list.Element
}

type entry struct {
	key   string
	value interface{}
}

func (e *entry) Len() int {
	return reCaChe.CalcLen(e.value)
}

func New(maxBytes int, onEvicted func(key string, value interface{})) reCaChe.CaChe {
	return &lru{
		maxBytes:  maxBytes,
		onEvicted: onEvicted,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
	}
}

func (l *lru) Set(key string, value interface{}) {
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e)
		en := e.Value.(*entry)
		l.usedBytes = l.usedBytes - reCaChe.CalcLen(en.value) + reCaChe.CalcLen(value)
		en.value = value
		return
	}

	en := &entry{key, value}
	e := l.ll.PushBack(en)
	l.cache[key] = e

	l.usedBytes += en.Len()
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.DelOldest()
	}
}

func (l *lru) Get(key string) interface{}{
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e)
		return e.Value.(*entry).value
	}

	return nil
}

func (l *lru) Del(key string) {
	if e, ok := l.cache[key]; ok {
		l.removeElement(e)
	}
}

func (l *lru) DelOldest() {
	l.removeElement(l.ll.Front())
}

func (l *lru) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	l.ll.Remove(e)
	en := e.Value.(*entry)
	l.usedBytes -= en.Len()
	delete(l.cache, en.key)

	if l.onEvicted != nil {
		l.onEvicted(en.key, en.value)
	}
}

func (l *lru) Len() int {
	return l.ll.Len()
}