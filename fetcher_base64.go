package freev2ray

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

func v2rayConfigFromVmessURL(vmessURL string) (V2rayConfigMap, error) {
	if !strings.HasPrefix(vmessURL, "vmess://") {
		return nil, errors.New("not vmess url")
	}

	vmessJSONBytes, err := base64.URLEncoding.DecodeString(vmessURL[8:])
	if err != nil {
		return nil, err
	}
	vmessJSON := string(vmessJSONBytes)

	addr := gjson.Get(vmessJSON, "add").String()

	if match, _ := regexp.Match("((25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])\\.){3}(25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9][0-9]|[0-9])", []byte(addr)); match {
	} else if match, _ := regexp.Match("([0-9A-Za-z\\-_\\.]+)\\.([0-9a-z]+\\.[a-z]{2,3}(\\.[a-z]{2})?)", []byte(addr)); match {
	} else {
		return nil, errors.New("addr not ip or domain:" + addr)
	}
	configMap := V2rayConfigMap{
		"outbounds.0.protocol":                         "vmess",
		"outbounds.0.settings.vnext.0.address":         gjson.Get(vmessJSON, "add").String(),
		"outbounds.0.settings.vnext.0.port":            int(gjson.Get(vmessJSON, "port").Int()),
		"outbounds.0.settings.vnext.0.users.0.id":      gjson.Get(vmessJSON, "id").String(),
		"outbounds.0.settings.vnext.0.users.0.alterId": int(gjson.Get(vmessJSON, "aid").Int()),
		"outbounds.0.streamSettings.network":           gjson.Get(vmessJSON, "net").String(),
	}

	if gjson.Get(vmessJSON, "tls").String() == "tls" {
		configMap["outbounds.0.streamSettings.security"] = "tls"
		configMap["outbounds.0.streamSettings.tlsSettings.allowInsecure"] = true
		if host := gjson.Get(vmessJSON, "host").String(); host != "" {
			configMap["outbounds.0.streamSettings.tlsSettings.serverName"] = host
		}
	}

	if net := gjson.Get(vmessJSON, "net").String(); net == "ws" {
		configMap["outbounds.0.streamSettings.wsSettings.connectionReuse"] = true
		configMap["outbounds.0.streamSettings.wsSettings.path"] = gjson.Get(vmessJSON, "path").String()
		configMap["outbounds.0.streamSettings.wsSettings.headers.Host"] = gjson.Get(vmessJSON, "host").String()
	} else if net == "tcp" {
		if gjson.Get(vmessJSON, "type").String() == "http" {
			return nil, errors.New("unsupported tcp with http")
		}
	} else {
		return nil, errors.New("unsupported net type " + net)
	}

	return configMap, nil
}

func v2rayConfigFromTrojanURL(trojanURL string) (V2rayConfigMap, error) {
	if !strings.HasPrefix(trojanURL, "trojan://") {
		return nil, errors.New(trojanURL + " not trojan url")
	}

	url := strings.Split(strings.Split(trojanURL[9:], "#")[0], "@")
	if len(url) != 2 {
		return nil, errors.New(trojanURL + "error trojan url format")
	}

	password := url[0]
	addrAndPort := strings.Split(url[1], ":")
	var addr string
	var port int
	if len(addrAndPort) == 1 {
		addr = addrAndPort[0]
		port = 443
	} else if len(addrAndPort) == 2 {
		addr = addrAndPort[0]
		var err error
		port, err = strconv.Atoi(addrAndPort[1])
		if err != nil {
			return nil, errors.New(trojanURL + "error trojan url format")
		}
	} else {
		return nil, errors.New(trojanURL + "error trojan url format")
	}

	return V2rayConfigMap{
		"outbounds.0.protocol":                                 "trojan",
		"outbounds.0.settings.servers.0.address":               addr,
		"outbounds.0.settings.servers.0.port":                  port,
		"outbounds.0.settings.servers.0.password":              password,
		"outbounds.0.settings.servers.0.level":                 0,
		"outbounds.0.streamSettings.security":                  "tls",
		"outbounds.0.streamSettings.tlsSettings.allowInsecure": true,
		"outbounds.0.streamSettings.tlsSettings.serverName":    addr,
	}, nil
}

func NewBase64VmessFetcher(url string, index int) *Base64Fetcher {
	return &Base64Fetcher{
		url:                url,
		index:              index,
		v2rayConfigFromURL: v2rayConfigFromVmessURL,
	}
}

func NewBase64TrojanFetcher(url string, index int) *Base64Fetcher {
	return &Base64Fetcher{
		url:                url,
		index:              index,
		v2rayConfigFromURL: v2rayConfigFromTrojanURL,
	}
}

type Base64Fetcher struct {
	url                string
	index              int
	v2rayConfigFromURL func(Url string) (V2rayConfigMap, error)
}

// Fetch 从网络上获取免费trojan节点信息
func (f *Base64Fetcher) Fetch() (V2rayConfigMap, time.Duration, error) {
	http := NewHttpClient()
	duration := time.Minute * 5

	if f.url == "" {
		f.url = "https://cdn.jsdelivr.net/gh/freefq/free@master/v2"
	}

	log.Println("Get url info from:", f.url)
	rsp, err := http.Get(f.url)
	if err != nil {
		return nil, duration, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()
	if err != nil {
		return nil, duration, err
	}

	dec1, err := base64.URLEncoding.DecodeString(string(body))
	if err != nil {
		return nil, duration, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(dec1))
	index := 0
	for scanner.Scan() {
		if info, err := f.v2rayConfigFromURL(scanner.Text()); err == nil {
			if index >= f.index {
				return info, duration, nil
			} else {
				index++
			}
		}
	}

	log.Println("No valid url")
	return nil, duration, nil
}
