/**
 * @Title  fast_test
 * @description  #
 * @Author  沈来
 * @Update  2020/8/18 17:02
 **/
package reCaChe

import (
	"fmt"
	"math/rand"
	"reCaChe/fast"
	"reCaChe/lru"
	"testing"
	"time"
)

func BenchmarkTourCacheSetParallel(b *testing.B) {
	cache := NewTourCache(nil, lru.New(b.N*100, nil))
	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			cache.Set(parallelKey(id, counter), value())
			counter ++
		}
	})
}

func BenchmarkFastCacheSetParallel(b *testing.B) {
	cache := fast.NewFastCache(b.N, 1024, nil)
	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(1000)
		counter := 0
		for pb.Next() {
			cache.Set(parallelKey(id, counter), value())
			counter ++
		}
	})
}

func key(i int) string {
	return fmt.Sprintf("key-%010d",i)
}

func value() []byte {
	return make([]byte, 100)
}

func parallelKey(threadID int, counter int) string {
	return fmt.Sprintf("key-%04d-%06d",threadID, counter)
}