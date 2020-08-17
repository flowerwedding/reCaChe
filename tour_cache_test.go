/**
 * @Title  tour_cache_test
 * @description  测试
 * @Author  沈来
 * @Update  2020/8/17 18:12
 **/
package reCaChe

import (
	"github.com/matryer/is"
	"log"
	"reCaChe/lru"
	"sync"
	"testing"
)

//放在别的路径下测试，否则会循环导包
func TestTourCacheGet(t *testing.T) {
	db := map[string]string{
		"key1":"val1",
		"key2":"val2",
		"key3":"val3",
		"key4":"val4",
	}
	getter := GetFunc(func(key string) interface{}{
		log.Println("[From DB] find key",key)

		if val, ok := db[key]; ok {
			return val
		}
		return nil
	})
	tourCache := NewTourCache(getter, lru.New(0,nil))

	its := is.New(t)

	var wg sync.WaitGroup

	for k, v := range db {
		wg.Add(1)
		go func(k, v string){
			defer wg.Done()
			its.Equal(tourCache.Get(k), v)

			its.Equal(tourCache.Get(k), v)
		}(k, v)
	}
	wg.Wait()

	its.Equal(tourCache.Get("unknown"), nil)
	its.Equal(tourCache.Get("unknown"), nil)

	its.Equal(tourCache.Stat().NGet, 10)
	its.Equal(tourCache.Stat().NHit, 4)
}