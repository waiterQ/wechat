package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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
		DeviceID(),
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
	WebInitConf = &c
	for i := 0; i < len(c.ContactList); i++ {
		m_member[c.ContactList[i].UserName] = c.ContactList[i]
	}
	me_userName = c.User.UserName
	return
}

type BaseRequest struct {
	Uin      string
	Sid      string
	Skey     string
	DeviceID string
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

type BaseResp struct {
	Ret    int
	ErrMsg string
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
	Count        int // batchgetcontract type=ex
	MemberList   []*Member
	Seq          int
}

func GetContract() (code string, err error) {
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
	// fmt.Println(contat)
	if contat.BaseResponse.Ret != 0 {
		return fmt.Sprintf("%d", contat.BaseResponse.Ret), errors.New(contat.BaseResponse.ErrMsg)
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
	r := BatchContactReq{}
	r.BaseReq = &BaseRequest{
		conf.Wxuin,
		conf.Wxsid,
		conf.Skey,
		DeviceID(),
	}
	for i := 0; i < len(WebInitConf.ContactList); i++ {
		if strings.Contains(WebInitConf.ContactList[i].UserName, "@") {
			one := ContactOne{
				UserName: WebInitConf.ContactList[i].UserName,
			}
			r.List = append(r.List, one)
		}
	}
	r.Count = len(r.List)
	bs, _ := json.Marshal(r)

	req, _ := http.NewRequest("post", morecontact_url, bytes.NewReader(bs))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
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
	if contat.BaseResponse.Ret != 0 && contat.BaseResponse.Ret != -1 {
		return errors.New(contat.BaseResponse.ErrMsg)
	}
	if contat.Count > 0 {
		for i := 0; i < len(contat.MemberList); i++ {
			m_member[contat.MemberList[i].UserName] = *contat.MemberList[i]
		}
	}
	return
}

type BatchContactReq struct {
	BaseReq *BaseRequest `json:"BaseRequest"`
	Count   int
	List    []ContactOne
}
type ContactOne struct {
	UserName string
}
