package combo

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"test/wechat"
	"test/wechat/utils"
)

// 键入命令处理

func ExecIns(exitCh chan<- error, msgCh chan<- []string) (err error) {
	for {
		cmd, values := InputPrepare()
		if cmd == "" {
			err = errors.New("缺少键入")
			continue
		}
		switch cmd {
		case "exit":
			err = Logout()
			if err != nil {
				fmt.Println(err)
			}
			exitCh <- errors.New("手动退出")
			return
		case "send":
			err = MsgProcess(values, msgCh)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "changeto":
			err = ChangeTo(values[0])
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "roomlist":
			code, list := ChatroomQry(values)
			if code == -1 {
				fmt.Println("没有找到聊天室")
			} else {
				msg := "可选聊天室:\n"
				for i := 0; i < len(list); i++ {
					msg += fmt.Sprintf("%s:[%s]\n", list[i].UserName, list[i].NickName)
				}
				fmt.Println(msg)
			}
		case "query": // 暂时查询username nickname
			_, list := ContactQry(values[0])
			msg := "结果:\n"
			for i := 0; i < len(list); i++ {
				msg += fmt.Sprintf("%s:[%s]\n", list[i].UserName, list[i].NickName)
			}
			fmt.Println(msg)
		case "withdraw":
			err := Withdraw(values)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("消息撤回成功")
		default:
			err = errors.New("未知指令")
			fmt.Println(err)
		}
	}
	return
}

func InputPrepare() (cmd string, values []string) {
	origin := utils.ReadStdin()
	origin = strings.TrimLeft(origin, " ")
	if origin == "" {
		return
	}
	if origin == "exit" || strings.HasPrefix(origin, "/exit") {
		cmd = "exit"
		return
	}
	if strings.HasPrefix(origin, "send ") || strings.HasPrefix(origin, "say ") || strings.HasPrefix(origin, "/send ") || strings.HasPrefix(origin, "/say ") {
		ss := strings.SplitN(origin, " ", 2)
		cmd = "send"
		if strings.Contains(ss[1], " to:") {
			ss2 := strings.SplitN(ss[1], " to:", 2)
			values = ss2
		} else {
			values = append(values, ss[1])
		}
		return
	}
	if strings.HasPrefix(origin, "changeto ") || strings.HasPrefix(origin, "/changeto ") {
		cmd = "changeto"
		ss := strings.SplitN(origin, " ", 2)
		values = append(values, ss[1])
		return
	}
	if origin == "/roomlist" || strings.HasPrefix(origin, "/roomlist ") {
		cmd = "roomlist"
		if origin == "/roomlist" {
			return
		}
		ss := strings.SplitN(origin, " ", 2)
		values = append(values, ss[1])
		return
	}
	if strings.HasPrefix(origin, "/query ") { // todo a lot can add 多级查询
		cmd = "query"
		ss := strings.SplitN(origin, " ", 2)
		values = append(values, ss[1])
		return
	}
	if origin == "withdraw" || origin == "/withdraw" || strings.HasPrefix(origin, "withdraw ") || strings.HasPrefix(origin, "/withdraw ") {
		cmd = "withdraw"
		if origin == "withdraw" || origin == "/withdraw" {
			return
		}
		ss := strings.SplitN(origin, " ", 2)
		values = append(values, ss[1])
		return
	}
	cmd = "send"
	if strings.Contains(origin, " to:") {
		ss := strings.SplitN(origin, " to:", 2)
		values = ss
	} else {
		values = append(values, origin)
	}
	return
}

// 发送消息
func MsgProcess(pars []string, msgCh chan<- []string) (err error) {
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
	msgCh <- []string{sendMsg, chat_obj}
	return
}

// 切换当前聊天对象
func ChangeTo(obj string) (err error) {
	if obj == "" {
		err = errors.New("请选择一个聊天对象")
		return
	}
	code, list := ContactQry(obj)
	switch code {
	case -1:
		err = errors.New("未找到聊天对象")
		return
	case 0:
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

	return
}

func Withdraw(pars []string) (err error) {
	var cliMsgid, svrMsgid, toUserName string
	if len(pars) == 0 {
		if wechat.LastSendMsg.SendTime.Add(time.Minute * 2).Before(time.Now()) {
			err = errors.New("[消息超过两分钟无法撤回]")
			return
		}
		cliMsgid, svrMsgid, toUserName = wechat.LastSendMsg.CliMsgid, wechat.LastSendMsg.SvrMsgid, wechat.LastSendMsg.Tousername
	} else {
		rcd, ok := wechat.Me_Said[pars[0]]
		if !ok {
			err = errors.New("未找到msgid")
			return
		}
		if rcd.SendTime.Add(time.Minute * 2).Before(time.Now()) {
			err = errors.New("[消息超过两分钟无法撤回]")
			return
		}
		cliMsgid, svrMsgid, toUserName = rcd.CliMsgid, rcd.SvrMsgid, rcd.Tousername
	}
	msg, err := wechat.Withdraw(cliMsgid, svrMsgid, toUserName)
	if err != nil {
		return
	}
	if msg != "" {
		err = errors.New(msg)
	}
	return
}
