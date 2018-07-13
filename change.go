package wechat

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

type StatusNotifyResp struct {
	BaseResponse BaseResp `json:"BaseResponse"`
	MsgID        string   `json:"MsgID"`
}

// 切换当前对话框 第一次是当前账号到当前账号 WebInitConf.User.UserName
func WebWxStatusNotify(from, to string) (rsp *StatusNotifyResp, err error) {
	xm := url.Values{}
	xm.Add("pass_ticket", conf.PassTicket)
	statusnotyfy_url := CgiUrl + "/webwxstatusnotify?" + xm.Encode()
	breq := BaseRequest{
		conf.Wxuin,
		conf.Wxsid,
		conf.Skey,
		DeviceID(),
	}
	m := make(map[string]interface{})
	m["BaseRequest"] = breq
	m["ClientMsgId"] = time.Now().UnixNano() / 1e6
	m["Code"] = 3
	m["FromUserName"] = from // WebInitConf.User.UserName
	m["ToUserName"] = to     // WebInitConf.User.UserName
	breqdata, _ := json.Marshal(m)
	req, _ := http.NewRequest("post", statusnotyfy_url, bytes.NewBuffer(breqdata))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	rsp = &StatusNotifyResp{}
	err = json.NewDecoder(resp.Body.(io.Reader)).Decode(rsp)
	if err != nil {
		return
	}
	return
}
