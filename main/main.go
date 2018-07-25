package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mdp/qrterminal"

	"test/wechat"
	"test/wechat/combo"
	lerrors "test/wechat/errors"
	"test/wechat/ins"
)

func main() {
	uuid, err := wechat.GetLoginUuid()
	if err != nil {
		fmt.Println("[wechat.GetLoginUuid]", err)
		return
	}

	qrterminal.Generate("https://login.weixin.qq.com/l/"+uuid, qrterminal.L, os.Stdout)
	// err = wechat.DownloadImgQR(uuid)
	// if err != nil {
	// 	fmt.Println("[wechat.GetImgQR(uuid)]", err)
	// 	return
	// }
	fmt.Println("扫描二维码授权登录")
	var redirect_url, code string
	for {
		time.Sleep(time.Second)
		redirect_url, code, err = wechat.GetLoginUrlAfterAuth(uuid)
		if err != nil {
			fmt.Println("[wechat.GetLoginUrlAfterAuth(uuid)]", err)
			return
		}
		if code == "200" {
			break
		}
	}

	// report0
	err = wechat.Report0()
	if err != nil {
		fmt.Println("[wechat.Report0()]", err)
		return
	}

	fmt.Println("redirect_url=", redirect_url)
	err = wechat.Login(redirect_url)
	if err != nil {
		fmt.Println("[wechat.Login(redirect_url)]", err)
		return
	}

	err = wechat.WebWxInit()
	if err != nil {
		fmt.Println("[wechat.WebWxInit()]", err)
		return
	}
	statusNotify_resp, err := wechat.WebWxStatusNotify(wechat.WebInitConf.User.UserName, wechat.WebInitConf.User.UserName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if statusNotify_resp.BaseResponse.Ret != 0 {
		fmt.Println("statusNotify_resp.BaseResponse", statusNotify_resp.BaseResponse)
		return
	}
	fmt.Println("statusNotify_resp.MsgID", statusNotify_resp.MsgID)
	code = ""
	code, err = wechat.GetContract()
	if err != nil {
		fmt.Println(code, err)
		return
	}
	fmt.Println("wechat.GetContract().Done")
	// report1
	wechat.Report1()
	if err != nil {
		fmt.Println("[wechat.Report1()]", err)
		return
	}

	syncResp, err := wechat.PostWebWxSync()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("wechat.PostWebWxSync().Done")
	rooms, err := combo.GetLastChatrooms(syncResp)
	if err != nil {
		fmt.Println("combo.GetLastChatrooms(syncResp)", err)
		return
	}
	fmt.Println("combo.GetLastChatrooms(syncResp).Done")
	if len(rooms) > 0 {
		err = wechat.BatchGetContact(rooms, false) // 群里的可能不是直接好友
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("wechat.BatchGetContact().Done")
	}
	wechat.SetSyncCookies()
	fmt.Println("初始化完成")

	// 注意顺序 [收到消息时的处理逻辑]
	wechat.RecHandls = append(wechat.RecHandls, wechat.DisplayMsg) // 显示消息

	errCh := make(chan lerrors.Lerror) // 接收致命错误
	exitCh := make(chan byte)          // 控制程序退出

	go SyncRecv(errCh)
	fmt.Println("接收消息启动完成")
	// after first send msg
	go FirstSendReport(errCh)
	fmt.Println("微信同步报告启动完成")

	msgCh := make(chan []string, 20)
	go SvrSend(msgCh, errCh)
	fmt.Println("发送消息启动完成")

	go ins.Start(exitCh, errCh, msgCh)
	fmt.Printf("键入指令启动完成\n\n")

	go errorListener(errCh, exitCh) // 大错误处理中心

	fmt.Println("启动完成：")

	<-exitCh
}

func SyncRecv(errCh chan<- lerrors.Lerror) {
	for {
		retCode, selector, respStr, err := wechat.GetSyncCheck()
		if err != nil {
			errCh <- lerrors.Transform(err, lerrors.FATAL)
			return
		}
		// fmt.Println(retCode, selector, respStr, err)
		switch retCode {
		case 1100:
			errCh <- lerrors.New("在微信上退出", lerrors.FATAL)
			return
		case 1101:
			if !wechat.ManualQuit {
				errCh <- lerrors.New("在其他设备上登录", lerrors.FATAL)
			} else {
				fmt.Println("手动退出")
			}
			return
		case 0:
			switch selector {
			case 2, 3, 4, 6: // 3的数据 profile
				if selector == 4 {
					fmt.Println("[通讯录更新了]")
				} else if selector == 6 {
					fmt.Println("//========= 红包来了! ========//")
				}
				syncResp, lerr := wechat.PostWebWxSync()
				if selector == 4 || selector == 6 {
					d, _ := json.Marshal(syncResp)
					fmt.Println(string(d))
				}
				if lerr != nil {
					errCh <- lerr
					if lerr.Level() > lerrors.ERROR {
						return
					}
				}
				wechat.SetSyncCookies()
				wechat.HandleRecvMsg(syncResp)
			// case 4: // 通讯录更新
			// 	fmt.Println("[通讯录更新了]")
			// case 6:
			// 	fmt.Println("//========= 红包来了! ========//")
			case 7:
				fmt.Println("在手机上操作了微信")
			case 0:
			default:
				fmt.Println(respStr)
				errCh <- lerrors.New("未知selector:" + fmt.Sprint(selector))
			}
		default:
			fmt.Println(respStr)
			errCh <- lerrors.New("retCode:"+fmt.Sprint(retCode), lerrors.FATAL)
			return
		}
	}
}

func FirstSendReport(errCh chan<- lerrors.Lerror) {
	for {
		time.Sleep(10 * time.Minute)
		if wechat.HaveFirstSendMsg {
			err := wechat.ReportSendMsg()
			if err != nil {
				errCh <- lerrors.Transform(err, lerrors.FATAL)
			}
			return
		}
	}
}

func SvrSend(msgCh <-chan []string, errCh chan<- lerrors.Lerror) {
	for ss := range msgCh {
		err := combo.Say(ss[1], ss[0])
		if err != nil {
			errCh <- lerrors.Transform(err) // 虽然接不到错误 但是先放着
		}
	}
}

func errorListener(errCh <-chan lerrors.Lerror, exitCh chan<- byte) {
	var count int // lerrors.ERROR次数 到2次或致命错误 就退出
	for {
		if count > 1 {
			exitCh <- 0
			return
		}
		for lerr := range errCh {
			if lerr.Level() == lerrors.FATAL {
				fmt.Printf("[FATAL] %s\n", lerr)
				exitCh <- 0
				return
			} else if lerr.Level() == lerrors.ERROR {
				fmt.Printf("[ERROR] %s\n", lerr)
				count += 1
			}
		}
	}
}
