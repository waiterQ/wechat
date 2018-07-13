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

var firstLoginTime time.Time
var FirstSendMsgTime time.Time
var HaveFirstSendMsg bool

// when open wx2.qq.com, before jslogin
func Report0() (err error) {
	rpt_url := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatreport?fun=new"
	SetCookiesForFirst(rpt_url)
	breq := BaseRequest{
		DeviceID: DeviceID(),
	}
	m := make(map[string]interface{})
	m["BaseRequest"] = breq
	m["Count"] = 2
	t1 := fmt.Sprint(time.Now().Add(-time.Second).UnixNano() / 1e6)
	firstLoginTime = time.Now()
	t2 := fmt.Sprint(firstLoginTime.UnixNano() / 1e6)
	m["List"] = []interface{}{ //
		map[string]interface{}{
			"Text": `{"type":"[app-runtime]","data":{"unload":{"listenerCount":117,"watchersCount":115,"scopesCount":30}}}`,
			"Type": 1,
		},
		map[string]interface{}{
			"Text": `{"type":"[app-timing]","data":{"appTiming":{"qrcodeStart":` + t2 + `,"qrcodeEnd":` + t2 + `},"pageTiming":{"navigationStart":` + t1 + `,"unloadEventStart":` + t1 + `,"unloadEventEnd":` + t1 + `,"redirectStart":0,"redirectEnd":0,"fetchStart":` + t1 + `,"domainLookupStart":` + t1 + `,"domainLookupEnd":` + t1 + `,"connectStart":` + t1 + `,"connectEnd":` + t1 + `,"secureConnectionStart":` + t1 + `,"requestStart":` + t1 + `,"responseStart":` + t1 + `,"responseEnd":` + t1 + `,"domLoading":` + t1 + `,"domInteractive":` + t2 + `,"domContentLoadedEventStart":` + t2 + `,"domContentLoadedEventEnd":` + t2 + `,"domComplete":` + t2 + `,"loadEventStart":` + t2 + `,"loadEventEnd":` + t2 + `,"timeToNonBlankPaint":` + t1 + `,"timeToDOMContentFlushed":` + t2 + `}}}`,
			"Type": 1,
		},
	}
	breqdata, _ := json.Marshal(m)
	req, _ := http.NewRequest("post", rpt_url, bytes.NewBuffer(breqdata))
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("Report0().[listCount=2].resp.StatusCode =" + fmt.Sprint(resp.StatusCode))
	}
	breq = BaseRequest{
		DeviceID: DeviceID(),
	}
	m = make(map[string]interface{})
	m["BaseRequest"] = breq
	m["Count"] = 0
	m["List"] = []string{}
	breqdata, _ = json.Marshal(m)

	req, _ = http.NewRequest("post", rpt_url, bytes.NewBuffer(breqdata))
	resp, err = Cli.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("Report0().[listCount=0].resp.StatusCode =" + fmt.Sprint(resp.StatusCode))
	}
	return
}

// after getcontact
func Report1() (err error) {
	rpt_url := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatreport?fun=new&pass_ticket=" + conf.PassTicket
	breq := BaseRequest{
		DeviceID: DeviceID(),
		Sid:      conf.Wxsid,
		Uin:      conf.Wxuin,
	}
	m := make(map[string]interface{})
	m["BaseRequest"] = breq
	m["Count"] = 1
	t1 := fmt.Sprint(firstLoginTime.UnixNano() / 1e6)
	t2 := fmt.Sprint(time.Now().Add(-3*time.Second).UnixNano() / 1e6)
	t3 := fmt.Sprint(time.Now().UnixNano() / 1e6)
	m["List"] = []interface{}{
		map[string]interface{}{
			"Text": `{"type":"[app-timing]","data":{"appTiming":{"qrcodeStart":` + t1 + `,"qrcodeEnd":` + t1 + `,"scan":` + t2 + `,"loginEnd":` + t3 + `,"initStart":` + t3 + `,"initEnd":` + t3 + `,"initContactStart":` + t3 + `},"pageTiming":{"navigationStart":` + t1 + `,"unloadEventStart":0,"unloadEventEnd":0,"redirectStart":0,"redirectEnd":0,"fetchStart":` + t1 + `,"domainLookupStart":` + t1 + `,"domainLookupEnd":` + t1 + `,"connectStart":` + t1 + `,"connectEnd":` + t1 + `,"secureConnectionStart":` + t1 + `,"requestStart":` + t1 + `,"responseStart":` + t1 + `,"responseEnd":` + t1 + `,"domLoading":` + t1 + `,"domInteractive":` + t1 + `,"domContentLoadedEventStart":` + t1 + `,"domContentLoadedEventEnd":` + t1 + `,"domComplete":` + t1 + `,"loadEventStart":` + t1 + `,"loadEventEnd":` + t1 + `,"timeToNonBlankPaint":` + t1 + `,"timeToDOMContentFlushed":` + t1 + `}}}`,
			"Type": 1,
		},
	}
	breqdata, _ := json.Marshal(m)
	req, _ := http.NewRequest("post", rpt_url, bytes.NewBuffer(breqdata))
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("Report1().resp.StatusCode =" + fmt.Sprint(resp.StatusCode))
	}
	return
}

// 登陆上之后，如果每10分钟内有发送消息，需要[发送框]
func ReportSendMsg() (err error) {
	rpt_url := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatreport?fun=new&pass_ticket=" + conf.PassTicket
	breq := BaseRequest{
		DeviceID: DeviceID(),
		Sid:      conf.Wxsid,
		Uin:      conf.Wxuin,
	}
	m := make(map[string]interface{})
	m["BaseRequest"] = breq
	m["Count"] = 1
	t1 := fmt.Sprint(FirstSendMsgTime.Add(-5*time.Second).UnixNano() / 1e6)
	t2 := fmt.Sprint(FirstSendMsgTime.UnixNano() / 1e6)
	m["List"] = []interface{}{
		map[string]interface{}{
			"Text": `{"type":"[action-record]","data":{"actions":[{"type":"click","action":"发送框","time":` + t1 + `},{"type":"keydown","action":"发送框-enter","time":` + t2 + `},{"type":"keydown","action":"发送框-enter","time":` + t2 + `}]}}`,
			"Type": 1,
		},
	}
	breqdata, _ := json.Marshal(m)
	req, _ := http.NewRequest("post", rpt_url, bytes.NewBuffer(breqdata))
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("ReportSendMsg().resp.StatusCode =" + fmt.Sprint(resp.StatusCode))
	}
	return
}

// 当天第二次登录
// after webwxgetbatchcontact
func Report2() (err error) {
	rpt_url := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatreport?fun=new"
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
	req, _ := http.NewRequest("post", rpt_url, bytes.NewBuffer(breqdata))
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
	rpt_url := "https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxstatreport?fun=new"
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
	req, _ := http.NewRequest("post", rpt_url, bytes.NewBuffer(breqdata))
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New("Report3().resp.StatusCode =" + fmt.Sprint(resp.StatusCode))
	}
	return
}
