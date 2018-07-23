package combo

import (
	"errors"
	"fmt"
	"time"

	"test/wechat"
)

// 第一次
func Say(to, content string) (err error) {
	if !wechat.HaveFirstSendMsg {
		wechat.HaveFirstSendMsg = true
		wechat.FirstSendMsgTime = time.Now()
	}
	return SayBase(to, content)
}

func SayBase(to, content string) (err error) {
	_, ok := wechat.M_chatroomViewed[to]
	if !ok {
		wechat.M_chatroomViewed[to] = 1
		mber := wechat.M_member[to]
		var mbers []string
		for i := 0; i < len(mber.MemberList); i++ {
			mbers = append(mbers, mber.MemberList[i].UserName)
		}
		err = wechat.BatchGetContact(mbers, false)
		if err != nil {
			fmt.Printf("SayBase.BatchGetContact(room:%s)\n", to)
		}
	}
	return SendMsg(to, content)
}

func SendMsg(to, content string) (err error) {
	where := ""
	if to[:2] == "@@" {
		where = "在 "
	} else {
		where = "对 "
	}
	where += wechat.NickName(to)
	var cliMsgid, svrMsgid string
	cliMsgid, svrMsgid, err = wechat.SendMsg(to, content)
	if err != nil {
		fmt.Printf("我%s: [消息发送失败.%s]\n", where, err)
		err = nil // 发送失败的错误忽略
	} else {
		wechat.MsCount[to] += 1
		fmt.Printf("我%s:(%s)\n 	%s\n", where, svrMsgid, content)

		wechat.Me_Said[svrMsgid] = wechat.MsgRecd{ // 我的 说话记录
			SvrMsgid:   svrMsgid,
			SendTime:   time.Now(),
			Tousername: to,
			CliMsgid:   cliMsgid}
		wechat.LastSendMsg = wechat.Me_Said[svrMsgid]
	}
	return
}

// 发送消息
func MsgProcess(pars []string) (std_values []string, err error) {
	sendMsg := pars[0]
	var chat_obj string
	if len(pars) == 1 {
		if wechat.Curr_chatObj.UserName == "" {
			err = errors.New("请选择一个聊天对象")
			return
		}
		chat_obj = wechat.Curr_chatObj.UserName
	} else {
		chat_obj = pars[1]
	}
	code, list := ContactQry(chat_obj)
	switch code {
	case -1:
		err = errors.New("未找到聊天对象")
		return
	case 0:
		chat_obj = list[0].UserName
		wechat.Curr_chatObj.UserName = list[0].UserName
		wechat.Curr_chatObj.NickName = list[0].NickName
	case 1:
		err = fmt.Errorf("对象不是你朋友,%s:[%s]", list[0].UserName, list[0].NickName)
		return
	case 2:
		errmsg := "找到多个对象\n"
		for i := 0; i < len(list); i++ {
			errmsg += fmt.Sprintf("%s:[%s]\n", list[i].UserName, list[i].NickName)
		}
		err = errors.New(errmsg)
		return
	default:
		err = errors.New("未知错误")
		return
	}
	std_values = []string{sendMsg, chat_obj}
	return
}
