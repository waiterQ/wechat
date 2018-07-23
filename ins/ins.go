package ins

import (
	"fmt"
	"strings"

	"test/wechat"
	lerrors "test/wechat/errors"
	"test/wechat/utils"
)

// 关于键入指令

func Start(exitCh chan<- byte, errCh chan<- lerrors.Lerror, msgCh chan<- []string) (err error) {
	for {
		originCmd := utils.ReadStdin()
		originCmd = strings.TrimLeft(originCmd, " ")
		if originCmd == "" {
			fmt.Println("缺少键入")
			continue
		}
		var inser Instruction
		for _, cmd := range Cmds {
			inser = cmd.New(originCmd)
			if inser != nil {
				break
			}
		}
		if inser == nil { // sendmsg default
			inser = &SendCmd{}
			if wechat.Curr_chatObj.UserName == "" {
				fmt.Println("请选择聊天对象")
				continue
			}
			originCmd = "send " + originCmd
		}
		values := inser.Prepare(originCmd, exitCh, errCh, msgCh)
		stop := inser.Exec(values)
		if stop {
			break
		}
	}
	return
}

func init() {
	Cmds = append(Cmds,
		new(ExitCmd),     // 退出
		new(SendCmd),     // 发送
		new(ChangetoCmd), // 改变当前聊天对象
		new(RoomlistCmd), // 聊天室列表
		new(QueryCmd),    // 查询
		new(WithdrawCmd), // 撤回
		new(CurrObj))     // 查看当前聊天对象
}
