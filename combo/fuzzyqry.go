package combo

import (
	// "fmt"
	"strings"
	"test/wechat"
)

// -1未找到 0找到唯一 1找到唯一不是朋友 2找到多个
// 已有联系人找 username nickname
func ContactQry(par string) (code int, list []QryPson) {
	mber, ok := wechat.M_member[par]
	if ok {
		if mber.IsFriend {
			list = append(list, QryPson{
				mber.UserName,
				mber.NickName,
				true})
			return
		}
		code = 1
		list = append(list, QryPson{
			mber.UserName,
			mber.NickName,
			false})
		return
	}
	var countFriend int
	code = -1 // 未找到
	for username, mber := range wechat.M_member {
		if strings.Contains(username, par) {
			if mber.IsFriend {
				countFriend += 1
				list = append(list, QryPson{
					mber.UserName,
					mber.NickName,
					true})
			} else {
				list = append(list, QryPson{
					mber.UserName,
					mber.NickName,
					false})
			}
			continue
		}
		if strings.Contains(mber.NickName, par) {
			if mber.IsFriend {
				countFriend += 1
				list = append(list, QryPson{
					mber.UserName,
					mber.NickName,
					true})
			} else {
				list = append(list, QryPson{
					mber.UserName,
					mber.NickName,
					false})
			}
			continue
		}
	}
	if len(list) == 0 {
		return
	} else if len(list) == 1 { // 找到唯一
		if countFriend == 1 {
			code = 0
			return
		} else {
			code = 1
			return
		}
	}
	code = 2
	return
}

type QryPson struct {
	UserName string
	NickName string
	IsFriend bool
}

// -1未找到 0有找到
// 群查询
func ChatroomQry(pars []string) (code int, list []QryPson) {
	if len(pars) == 0 {
		for _, mber := range wechat.M_member {
			if len(mber.UserName) > 2 && mber.UserName[:2] == "@@" {
				list = append(list, QryPson{
					mber.UserName,
					mber.NickName,
					false})
			}
		}
	} else {
		for _, mber := range wechat.M_member {
			if len(mber.UserName) > 2 && mber.UserName[:2] == "@@" {
				if strings.Contains(mber.UserName, pars[0]) {
					list = append(list, QryPson{
						mber.UserName,
						mber.NickName,
						false})
					continue
				}
				if strings.Contains(mber.NickName, pars[0]) {
					list = append(list, QryPson{
						mber.UserName,
						mber.NickName,
						false})
				}
			}
		}
	}
	if len(list) == 0 {
		code = -1
	}
	return
}
