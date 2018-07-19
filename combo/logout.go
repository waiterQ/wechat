package combo

import (
	"fmt"
	"test/wechat"
)

// exitCh
func Logout() error {
	err2 := wechat.ReportLogout()
	if err2 != nil {
		fmt.Println(err2)
	}
	return wechat.Logout()
}
