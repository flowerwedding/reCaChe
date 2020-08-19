/**
 * @Title  cache
 * @description  加锁，并发安全
 * @Author  沈来
 * @Update  2020/8/17 14:32
 **/
package reCaChe

import (
	"log"
	"sync"
)

type CaChe interface{
	Set (key string, value interface{})
	Get(key string) interface{}
	Del(key string)
	DelOldest()
	Len() int
}

const DefaultMaxBytes = 1 << 29

type safeCaChe struct {
	m          sync.RWMutex
	cache      CaChe
	nhit, nget int
}

func newSafeCaChe(cache CaChe) *safeCaChe{
	return &safeCaChe{
		cache : cache,
	}
}

func (sc *safeCaChe) set(key string, value interface{}){
	sc.m.Lock()
	defer sc.m.Unlock()
	sc.cache.Set(key, value)
}

func (sc *safeCaChe) get(key string) interface{}{
	sc.m.RLock()
	defer sc.m.RUnlock()
	sc.nget++
	if sc.cache == nil {
		return nil
	}

	v := sc.cache.Get(key)
	if v != nil {
		log.Println("[TourCache] hit")
		sc.nhit++
	}

	return v
}

func (sc *safeCaChe) stat() *Stat{
	sc.m.RLock()
	defer sc.m.RUnlock()
	return &Stat{
		NHit: sc.nhit,
		NGet: sc.nget,
	}
}

type Stat struct{
	NHit, NGet int
}