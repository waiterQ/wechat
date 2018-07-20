package combo

import (
	"test/wechat"
	lerrors "test/wechat/errors"
)

// exitCh
func Logout() lerrors.Lerror {
	wechat.ManualQuit = true // 手动退出
	err2 := wechat.ReportLogout()

	err := wechat.Logout()
	if err != nil {
		return lerrors.Transform(err, lerrors.FATAL)
	}
	if err2 != nil {
		return lerrors.Transform(err2)
	}
	return nil
}
