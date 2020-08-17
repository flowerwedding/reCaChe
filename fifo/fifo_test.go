/**
 * @Title  fifo_test
 * @description  测试
 * @Author  沈来
 * @Update  2020/8/17 15:08
 **/
package fifo

import (
	"github.com/matryer/is"
	"testing"
)

func TestSetGet(t *testing.T) {
	its := is.New(t)

	cache := New(24,nil)
	cache.DelOldest()
	cache.Set("k1",1)
	v := cache.Get("k1")
	its.Equal(v,1)

	cache.Del("k1")
	its.Equal(0,cache.Len()) // expect to be the same
}

func TestOnEvicted(t *testing.T) {
	its := is.New(t)

	keys := make([]string, 0, 8)
	onEvicted := func(key string, value interface{}) {
		keys = append(keys, key)
	}
	cache := New(16, onEvicted)

	cache.Set("k1",1)
	cache.Set("k2",2)
	cache.Get("k1")
	cache.Set("k3",3)
	cache.Get("k1")
	cache.Set("k4",4)

	expected := []string{"k1", "k2"}

	its.Equal(expected, keys)
	its.Equal(2,cache.Len())
}