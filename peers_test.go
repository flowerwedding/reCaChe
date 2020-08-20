/**
 * @Title  peers_test
 * @description  #
 * @Author  沈来
 * @Update  2020/8/20 15:36
 **/
package reCaChe

import (
	"flag"
	"log"
	"net/http"
	"reCaChe/lru"
	"testing"
)

var db2 = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *TourCache {
	getter := GetFunc(func(key string) interface{} {
		log.Println("[Slowdb2] search key", key)
		if v, ok := db2[key]; ok {
			return []byte(v)
		}
		return nil
	})
	return NewTourCache("scores", getter, lru.New(0, nil))
}

func startCacheServer(addr string, addrs []string, gee *TourCache) {
	peers := NewHTTPPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("cache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, gee *TourCache) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view := gee.Get(key)
			if view == nil {
				http.Error(w,"", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write([]byte(view.(string)))
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func TestHTTPPool(t *testing.T) {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	gee := createGroup()
	if api {
		go startAPIServer(apiAddr, gee)
	}
	startCacheServer(addrMap[port], []string(addrs), gee)
}