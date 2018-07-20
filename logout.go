package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func Logout() (err error) {
	logout_url := CgiUrl + "/webwxlogout?redirect=1&type=0&skey=" + conf.Skey
	m := make(map[string]interface{})
	m["sid"] = conf.Wxsid
	m["uin"] = conf.Wxuin
	formdata, _ := json.Marshal(m)
	req, _ := http.NewRequest("post", logout_url, bytes.NewBuffer(formdata))
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 301 && resp.StatusCode != 200 { // webwx是301 我请求是200
		err = errors.New("Logout().resp.StatusCode != 301, =" + fmt.Sprint(resp.StatusCode))
	}
	return
}
