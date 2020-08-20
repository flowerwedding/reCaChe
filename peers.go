/**
 * @Title  peers
 * @description  #
 * @Author  沈来
 * @Update  2020/8/20 14:37
 **/
package reCaChe

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"net/http"
	"net/url"
	"reCaChe/cachepb"
)

type PeerPicker interface {
	// PickPeer() 方法用于根据传入的 key 选择相应节点 PeerGet
	PickPeer(key string)(peer PeerGet, ok bool)
}

type PeerGet interface {
	//接口 PeerGet 的 Get() 方法用于从对应 tour 查找缓存值。PeerGet 就对应于 HTTP 客户端
	//Get(tour string, key string)(interface{}, error)
	Get(in *cachepb.Request, out *cachepb.Response) error
}

type httpGet struct {
	baseURL string//将要访问的远程节点的地址
}

func (h *httpGet) Get(in *cachepb.Request, out *cachepb.Response) error {
	u := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(in.GetGroup()), url.QueryEscape(in.GetKey()),
	)
	res, err := http.Get(u)//获取返回值
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)//转成[]bytes类型
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}

	return nil
}

var _PeerGet = (*httpGet)(nil)