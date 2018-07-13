package combo

import (
	"errors"
	"strings"
	"test/wechat"
)

// 获取最近联系群聊
func GetLastChatrooms(syncResp *wechat.WebWxSyncResp) (chatroomsName []string, err error) {
	if len(syncResp.AddMsgList) == 0 {
		err = errors.New("没有消息")
		return
	}
	var have bool
	for i := 0; i < len(syncResp.AddMsgList); i++ {
		if syncResp.AddMsgList[i].MsgType == 51 &&
			wechat.Me_userName == syncResp.AddMsgList[i].FromUserName && wechat.Me_userName == syncResp.AddMsgList[i].ToUserName {
			have = true
			if syncResp.AddMsgList[i].StatusNotifyUserName != "" {
				allUsersName := strings.Split(syncResp.AddMsgList[i].StatusNotifyUserName, ",")
				for j := 0; j < len(allUsersName); j++ {
					if allUsersName[j][:2] == "@@" {
						chatroomsName = append(chatroomsName, allUsersName[j])
					}
				}
			}
			break
		}
	}
	if !have {
		err = errors.New("没有最近联系人相关消息")
		return
	}
	return
}
