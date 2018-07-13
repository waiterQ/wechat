package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 每发送一条消息，或接收到消息 get.syncCheck会立即返回 同时post.syncCheck获取消息

func SendMsg(to, content string) (err error) {
	xm := url.Values{}
	xm.Add("pass_ticket", conf.PassTicket)
	sendMsg_url := CgiUrl + "/webwxsendmsg?" + xm.Encode()
	r := MsgReq{}
	r.BaseReq = &BaseRequest{
		conf.Wxuin,
		conf.Wxsid,
		conf.Skey,
		DeviceID(),
	}
	localMsgID := LocalMsgID()
	r.Msg.ClientMsgId = localMsgID
	r.Msg.LocalID = localMsgID
	r.Msg.Content = content
	r.Msg.ToUserName = to
	r.Msg.FromUserName = WebInitConf.User.UserName
	r.Msg.Type = 1 // 1为文本格式
	bs, _ := json.Marshal(r)
	req, _ := http.NewRequest("post", sendMsg_url, bytes.NewReader(bs))
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var rsp MsgResp
	err = json.NewDecoder(resp.Body.(io.Reader)).Decode(&rsp)
	if err != nil {
		return
	}
	if rsp.BaseResponse.Ret != 0 {
		return errors.New(rsp.BaseResponse.ErrMsg)
	}
	return
}

type MsgReq struct {
	BaseReq *BaseRequest `json:"BaseRequest"`
	Msg     struct {
		ClientMsgId  string
		Content      string
		FromUserName string
		LocalID      string
		ToUserName   string
		Type         int
	}
	Scene int
}

type MsgResp struct {
	BaseResponse *BaseResp
	MsgID        string
	LocalID      string
}

func LocalMsgID() string {
	return strconv.Itoa(int(time.Now().Unix())) + "0" + strconv.Itoa(rand.Int())[3:6]
}

type MsgRecv struct {
	MsgId        string
	FromUserName string
	ToUserName   string
	MsgType      int
	Content      string
	CreateTime   int64

	Status               int
	ImgStatus            int
	VoiceLength          int
	PlayLength           int
	FileName             string
	FizeSize             string
	MediaId              string
	Url                  string
	AppMsgType           int
	StatusNotifyCode     int         // 4所有联系人和chatroom 有排序 5单个
	StatusNotifyUserName string      // 逗号分隔的userName
	RecommendInfo        interface{} // ?
	AppInfo              struct {    // ?
		AppID string
		Type  int
	}
	HasProductId  int
	Ticket        string
	ImgHeight     float64
	Imgwidth      float64
	SubMsgType    int
	NewMsgId      int
	OriContent    string
	EncryFileName string
}

func HandleRecvMsg(syncResp *WebWxSyncResp) {
	for i := 0; i < len(syncResp.AddMsgList); i++ {
		var from, content string
		if strings.Contains(syncResp.AddMsgList[i].FromUserName, "@@") {
			ss := strings.SplitN(syncResp.AddMsgList[i].Content, ":<br/>", 2)
			from = fmt.Sprintf("chatroom[%s]的 %s", NickName(syncResp.AddMsgList[i].FromUserName), NickName(ss[0]))
			content = ss[1]
		} else {
			from = NickName(syncResp.AddMsgList[i].FromUserName)
			content = syncResp.AddMsgList[i].Content
		}
		// to = NickName(syncResp.AddMsgList[i].ToUserName)
		switch syncResp.AddMsgList[i].MsgType {
		case MSG_TEXT:
			if syncResp.AddMsgList[i].SubMsgType == MSG_LOCATION {
				ss := strings.SplitN(content, ":<br/>", 2)
				content = fmt.Sprintf("[%s %s]", ss[0], syncResp.AddMsgList[i].Url)
			} else {
				content = strings.Replace(content, "<br/>", "\n", -1)
				fmt.Printf("%s:\n	%s\n", from, content)
			}
		case MSG_WITHDRAW:
			content = "[撤回一条消息]"
			fmt.Printf("%s:\n	%s\n", from, content)
		case MSG_LINK:
			content = strings.Replace(content, "<br/>", "\n", -1)
			content = html.UnescapeString(content)
			fmt.Printf("%s:\n	%s\n", from, content)
		case MSG_IMG:
			content = fmt.Sprintf("%s/webwxgetmsgimg?&MsgID=%s&skey=%s&type=slave", CgiUrl, syncResp.AddMsgList[i].MsgId, WebInitConf.SKey)
			fmt.Printf("%s:\n	[收到一个图片 %s]\n", from, content)
			// 没有cookie无法访问
		case MSG_EMOTION:
			content = html.UnescapeString(content)
			content = strings.Replace(content, " ", "", -1)
			ss := strings.SplitN(content, `cdnurl="`, 2)
			content = ss[1]
			ss = strings.SplitN(content, `"`, 2)
			content = ss[0]
			fmt.Printf("%s:\n	[收到一个gif %s]\n", from, content)
		case MSG_VOICE:
			content = "[收到一段语音]"
			fmt.Printf("%s:\n	%s\n", from, content)
		case MSG_FV:
			content = "[收到好友验证消息]"
			fmt.Printf("%s:\n	%s\n", from, content)
		case MSG_SCC:
			content = "[收到一张好友明片]"
			fmt.Printf("%s:\n	%s\n", from, content)
		case MSG_VIDEO:
			content = "[收到一段视频]"
			fmt.Printf("%s:\n	%s\n", from, content)
		case MSG_VOIP, MSG_INIT, MSG_VOIPNOTIFY, MSG_VOIPINVITE:
		case MSG_SHORT_VIDEO:
			content = "[收到小视频]"
			fmt.Printf("%s:\n	%s\n", from, content)
		case MSG_SYSNOTICE:
			content = "[收到系统通知]"
			fmt.Printf("%s:\n	%s\n", from, content)
		case MSG_SYS:
			content = "[收到系统消息]"
			fmt.Printf("%s:\n	%s\n", from, content)
		default:
			content = fmt.Sprintf("[收到MsgType:%d]", syncResp.AddMsgList[i].MsgType)
			fmt.Printf("%s:\n	%s\n", from, content)
		}
	}
}

const (
	// msg types
	MSG_TEXT        = 1     // text message
	MSG_IMG         = 3     // image message
	MSG_VOICE       = 34    // voice message
	MSG_FV          = 37    // friend verification message
	MSG_PF          = 40    // POSSIBLEFRIEND_MSG
	MSG_SCC         = 42    // shared contact card
	MSG_VIDEO       = 43    // video message
	MSG_EMOTION     = 47    // gif
	MSG_LOCATION    = 48    // location message
	MSG_LINK        = 49    // shared link message
	MSG_VOIP        = 50    // VOIPMSG
	MSG_INIT        = 51    // wechat init message
	MSG_VOIPNOTIFY  = 52    // VOIPNOTIFY
	MSG_VOIPINVITE  = 53    // VOIPINVITE
	MSG_SHORT_VIDEO = 62    // short video message
	MSG_SYSNOTICE   = 9999  // SYSNOTICE
	MSG_SYS         = 10000 // system message
	MSG_WITHDRAW    = 10002 // withdraw notification message
)
