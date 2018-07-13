package wechat

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var Cli *http.Client
var conf *XmlConfig
var WebInitConf *InitResp
var CgiUrl string
var M_member map[string]Contact = make(map[string]Contact) // GetContact batchgetContact webwxinit.chatrooms
var M_chatroomViewed map[string]int = make(map[string]int) // chatroom第一次进入需要获取所有群成员
// var curr_obj string                                        // 当前聊天对象
var Me_userName string // 自己的userName
var randSourse rand.Source

type XmlConfig struct {
	XMLName     xml.Name `xml:"error"`
	Ret         int      `xml:"ret"`
	Message     string   `xml:"message"`
	Skey        string   `xml:"skey"`
	Wxsid       string   `xml:"wxsid"`
	Wxuin       int      `xml:"wxuin"`
	PassTicket  string   `xml:"pass_ticket"`
	IsGrayscale int      `xml:"isgrayscale"`
}

type InitResp struct {
	BaseResponse *struct {
		Ret    int
		ErrMsg string
	} `json:"BaseResponse"`
	User                User     `json:"User"`
	Count               int      `json:"Count"`
	ContactList         []Member `json:"ContactList"`
	SyncKey             SyncKey  `json:"SyncKey"`
	ChatSet             string   `json:"ChatSet"`
	SKey                string   `json:"SKey"`
	ClientVersion       int      `json:"ClientVersion"`
	SystemTime          int      `json:"SystemTime"`
	GrayScale           int      `json:"GrayScale"`
	InviteStartCount    int      `json:"InviteStartCount"`
	MPSubscribeMsgCount int      `json:"MPSubscribeMsgCount"`
	//MPSubscribeMsgList  string  `json:"MPSubscribeMsgList"`
	ClickReportInterval int `json:"ClickReportInterval"`
}

type SyncKey struct {
	Count int `json:"Count"`
	List  []struct {
		Key int `json:"Key"`
		Val int `json:"Val"`
	} `json:"List"`
}

func NickName(userName string) (name string) {
	if userName == Me_userName {
		name += "我"
		return
	}
	member, ok := M_member[userName]
	if ok {
		name += member.NickName
		return
	}
	name += userName
	return
}

func DeviceID() string {
	for {
		num := randSourse.Int63()
		if num > 100000000000000 {
			s := fmt.Sprint(num)
			return "e" + s[len(s)-15:]
		}
	}
}

func init() {
	randSourse = rand.NewSource(time.Now().Unix())
}
