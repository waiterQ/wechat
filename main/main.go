package main

import (
	"errors"
	"fmt"
	"test/wechat"
	"test/wechat/combo"
	"time"
)

func main() {
	uuid, err := wechat.GetLoginUuid()
	if err != nil {
		fmt.Println("[wechat.GetLoginUuid]", err)
		return
	}
	err = wechat.GetImgQR(uuid)
	if err != nil {
		fmt.Println("[wechat.GetImgQR(uuid)]", err)
		return
	}
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

	chErr := make(chan error)

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
	if len(rooms) > 0 {
		err = wechat.BatchGetContact(rooms, false) // 群里的可能不是直接好友
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("wechat.BatchGetContact().Done")
	}
	wechat.SetSyncCookies()

	go SyncRecv(chErr)
	// after first send msg
	go FirstSendReport(chErr)

	err = <-chErr
	if err != nil {
		fmt.Println(err)
	}
}

func SyncRecv(chErr chan<- error) {
	for {
		retCode, selector, respStr, err := wechat.GetSyncCheck()
		if err != nil {
			chErr <- err
			return
		}
		// fmt.Println(retCode, selector, respStr, err)
		switch retCode {
		case 1100:
			chErr <- errors.New("在微信上退出")
			return
		case 1101:
			chErr <- errors.New("在其他设备上登录")
			return
		case 0:
			switch selector {
			case 2, 3:
				syncResp, err := wechat.PostWebWxSync()
				if err != nil {
					chErr <- err
					return
				}
				wechat.SetSyncCookies()
				wechat.HandleRecvMsg(syncResp)
			case 4: // 通讯录更新
				fmt.Println("通讯录更新了")
			case 6:
				fmt.Println("//========= 红包来了! ========//")
			case 7:
				fmt.Println("在手机上操作了微信")
			case 0:
			default:
				fmt.Println(respStr)
				chErr <- errors.New("未知selector:" + fmt.Sprint(selector))
				return
			}
		default:
			fmt.Println(respStr)
			chErr <- errors.New("retCode:" + fmt.Sprint(retCode))
			return
		}
	}
}

func FirstSendReport(chErr chan<- error) {
	for {
		time.Sleep(10 * time.Minute)
		if wechat.HaveFirstSendMsg {
			err := wechat.ReportSendMsg()
			if err != nil {
				chErr <- err
			}
			return
		}
	}
}
