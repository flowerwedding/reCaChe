/**
 * @Title  tour_cache
 * @description  提供给用户端使用的接口，二次封装
 * @Author  沈来
 * @Update  2020/8/17 18:04
 **/
package reCaChe

import "sync"

type Getter interface{
	Get(key string) interface{}
}

type GetFunc func(key string) interface{}

//实现Getter接口的Get方法
func (f GetFunc) Get(key string) interface{}{
	return f(key)
}

var (
	//缓存组，不同的缓存空间有不同的名字
	tours = make(map[string]*TourCache)
	mu sync.RWMutex
)

type TourCache struct{
	name        string
	mainCache   *safeCaChe//并发安全的缓存实现
	getter      Getter//回调，用于缓存未命中时从数据源获取的数据
}

//如果放在第一次set里面判断nil，可以延迟初始化，建设内存
func NewTourCache(name string,getter Getter,cache CaChe) *TourCache {
	mu.Lock()
	defer mu.Unlock()
	tour := &TourCache{
		name:      name,
		mainCache: newSafeCaChe(cache),
		getter:    getter,
	}
	tours[name] = tour
	return tour
}

func GetTour(name string) *TourCache{
	mu.RLock()
	defer mu.RUnlock()
	tour := tours[name]
	return tour
}

func (t *TourCache) Get(key string) interface{} {
	val := t.mainCache.get(key)
	if val != nil {
		return val
	}

	//缓存中数据不存在，调用回调函数获取数据并将数据写入缓存，最后返回获取的数据
	if t.getter != nil {
		val = t.getter.Get(key)
		if val == nil {
			return nil
		}
		t.mainCache.set(key, val)
		return val
	}

	return nil
}

func (t *TourCache) Set(key string, val interface{}) {
	if val == nil {
		return
	}
	t.mainCache.set(key,val)
}

func (t *TourCache) Stat() *Stat {
	return t.mainCache.stat()
}