package combo

import (
	"fmt"
	"test/wechat"
	"time"
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
			return
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
	err = wechat.SendMsg(to, content)
	if err != nil {
		fmt.Printf("我%s: [消息发送失败.%s]\n", where, err)
	} else {
		fmt.Printf("我%s: %s\n", where, content)
	}
	return
}
