package idiom

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// 异常码（code）	说明
// 40001	参数key错误
// 40002	请求内容info为空
// 40004	当天请求次数已使用完
// 40007	数据格式异常/请按规定的要求进行加密

func Chengyujielong(str, userName string) (rsp *IdiomResp, err error) {
	req := IdiomReq{
		Key:    conf.Apikey,
		Info:   str,
		Userid: userName,
	}
	d, _ := json.Marshal(req)
	resp, err := http.Post("http://www.tuling123.com/openapi/api", "application/json; charset=UTF-8",
		bytes.NewReader(d))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	rsp = &IdiomResp{}
	err = json.NewDecoder(resp.Body).Decode(rsp)
	if err != nil {
		rsp = nil
		return
	}
	return
}

type IdiomReq struct {
	Key    string `json:"key"`
	Info   string `json:"info"`
	Userid string `json:"userid"`
}

type IdiomResp struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func init() {
	b, _ := ioutil.ReadFile(strings.Split(os.Getenv("GOPATH"), ";")[0] + "/src/test/wechat/kit/idiom/api.conf") // ./api.conf
	json.Unmarshal(b, &conf)
}

var conf Conf

type Conf struct {
	Apikey string
}
