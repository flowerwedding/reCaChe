/**
 * @Title  http
 * @description  节点间通信
 * @Author  沈来
 * @Update  2020/8/19 16:21
 **/
package reCaChe

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net/http"
	pb "reCaChe/cachepb"
	"reCaChe/consistent"
	"strings"
	"sync"
)

const defaultBasePath = "/_reCache/"
const defaultReplicas = 50


type HTTPPool struct {
	self      string//记录自己的地址，IP和端口
	basePath  string//节点间通讯地址的前缀，默认是...
	mu          sync.Mutex
	peers       *consistent.Map//根据具体的key选择节点
	httpGetters map[string]*httpGet//映射远程节点httpGet
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

	body, err := proto.Marshal(&pb.Response{Value: []byte(view.(string))})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//使用 w.Write() 将缓存值作为 httpResponse 的 body 返回
	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(body)
}

//Set() 方法实例化了一致性哈希算法，并且添加了传入的节点
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistent.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGet, len(peers))
	//并为每一个节点创建了一个 HTTP 客户端 httpGetter
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGet{baseURL: peer + p.basePath}
	}
}

//根据具体的 key，选择节点，返回节点对应的 HTTP 客户端
func (p *HTTPPool) PickPeer(key string) (*httpGet, bool){
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s",peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _PeerPicker = (*HTTPPool)(nil)
