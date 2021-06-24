package freev2ray

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/tidwall/gjson"
)

// 使用国外邮箱发送邮件到 ss@rohankdd.com 可自动获取镜像站。

//FreessVmessFetcher 从freess获取Vmess节点
type FreessVmessFetcher struct {
	MirrorURL string
}

//https://free-ss.pw/v/443.json

// Fetch 从网络上获取免费V2ray节点信息
func (f FreessVmessFetcher) Fetch() (V2rayConfigMap, time.Duration, error) {
	http := NewHttpClient()
	duration := time.Minute * 5

	if f.MirrorURL == "" {
		f.MirrorURL = "https://free-ss.pw"
	}

	jsonURL := f.MirrorURL + "/v/443.json"

	log.Println("Get url info from:", jsonURL)
	rsp, err := http.Get(jsonURL)
	if err != nil {
		return nil, duration, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()
	if err != nil {
		return nil, duration, err
	}

	configMap := V2rayConfigMap{
		"outbounds.0.protocol": "vmess",
	}
	k := "outbounds.0.settings.vnext.0.address"
	configMap[k] = gjson.Get(string(body), k).String()
	k = "outbounds.0.settings.vnext.0.port"
	configMap[k] = gjson.Get(string(body), k).Int()
	k = "outbounds.0.settings.vnext.0.users.0.id"
	configMap[k] = gjson.Get(string(body), k).String()
	k = "outbounds.0.settings.vnext.0.users.0.alterId"
	configMap[k] = gjson.Get(string(body), k).Int()

	k = "outbounds.0.streamSettings.network"
	configMap[k] = gjson.Get(string(body), k).String()
	if configMap[k] == "ws" {
		k = "outbounds.0.streamSettings.wsSettings.path"
		configMap[k] = gjson.Get(string(body), k).String()
		k = "outbounds.0.streamSettings.wsSettings.headers.Host"
		configMap[k] = gjson.Get(string(body), k).String()
	}

	k = "outbounds.0.streamSettings.security"
	configMap[k] = gjson.Get(string(body), k).String()
	if configMap[k] == "tls" {
		k = "outbounds.0.streamSettings.tlsSettings.serverName"
		configMap[k] = gjson.Get(string(body), k).String()
	}

	h, m, s := time.Now().Clock()
	duration = (time.Duration(23-h)*3600 + time.Duration(59-m)*60 + time.Duration(70-s)) * time.Second

	return configMap, duration, nil

}
