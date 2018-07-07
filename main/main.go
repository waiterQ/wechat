package main

import (
	"fmt"
	// "os"
	"test/wechat"
	"time"
)

// func init() {
// 	os.Remove("./4dpdzifRHg==.jpg")
// 	time.Sleep(time.Hour)
// }

func main() {
	uuid, err := wechat.GetLoginUuid()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = wechat.GetImgQR(uuid)
	if err != nil {
		fmt.Println(err)
		return
	}
	var redirect_url, code string
	for {
		time.Sleep(time.Second)
		redirect_url, code, err = wechat.GetLoginUrlAfterAuth(uuid)
		if err != nil {
			fmt.Println(err)
			return
		}
		if code == "200" {
			break
		}
	}

	fmt.Println("redirect_url=", redirect_url)
	err = wechat.Login(redirect_url)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = wechat.WebWxInit()
	if err != nil {
		fmt.Println(err)
		return
	}
	statusNotify_resp, err := wechat.WebWxStatusNotify()
	if err != nil {
		fmt.Println(err)
		return
	}
	if statusNotify_resp.BaseResponse.Ret != 0 {
		fmt.Println("statusNotify_resp.BaseResponse", statusNotify_resp.BaseResponse)
		return
	}
	fmt.Println("statusNotify_resp.MsgID", statusNotify_resp.MsgID)
	contact_resp, err := wechat.GetContract()
	if err != nil {
		fmt.Println(err)
		return
	}
	if contact_resp.BaseResponse.Ret != 0 {
		fmt.Println("contact_resp.BaseResponse", contact_resp.BaseResponse)
		return
	}
	fmt.Println("contact_resp", contact_resp)
	return
}
