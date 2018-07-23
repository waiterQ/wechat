package wechat

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var qtpath string
var cookies []*http.Cookie

var AppID string = "wx782c26e4c19acffb" // 似乎所有的webwx的AppID都是一样的

func init() {
	transport := *(http.DefaultTransport.(*http.Transport))
	transport.ResponseHeaderTimeout = 1 * time.Minute
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	jar, _ := cookiejar.New(nil)
	Cli = &http.Client{
		Transport: &transport,
		Jar:       jar,
		Timeout:   1 * time.Minute,
	}
}

func GetLoginUuid() (loginuuid string, err error) {
	params := url.Values{}
	params.Set("appid", AppID)
	params.Set("fun", "new")
	params.Set("lang", "zh_CN")
	params.Set("_", strconv.FormatInt(time.Now().Unix(), 10))

	loginjs_url := "https://login.weixin.qq.com/jslogin"
	resp, err := Cli.PostForm(loginjs_url, params)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resultStr := string(respBody)
	fmt.Println("GetLoginUuid().resultStr=", resultStr)
	re := regexp.MustCompile(`window.QRLogin.code = (\d+); window.QRLogin.uuid = "(\S+?)"`)
	pm := re.FindStringSubmatch(resultStr)
	fmt.Printf("pm=%v\n", pm)

	if len(pm) != 3 {
		return "", errors.New(fmt.Sprintf("no loginuuid, len(pm)=%d", len(pm)))
	} else {
		if pm[1] != "200" {
			return "", errors.New("status != 200")
		}
		return pm[2], nil
	}
}

func DownloadImgQR(loginuuid string) (err error) {
	params := url.Values{}
	params.Set("t", "webwx")
	params.Set("_", strconv.FormatInt(time.Now().Unix(), 10))

	qrUrl := "https://login.weixin.qq.com/qrcode/"

	req, _ := http.NewRequest("POST", qrUrl+loginuuid, strings.NewReader(params.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := Cli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	qtpath = TmpPath + "/" + loginuuid + ".jpg"
	file, _ := os.OpenFile(qtpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer file.Close()
	file.Write(data)
	return nil
}

func GetLoginUrlAfterAuth(uuid string) (redirect_uri, code string, err error) {
	loginUrl := fmt.Sprintf("https://login.weixin.qq.com/cgi-bin/mmwebwx-bin/login?tip=0&uuid=%s&_=%d", uuid, time.Now().Unix())
	// tip等于0或1 似乎没有什么影响
	// rt = tip
	resp, err := Cli.Get(loginUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	resultStr := string(data)
	fmt.Printf("LoginAfterAuth(%s).resultStr=%s\n", uuid, resultStr)
	re := regexp.MustCompile(`window.code=(\d+);`)
	pm := re.FindStringSubmatch(resultStr)
	if len(pm) != 0 {
		code = pm[1]
	} else {
		err = errors.New("can't find the code")
		return
	}
	// rt = 0
	switch code {
	case "201":
		fmt.Println("扫描成功，请在手机上点击确认登陆")
	case "200":
		reRedirect := regexp.MustCompile(`window.redirect_uri="(\S+?)"`)
		pmSub := reRedirect.FindStringSubmatch(resultStr)
		if len(pmSub) != 0 {
			redirect_uri = pmSub[1]
			u, _ := url.Parse(redirect_uri)
			CgiUrl = u.Scheme + "://" + u.Host + "/cgi-bin/mmwebwx-bin"
			for i := 0; i < len(URLPool); i++ { // 设置sync同步消息用
				if URLPool[i].IndexUrl == u.Host {
					syncHost = u.Host
					syncUrl = u.Scheme + "://" + URLPool[i].SyncUrl + "/cgi-bin/mmwebwx-bin/synccheck"
					break
				}
			}
		} else {
			err = errors.New("regex error in window.redirect_uri")
			return
		}
		redirect_uri += "&fun=new" // 没有这个会返回一个页面
		fmt.Printf("os.Remove(%s) %v\n", qtpath, os.Remove(qtpath))
	case "400":
		err = errors.New("超时，二维码失效")
		fmt.Printf("os.Remove(%s) %v\n", qtpath, os.Remove(qtpath))
	case "408":
	// case "0":
	// 	err = errors.New("超时了")
	default:
		err = errors.New("未知错误")
	}
	return
}

func Login(url string) (err error) {
	resp, err := Cli.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	cs := resp.Cookies()
	cookies = append(cookies, cs...)
	config := XmlConfig{}

	if err = xml.NewDecoder(resp.Body.(io.Reader)).Decode(&config); err != nil {
		return
	}
	conf = &config
	return
}
