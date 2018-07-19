package wechat

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func Withdraw(cliMsgid, svrMsgid, toUserName string) (msg string, err error) {
	withdraw_url := CgiUrl + "/webwxrevokemsg?pass_ticket=" + conf.PassTicket
	breq := BaseRequest{
		DeviceID: DeviceID(),
		Sid:      conf.Wxsid,
		Skey:     conf.Skey,
		Uin:      conf.Wxuin,
	}
	m := make(map[string]interface{})
	m["BaseRequest"] = breq
	m["ClientMsgId"] = cliMsgid
	m["SvrMsgId"] = svrMsgid
	m["ToUserName"] = toUserName
	jsondata, _ := json.Marshal(m)
	req, _ := http.NewRequest("post", withdraw_url, bytes.NewBuffer(jsondata))
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	wd := &WithdrawResp{}
	err = json.NewDecoder(resp.Body.(io.Reader)).Decode(wd)
	if err != nil {
		return
	}
	if wd.BaseResponse.Ret != 0 {
		msg = "撤回失败"
	}
	return
}

type WithdrawResp struct {
	BaseResponse *BaseResp
	Introduction string
	SysWording   string
}

// {
// "BaseResponse": {
// "Ret": 0,
// "ErrMsg": ""
// }
// ,
// "Introduction": "你可以撤回2分钟内发送的消息（部分旧版本微信不支持这个功能）。",
// "SysWording": "已撤回"
// }

// {
// "BaseResponse": {
// "Ret": -1,
// "ErrMsg": ""
// }
// ,
// "Introduction": "",
// "SysWording": ""
// }
