package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func WebWxInit() (err error) {
	xm := url.Values{}
	xm.Add("pass_ticket", conf.PassTicket)
	xm.Add("skey", conf.Skey)
	xm.Add("r", fmt.Sprintf("-%d", time.Now().Unix()))
	init_url := CgiUrl + "/webwxinit?" + xm.Encode()

	breq := BaseRequest{
		conf.Wxuin,
		conf.Wxsid,
		conf.Skey,
		deviceID,
	}
	m := make(map[string]interface{})
	m["BaseRequest"] = breq
	breqdata, _ := json.Marshal(m)

	req, _ := http.NewRequest("post", init_url, bytes.NewBuffer(breqdata))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	u, _ := url.Parse(init_url)
	Cli.Jar.SetCookies(u, cookies)
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	c := InitResp{}
	err = json.NewDecoder(resp.Body.(io.Reader)).Decode(&c)
	if err != nil {
		return
	}
	webInitConf = &c
	// bs, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("WebInit", string(bs))
	return
}

type BaseRequest struct {
	Uin      string
	Sid      string
	Skey     string
	DeviceID string
}

type InitResp struct {
	BaseResponse *struct {
		Ret    int
		ErrMsg string
	} `json:"BaseResponse"`
	User                User    `json:"User"`
	Count               int     `json:"Count"`
	ContactList         []User  `json:"ContactList"`
	SyncKey             SyncKey `json:"SyncKey"`
	ChatSet             string  `json:"ChatSet"`
	SKey                string  `json:"SKey"`
	ClientVersion       int     `json:"ClientVersion"`
	SystemTime          int     `json:"SystemTime"`
	GrayScale           int     `json:"GrayScale"`
	InviteStartCount    int     `json:"InviteStartCount"`
	MPSubscribeMsgCount int     `json:"MPSubscribeMsgCount"`
	//MPSubscribeMsgList  string  `json:"MPSubscribeMsgList"`
	ClickReportInterval int `json:"ClickReportInterval"`
}

type User struct {
	UserName          string `json:"UserName"`
	Uin               int64  `json:"Uin"`
	NickName          string `json:"NickName"`
	HeadImgUrl        string `json:"HeadImgUrl" xml:""`
	RemarkName        string `json:"RemarkName" xml:""`
	PYInitial         string `json:"PYInitial" xml:""`
	PYQuanPin         string `json:"PYQuanPin" xml:""`
	RemarkPYInitial   string `json:"RemarkPYInitial" xml:""`
	RemarkPYQuanPin   string `json:"RemarkPYQuanPin" xml:""`
	HideInputBarFlag  int    `json:"HideInputBarFlag" xml:""`
	StarFriend        int    `json:"StarFriend" xml:""`
	Sex               int    `json:"Sex" xml:""`
	Signature         string `json:"Signature" xml:""`
	AppAccountFlag    int    `json:"AppAccountFlag" xml:""`
	VerifyFlag        int    `json:"VerifyFlag" xml:""`
	ContactFlag       int    `json:"ContactFlag" xml:""`
	WebWxPluginSwitch int    `json:"WebWxPluginSwitch" xml:""`
	HeadImgFlag       int    `json:"HeadImgFlag" xml:""`
	SnsFlag           int    `json:"SnsFlag" xml:""`
}

type SyncKey struct {
	Count int `json:"Count"`
	List  []struct {
		Key int `json:"Key"`
		Val int `json:"Val"`
	} `json:"List"`
}

type BaseResp struct {
	Ret    int
	ErrMsg string
}

type StatusNotifyResp struct {
	BaseResponse BaseResp `json:"BaseResponse"`
	MsgID        string   `json:"MsgID"`
}

func WebWxStatusNotify() (rsp *StatusNotifyResp, err error) {
	xm := url.Values{}
	xm.Add("pass_ticket", conf.PassTicket)
	statusnotyfy_url := CgiUrl + "/webwxstatusnotify?" + xm.Encode()
	breq := BaseRequest{
		conf.Wxuin,
		conf.Wxsid,
		conf.Skey,
		deviceID,
	}
	m := make(map[string]interface{})
	m["BaseRequest"] = breq
	m["ClientMsgId"] = time.Now().UnixNano() / 1e6
	m["Code"] = 3
	m["FromUserName"] = webInitConf.User.UserName
	m["ToUserName"] = webInitConf.User.UserName
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

type Member struct {
	Uin              int64
	UserName         string
	NickName         string
	HeadImgUrl       string
	ContactFlag      int
	MemberCount      int
	MemberList       []Member
	RemarkName       string
	HideInputBarFlag int
	Sex              int
	Signature        string
	VerifyFlag       int
	OwnerUin         int
	PYInitial        string
	PYQuanPin        string
	RemarkPYInitial  string
	RemarkPYQuanPin  string
	StarFriend       int
	AppAccountFlag   int
	Statues          int
	AttrStatus       int
	Province         string
	City             string
	Alias            string
	SnsFlag          int
	UniFriend        int
	DisplayName      string
	ChatRoomId       int
	KeyWord          string
	EncryChatRoomId  string
	IsOwner          int
}

type ContactResp struct {
	BaseResponse *BaseResp
	MemberCount  int
	MemberList   []*Member
	Seq          int
}

func GetContract() (err error) {
	xm := url.Values{}
	xm.Add("pass_ticket", conf.PassTicket)
	xm.Add("seq", "0")
	xm.Add("skey", conf.Skey)
	xm.Add("r", fmt.Sprintf("%d", time.Now().UnixNano()/1e6))
	getcontact_url := CgiUrl + "/webwxgetcontact?" + xm.Encode()
	req, _ := http.NewRequest("get", getcontact_url, nil)
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	contat := &ContactResp{}
	err = json.NewDecoder(resp.Body.(io.Reader)).Decode(contat)
	if err != nil {
		return
	}
	if contat.MemberCount > 0 {
		for i := 0; i < len(contat.MemberList); i++ {
			m_member[contat.MemberList[i].UserName] = *contat.MemberList[i]
		}
	}
	return
}

func BatchGetContact() (err error) {
	xm := url.Values{}
	xm.Add("pass_ticket", conf.PassTicket)
	xm.Add("type", "ex")
	xm.Add("r", fmt.Sprintf("%d", time.Now().UnixNano()/1e6))
	xm.Add("lang", "zh_CN")
	morecontact_url := CgiUrl + "/webwxbatchgetcontact?" + xm.Encode()
	req, _ := http.NewRequest("post", morecontact_url, nil)
}
