/**
 * @Title  singleflight
 * @description  #
 * @Author  沈来
 * @Update  2020/8/20 15:59
 **/
package singleflight

import "sync"

//call 代表正在进行中，或已经结束的请求。使用 sync.WaitGroup 锁避免重入
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

//Group 是 singleflight 的主数据结构，管理不同 key 的请求(call)
type Group struct {
	mu sync.Mutex       // protects m
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	if c, ok := g.m[key]; ok {
		c.wg.Wait()   // 如果请求正在进行中，则等待
		return c.val, c.err  // 请求结束，返回结果
	}
	c := new(call)
	c.wg.Add(1)       // 发起请求前加锁
	g.m[key] = c      // 添加到 g.m，表明 key 已经有对应的请求在处理

	c.val, c.err = fn() // 调用 fn，发起请求
	c.wg.Done()         // 请求结束

	delete(g.m, key)    // 更新 g.m

	return c.val, c.err // 返回结果
}