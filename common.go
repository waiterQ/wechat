package wechat

import (
	"net/http"
)

var Cli *http.Client
var conf *XmlConfig
var webInitConf *InitResp
var CgiUrl string
var m_member map[string]Member = make(map[string]Member)
var deviceID string = "e435753412807776"
