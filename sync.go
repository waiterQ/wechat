package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	lerrors "test/wechat/errors"
)

var (
	URLPool = []UrlGroup{
		{"wx2.qq.com", "file.wx2.qq.com", "webpush.wx2.qq.com"},
		{"wx8.qq.com", "file.wx8.qq.com", "webpush.wx8.qq.com"},
		{"qq.com", "file.wx.qq.com", "webpush.wx.qq.com"},
		{"web2.wechat.com", "file.web2.wechat.com", "webpush.web2.wechat.com"},
		{"wechat.com", "file.web.wechat.com", "webpush.web.wechat.com"},
	}
)

type UrlGroup struct {
	IndexUrl  string
	UploadUrl string
	SyncUrl   string
}

var syncUrl string
var syncHost string
var sync_cookies []*http.Cookie

// func FirstSyncCookies() {
// 	u, _ := url.Parse(syncUrl)
// 	Cli.Jar.SetCookies(u, cookies)
// 	return
// }

func SetSyncCookies() {
	u, _ := url.Parse(syncUrl)
	// for i := 0; i < len(sync_cookies); i++ {
	// 	fmt.Println(*sync_cookies[i])
	// }
	Cli.Jar.SetCookies(u, sync_cookies)
	return
}

// sync的逻辑好像是get.syncCheck.select=2 则需要post.webwxsync =0则不需要
// select为2代表有消息 为0代表暂无消息

func GetSyncCheck() (retCode, selector int, respStr string, err error) {
	xm := url.Values{}
	xm.Set("_", strconv.FormatInt(time.Now().Unix(), 10))
	xm.Set("DeviceID()", DeviceID())
	xm.Set("r", fmt.Sprintf("%d", time.Now().Unix()))
	xm.Set("sid", conf.Wxsid)
	xm.Set("skey", conf.Skey)
	xm.Set("uin", fmt.Sprint(conf.Wxuin))
	var synckeyStr string
	for i := 0; i < len(WebInitConf.SyncKey.List); i++ {
		synckeyStr += fmt.Sprintf("%d_%d|", WebInitConf.SyncKey.List[i].Key, WebInitConf.SyncKey.List[i].Val)
	}
	synckeyStr = synckeyStr[:len(synckeyStr)-1]
	xm.Set("synckey", synckeyStr)
	sync_url := syncUrl + "?" + xm.Encode()

	req, _ := http.NewRequest("get", sync_url, nil)
	// req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bs, _ := ioutil.ReadAll(resp.Body)
	respStr = string(bs)
	rgx := regexp.MustCompile(`window.synccheck={retcode:"(\d+)",selector:"(\d+)"}`)
	pmSub := rgx.FindStringSubmatch(respStr)
	if len(pmSub) < 2 {
		err = errors.New("window.synccheck.len(pmSub)<2")
		return
	}
	retCode, _ = strconv.Atoi(pmSub[1])
	selector, _ = strconv.Atoi(pmSub[2])
	return
}

func PostWebWxSync() (syncResp *WebWxSyncResp, lerr lerrors.Lerror) {
	xm := url.Values{}
	xm.Set("sid", conf.Wxsid)
	xm.Set("skey", conf.Skey)
	xm.Add("pass_ticket", conf.PassTicket)
	webWxSyncUrl := CgiUrl + "/webwxsync?" + xm.Encode()

	baseReq := BaseRequest{
		conf.Wxuin,
		conf.Wxsid,
		conf.Skey,
		DeviceID()}
	m := make(map[string]interface{})
	m["BaseRequest"] = baseReq
	m["rr"] = fmt.Sprintf("%d", ^time.Now().Unix())
	m["SyncKey"] = WebInitConf.SyncKey
	bs, _ := json.Marshal(m)
	req, _ := http.NewRequest("post", webWxSyncUrl, bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := Cli.Do(req)
	if err != nil {
		lerr = lerrors.Transform(err, lerrors.FATAL)
		return
	}
	defer resp.Body.Close()

	e_syncResp := WebWxSyncResp{}
	err = json.NewDecoder(resp.Body.(io.Reader)).Decode(&e_syncResp)
	if err != nil {
		lerr = lerrors.Transform(err, lerrors.FATAL)
		return
	}
	if e_syncResp.BaseResponse.Ret != 0 {
		lerr = lerrors.New(e_syncResp.BaseResponse.ErrMsg)
		return
	}
	sync_cookies = resp.Cookies()
	WebInitConf.SyncKey = e_syncResp.SyncKey

	syncResp = &e_syncResp

	return
}

type WebWxSyncResp struct {
	BaseResponse BaseResp `json:"BaseResponse"`

	SyncKey      SyncKey `json:"SyncKey"`
	SyncCheckKey SyncKey `json:"SyncCheckKey"`
	SKey         string

	AddMsgCount int
	AddMsgList  []MsgRecv

	ModContactCount        int
	ModContactList         []interface{} // todo
	DelContactCount        int
	DelContactList         []interface{}
	ModChatRoomMemberCount int
	ModChatRoomMemberList  []interface{}
	Profile                interface{}
	ContinueFlag           int
}
