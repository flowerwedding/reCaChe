/**
 * @Title  main
 * @description  #
 * @Author  沈来
 * @Update  2020/8/19 16:42
 **/
package reCaChe

import (
	"log"
	"net/http"
	"reCaChe/lru"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestNewHTTPPool(t *testing.T) {
	getter := GetFunc(func(key string) interface{}{
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return v
		}
		return nil
	})
	NewTourCache("scores", getter, lru.New(0, nil))

	addr := "localhost:9999"
	peers := NewHTTPPool(addr)
	log.Println("reCache is running at", addr)
	log.Fatal(http.ListenAndServe(addr,peers))
}