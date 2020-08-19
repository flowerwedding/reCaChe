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
	maxBytes int//maxBytes 是允许使用的最大内存

	onEvicted func(key string, value interface{})//OnEvicted 是某条记录被移除时的回调函数，可以为 nil

	usedBytes int//usedBytes 是当前已使用的内存

	ll    *list.List
	cache map[string]*list.Element//键是字符串，值是双向链表中对应节点的指针
}

//键值对 entry 是双向链表节点的数据类型，在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射
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
	//如果键存在，则更新对应节点的值，并将该节点移到队尾
	if e, ok := l.cache[key]; ok {
		l.ll.MoveToBack(e)
		en := e.Value.(*entry)
		l.usedBytes = l.usedBytes - reCaChe.CalcLen(en.value) + reCaChe.CalcLen(value)
		en.value = value
		return
	}

	//不存在则是新增场景，首先队尾添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关系
	en := &entry{key, value}
	e := l.ll.PushBack(en)
	l.cache[key] = e

	//更新 l.usedBytes，如果超过了设定的最大值 l.maxBytes，则移除最少访问的节点
	l.usedBytes += en.Len()
	if l.maxBytes > 0 && l.usedBytes > l.maxBytes {
		l.DelOldest()
	}
}

func (l *lru) Get(key string) interface{}{
	//第一步从字典中找到对应的双向链表的节点，第二步，将该节点移动到队尾
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
	l.removeElement(l.ll.Front())//c.ll.Front() 取到队首节点
}

func (l *lru) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	l.ll.Remove(e)//从链表中删除
	en := e.Value.(*entry)
	l.usedBytes -= en.Len()//更新当前所占用的内存
	delete(l.cache, en.key)//从字典中 l.cache 删除该节点的映射关系

	if l.onEvicted != nil {//使用回调函数
		l.onEvicted(en.key, en.value)
	}
}

//添加的数据的数量
func (l *lru) Len() int {
	return l.ll.Len()
}