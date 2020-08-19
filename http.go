/**
 * @Title  http
 * @description  节点间通信
 * @Author  沈来
 * @Update  2020/8/19 16:21
 **/
package reCaChe

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_reCache/"

type HTTPPool struct {
	self      string//记录自己的地址，IP和端口
	basePath  string//节点间通讯地址的前缀，默认是...
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:      self,
		basePath:  defaultBasePath,
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}){
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format,v ...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter,r *http.Request){
	///<basepath>/<tourname>/<key>路径格式
	//前缀是否正确
	if !strings.HasPrefix(r.URL.Path, p.basePath){
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s",r.Method, r.URL.Path)

	parts := strings.SplitN(r.URL.Path[len(p.basePath):],"/",2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	//获得实例
	tourName := parts[0]
	key := parts[1]

	tour := GetTour(tourName)
	if tour == nil {
		http.Error(w,"no such tour: " + tourName, http.StatusNotFound)
		return
	}

	//获取缓存数据
	view := tour.Get(key)
	if view == nil {
		http.Error(w,"",http.StatusInternalServerError)
	}

	//使用 w.Write() 将缓存值作为 httpResponse 的 body 返回
	w.Header().Set("Content-Type","application/octet-stream")
	_, _ = w.Write([]byte(view.(string)))
}