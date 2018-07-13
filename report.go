package wechat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func SetCookiesForFirst(rawUrl string) {
	u, _ := url.Parse(rawUrl)
	c1 := http.Cookie{
		Name:       "MM_WX_NOTIFY_STATE",
		Value:      "1",
		Path:       "/",
		Domain:     u.Host,
		Expires:    time.Now().AddDate(1, 0, 0),
		RawExpires: time.Now().AddDate(1, 0, 0).Format(time.RFC1123),
	}
	c2 := http.Cookie{
		Name:       "MM_WX_SOUND_STATE",
		Value:      "1",
		Path:       "/",
		Domain:     u.Host,
		Expires:    time.Now().AddDate(1, 0, 0),
		RawExpires: time.Now().AddDate(1, 0, 0).Format(time.RFC1123),
	}
	c3 := http.Cookie{
		Name:       "mm_lang",
		Value:      "zh_CN",
		Path:       "/",
		Domain:     u.Host,
		Expires:    time.Now().AddDate(1, 0, 0),
		RawExpires: time.Now().AddDate(1, 0, 0).Format(time.RFC1123),
	}
	cs := []*http.Cookie{&c1, &c2, &c3}
	Cli.Jar.SetCookies(u, cs)
	return
}

// before https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid=IaKFGQz7Bw==&tip=0&r=1943993128&_=1531359321280
func Report1() (err error) {
	r1_url := "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxstatreport?fun=new"
	SetCookiesForFirst(r1_url)
	breq := BaseRequest{
		DeviceID: DeviceID(),
	}
	m := make(map[string]interface{})
	m["BaseRequest"] = breq
	m["Count"] = 0
	m["List"] = []string{}
	breqdata, _ := json.Marshal(m)

	req, _ := http.NewRequest("post", r1_url, bytes.NewBuffer(breqdata))
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("Report1().resp.StatusCode =" + fmt.Sprint(resp.StatusCode))
	}
	return
}

// after webwxgetbatchcontact
func Report2() (err error) {
	r2_url := CgiUrl + "/webwxstatreport?fun=new"
	SetCookiesForFirst(r2_url)
	breq := BaseRequest{
		conf.Wxuin,
		conf.Wxsid,
		"",
		DeviceID(),
	}
	m := make(map[string]interface{})
	m["BaseRequest"] = breq
	m["Count"] = 1
	m["List"] = []interface{}{ // todo
		map[string]interface{}{
			"Text": `{"type":"[app-runtime]","data":{"0":{"listenerCount":297,"watchersCount":430,"scopesCount":112},"15000":{"listenerCount":297,"watchersCount":430,"scopesCount":112},"600000":{"listenerCount":297,"watchersCount":433,"scopesCount":113},"unload":{"listenerCount":351,"watchersCount":688,"scopesCount":169}}}`,
			"Type": 1,
		},
	}
	breqdata, _ := json.Marshal(m)
	req, _ := http.NewRequest("post", r2_url, bytes.NewBuffer(breqdata))
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("Report2().resp.StatusCode =" + fmt.Sprint(resp.StatusCode))
	}
	return
}

// after report2
func Report3() {
	r2_url := CgiUrl + "/webwxstatreport?fun=new"
	SetCookiesForFirst(r2_url)
	breq := BaseRequest{
		conf.Wxuin,
		conf.Wxsid,
		"",
		DeviceID(),
	}
	t1 := time.Now().Add(-2*time.Second).UnixNano() / 1e6
	t2 := time.Now().Add(-time.Second).UnixNano() / 1e6
	t3 := time.Now().UnixNano() / 1e6
	m := make(map[string]interface{})
	m["BaseRequest"] = breq
	m["Count"] = 1
	m["List"] = []interface{}{
		map[string]interface{}{
			"Text": `{"type":"[app-timing]","data":{"appTiming":{"initStart":` + fmt.Sprint(t2) + `,"initEnd":` + fmt.Sprint(t3) + `,"initContactStart":` + fmt.Sprint(t3) + `},"pageTiming":{"navigationStart":` + fmt.Sprint(t1) + `,"unloadEventStart":0,"unloadEventEnd":0,"redirectStart":` + fmt.Sprint(t1) + `,"redirectEnd":` + fmt.Sprint(t1) + `,"fetchStart":` + fmt.Sprint(t1) + `,"domainLookupStart":` + fmt.Sprint(t1) + `,"domainLookupEnd":` + fmt.Sprint(t1) + `,"connectStart":` + fmt.Sprint(t1) + `,"connectEnd":` + fmt.Sprint(t1) + `,"secureConnectionStart":` + fmt.Sprint(t1) + `,"requestStart":` + fmt.Sprint(t1) + `,"responseStart":` + fmt.Sprint(t1) + `,"responseEnd":` + fmt.Sprint(t1) + `,"domLoading":` + fmt.Sprint(t1) + `,"domInteractive":` + fmt.Sprint(t2) + `,"domContentLoadedEventStart":` + fmt.Sprint(t2) + `,"domContentLoadedEventEnd":` + fmt.Sprint(t2) + `,"domComplete":` + fmt.Sprint(t2) + `,"loadEventStart":` + fmt.Sprint(t2) + `,"loadEventEnd":` + fmt.Sprint(t2) + `,"timeToNonBlankPaint":` + fmt.Sprint(t2) + `,"timeToDOMContentFlushed":` + fmt.Sprint(t2) + `}}}`,
			"Type": 1,
		},
	}
	breqdata, _ := json.Marshal(m)
	req, _ := http.NewRequest("post", r2_url, bytes.NewBuffer(breqdata))
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("Report3().resp.StatusCode =" + fmt.Sprint(resp.StatusCode))
	}
	return
}

func Reprot4() {

}
