package ins

import (
	"fmt"
	"strings"

	"test/wechat"
	"test/wechat/combo"
	lerrors "test/wechat/errors"
)

// exit
type ExitCmd struct {
	ExitCh chan<- byte
	ErrCh  chan<- lerrors.Lerror
}

func (cmd *ExitCmd) New(originCmd string) (inser Instruction) {
	if originCmd == "exit" || strings.HasPrefix(originCmd, "/exit") {
		inser = &ExitCmd{}
		return
	}
	return
}

// ctrls for channel exitCh and errCh
func (cmd *ExitCmd) Prepare(originCmd string, ctrls ...interface{}) (values []string) {
	cmd.ExitCh = ctrls[0].(chan<- byte)
	cmd.ErrCh = ctrls[1].(chan<- lerrors.Lerror)
	return
}

func (cmd *ExitCmd) Exec(values []string) (stop bool) {
	stop = true
	lerr := combo.Logout()
	if lerr != nil {
		cmd.ErrCh <- lerr
		return
	}
	cmd.ExitCh <- 0
	return
}

// send message
type SendCmd struct {
	MsgCh chan<- []string
}

func (cmd *SendCmd) New(originCmd string) (inser Instruction) {
	if strings.HasPrefix(originCmd, "send ") || strings.HasPrefix(originCmd, "say ") || strings.HasPrefix(originCmd, "/send ") || strings.HasPrefix(originCmd, "/say ") {
		inser = &SendCmd{}
		return
	}
	return
}

func (cmd *SendCmd) Prepare(originCmd string, ctrls ...interface{}) (values []string) {
	cmd.MsgCh = ctrls[2].(chan<- []string)
	ss := strings.SplitN(originCmd, " ", 2)
	vStr := ""
	if len(ss) == 1 {
		vStr = originCmd
	} else {
		vStr = ss[1]
	}
	if strings.Contains(vStr, " to:") {
		ss2 := strings.SplitN(vStr, " to:", 2)
		values = ss2
	} else {
		values = append(values, vStr)
	}
	return
}

func (cmd *SendCmd) Exec(values []string) (stop bool) {
	stdVs, err := combo.MsgProcess(values)
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd.MsgCh <- stdVs
	return
}

// changeto
type ChangetoCmd struct{}

func (cmd *ChangetoCmd) New(originCmd string) (inser Instruction) {
	if strings.HasPrefix(originCmd, "changeto ") || strings.HasPrefix(originCmd, "/changeto ") {
		inser = &ChangetoCmd{}
		return
	}
	return
}

func (cmd *ChangetoCmd) Prepare(originCmd string, ctrls ...interface{}) (values []string) {
	ss := strings.SplitN(originCmd, " ", 2)
	values = append(values, ss[1])
	return
}

func (cmd *ChangetoCmd) Exec(values []string) (stop bool) {
	err := combo.ChangeTo(values[0])
	if err != nil {
		fmt.Println(err)
	}
	return
}

// roomlist
type RoomlistCmd struct{}

func (cmd *RoomlistCmd) New(originCmd string) (inser Instruction) {
	if originCmd == "roomlist" || strings.HasPrefix(originCmd, "roomlist ") || originCmd == "/roomlist" || strings.HasPrefix(originCmd, "/roomlist ") {
		inser = &RoomlistCmd{}
		return
	}
	return
}

func (cmd *RoomlistCmd) Prepare(originCmd string, ctrls ...interface{}) (values []string) {
	if originCmd == "/roomlist" || originCmd == "roomlist" {
		return
	}
	ss := strings.SplitN(originCmd, " ", 2)
	values = append(values, ss[1])
	return
}

func (cmd *RoomlistCmd) Exec(values []string) (stop bool) {
	code, list := combo.ChatroomQry(values)
	if code == -1 {
		fmt.Println("没有找到聊天室")
	} else {
		msg := "可选聊天室:\n"
		for i := 0; i < len(list); i++ {
			msg += fmt.Sprintf("%s:[%s]\n", list[i].UserName, list[i].NickName)
		}
		fmt.Println(msg)
	}
	return
}

// query (todo more) 暂时查询username nickname
type QueryCmd struct{}

func (cmd *QueryCmd) New(originCmd string) (inser Instruction) {
	if strings.HasPrefix(originCmd, "query ") || strings.HasPrefix(originCmd, "/query ") { // todo a lot can add 多级查询
		inser = &QueryCmd{}
		return
	}
	return
}

func (cmd *QueryCmd) Prepare(originCmd string, ctrls ...interface{}) (values []string) {
	ss := strings.SplitN(originCmd, " ", 2)
	values = append(values, ss[1])
	return
}

func (cmd *QueryCmd) Exec(values []string) (stop bool) {
	_, list := combo.ContactQry(values[0])
	msg := "结果:\n"
	for i := 0; i < len(list); i++ {
		msg += fmt.Sprintf("%s:[%s]\n", list[i].UserName, list[i].NickName)
	}
	fmt.Println(msg)
	return
}

// withdraw message
type WithdrawCmd struct {
	ErrCh chan<- lerrors.Lerror
}

func (cmd *WithdrawCmd) New(originCmd string) (inser Instruction) {
	if originCmd == "wd" || originCmd == "withdraw" || originCmd == "/withdraw" || strings.HasPrefix(originCmd, "withdraw ") || strings.HasPrefix(originCmd, "/withdraw ") {
		inser = &WithdrawCmd{}
		return
	}
	return
}

func (cmd *WithdrawCmd) Prepare(originCmd string, ctrls ...interface{}) (values []string) {
	if originCmd == "wd" || originCmd == "withdraw" || originCmd == "/withdraw" {
		return
	}
	ss := strings.SplitN(originCmd, " ", 2)
	values = append(values, ss[1])
	return
}

func (cmd *WithdrawCmd) Exec(values []string) (stop bool) {
	lerr := combo.Withdraw(values)
	if lerr != nil {
		if lerr.Level() > lerrors.WARN {
			cmd.ErrCh <- lerr
		} else {
			fmt.Println(lerr)
		}
		return
	}
	fmt.Println("消息撤回成功")
	return
}

type CurrObj struct{}

func (cmd *CurrObj) New(originCmd string) (inser Instruction) {
	if originCmd == "currobj" || strings.HasPrefix(originCmd, "/currobj") {
		inser = &CurrObj{}
		return
	}
	return
}

func (cmd *CurrObj) Prepare(originCmd string, ctrls ...interface{}) (values []string) { return }

func (cmd *CurrObj) Exec(values []string) (stop bool) {
	fmt.Printf("当前对象为: %s(%s)\n", wechat.Curr_chatObj.NickName, wechat.Curr_chatObj.UserName)
	return
}
