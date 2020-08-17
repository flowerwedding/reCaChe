/**
 * @Title  tour_cache
 * @description  提供给用户端使用的接口
 * @Author  沈来
 * @Update  2020/8/17 18:04
 **/
package reCaChe

type Getter interface{
	Get(key string) interface{}
}

type GetFunc func(key string) interface{}

func (f GetFunc) Get(key string) interface{}{
	return f(key)
}

type TourCache struct{
	mainCache   *safeCaChe
	getter      Getter
}

func NewTourCache(getter Getter,cache CaChe) *TourCache {
	return &TourCache{
		mainCache: newSafeCaChe(cache),
		getter:    getter,
	}
}

func (t *TourCache) Get(key string) interface{} {
	val := t.mainCache.get(key)
	if val != nil {
		return val
	}

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

func (t *TourCache) Stat() *Stat {
	return t.mainCache.stat()
}