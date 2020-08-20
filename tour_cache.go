/**
 * @Title  tour_cache
 * @description  提供给用户端使用的接口，二次封装
 * @Author  沈来
 * @Update  2020/8/17 18:04
 **/
package reCaChe

import (
	"log"
	"reCaChe/cachepb"
	"reCaChe/singleflight"
	"sync"
)

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
	peers       PeerPicker
	loader      *singleflight.Group
}

//如果放在第一次set里面判断nil，可以延迟初始化，建设内存
func NewTourCache(name string,getter Getter,cache CaChe) *TourCache {
	mu.Lock()
	defer mu.Unlock()
	tour := &TourCache{
		name:      name,
		mainCache: newSafeCaChe(cache),
		getter:    getter,
		loader:    &singleflight.Group{},
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

func (t *TourCache) RegisterPeers(peers PeerPicker) {
	if t.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	t.peers = peers
}

func (t *TourCache) load(key string) (value interface{}, err error){
	viewi, err := t.loader.Do(key, func() (interface{}, error) {
		if t.peers != nil {
			if peer, ok := t.peers.PickPeer(key); ok {
				if value, err  = t.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[cache] Failed to get from peer", err)
			}
		}

		return t.getLocally(key)
	})

	if err == nil {
		return viewi, nil
	}
	return
}

func (t *TourCache) getFromPeer(peer PeerGet, key string) (interface{}, error) {
	req := &cachepb.Request{
		Group: t.name,
		Key:   key,
	}
	res := &cachepb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return "", err
	}
	return res.Value, nil
}

func (t *TourCache) getLocally(key string) (interface{}, error) {
	bytes := t.getter.Get(key)
	if bytes == nil {
		return nil, nil
	}
	t.populateCache(key, bytes)
	return bytes, nil
}

func (t *TourCache) populateCache(key string, value interface{}) {
	t.mainCache.set(key, value)
}