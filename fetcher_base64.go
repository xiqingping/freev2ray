package freev2ray

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

func base64Decode(b64s string) ([]byte, error) {
	if ret, err := base64.URLEncoding.DecodeString(b64s); err == nil {
		return ret, err
	} else if ret, err := base64.RawURLEncoding.DecodeString(b64s); err == nil {
		return ret, err
	} else if ret, err := base64.StdEncoding.DecodeString(b64s); err == nil {
		return ret, err
	} else if ret, err := base64.RawStdEncoding.DecodeString(b64s); err == nil {
		return ret, err
	} else {
		return nil, err
	}
}

func v2rayConfigFromSSURL(url string) (V2rayConfigMap, error) {
	if !strings.HasPrefix(url, "ss://") {
		return nil, errors.New("not a ss url")
	}
	urls := strings.Split(strings.Split(url[5:], "#")[0], "@")
	if len(urls) != 2 {
		return nil, errors.New("error ss url format")
	}

	sMethodAndPassword, err := base64Decode(urls[0])
	if err != nil {
		return nil, errors.New("decode method and password error")
	}

	methodAndPassword := strings.Split(string(sMethodAndPassword), ":")
	if len(methodAndPassword) != 2 {
		return nil, errors.New("error method and password format")
	}

	addrAndPort := strings.Split(urls[1], ":")
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
			return nil, errors.New("error trojan port format")
		}
	} else {
		return nil, errors.New("error trojan addr and port format")
	}

	return V2rayConfigMap{
		"outbounds.0.protocol":                    "shadowsocks",
		"outbounds.0.settings.servers.0.email":    "love@v2ray.com",
		"outbounds.0.settings.servers.0.address":  addr,
		"outbounds.0.settings.servers.0.port":     port,
		"outbounds.0.settings.servers.0.method":   methodAndPassword[0],
		"outbounds.0.settings.servers.0.password": methodAndPassword[1],
		"outbounds.0.settings.servers.0.level":    0,
		"outbounds.0.settings.servers.0.ota":      false,
		"outbounds.0.streamSettings.security":     "none",
		"outbounds.0.streamSettings.network":      "tcp",
	}, nil
}

func v2rayConfigFromVmessURL(url string) (V2rayConfigMap, error) {
	if !strings.HasPrefix(url, "vmess://") {
		return nil, errors.New("not vmess url")
	}

	vmessJSONBytes, err := base64Decode(url[8:])
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
		if typ := gjson.Get(vmessJSON, "type").String(); typ == "http" {
			configMap["outbounds.0.streamSettings.tcpSettings.header.type"] = "http"
			if hosts := strings.ReplaceAll(gjson.Get(vmessJSON, "host").String(), " ", ""); hosts != "" {
				configMap["outbounds.0.streamSettings.tcpSettings.header.request.headers.Host"] = strings.Split(hosts, ",")
			}
			if paths := strings.ReplaceAll(gjson.Get(vmessJSON, "path").String(), " ", ""); paths != "" {
				configMap["outbounds.0.streamSettings.tcpSettings.header.request.path"] = strings.Split(paths, ",")
			} else {
				configMap["outbounds.0.streamSettings.tcpSettings.header.request.path"] = []string{"/"}
			}

			configMap["outbounds.0.streamSettings.tcpSettings.header.request.version"] = "1.1"
			configMap["outbounds.0.streamSettings.tcpSettings.header.request.method"] = "GET"
			configMap["outbounds.0.streamSettings.tcpSettings.header.request.headers.User-Agent"] = []string{
				"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.75 Safari/537.36",
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:57.0) Gecko/20100101 Firefox/57.0",
			}
			configMap["outbounds.0.streamSettings.tcpSettings.header.request.headers.Accept"] = []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"}
			configMap["outbounds.0.streamSettings.tcpSettings.header.request.headers.Accept-language"] = []string{"zh-CN,zh;q=0.8,en-US;q=0.6,en;q=0.4"}
			configMap["outbounds.0.streamSettings.tcpSettings.header.request.headers.Accept-Encoding"] = []string{"gzip, deflate, br"}
			configMap["outbounds.0.streamSettings.tcpSettings.header.request.headers.Cache-Control"] = []string{"no-cache"}
			configMap["outbounds.0.streamSettings.tcpSettings.header.request.headers.Connection"] = []string{"keep-alive"}
			configMap["outbounds.0.streamSettings.tcpSettings.header.request.headers.Pragma"] = "no-cache"
		} else if typ != "" {
			return nil, errors.New("unsupported tcp with " + typ)
		}
	} else {
		return nil, errors.New("unsupported net type " + net)
	}

	return configMap, nil
}

func v2rayConfigFromTrojanURL(url string) (V2rayConfigMap, error) {
	if !strings.HasPrefix(url, "trojan://") {
		return nil, errors.New("not trojan url")
	}

	urls := strings.Split(strings.Split(url[9:], "#")[0], "@")
	if len(urls) != 2 {
		return nil, errors.New("error trojan url format")
	}

	password := urls[0]
	addrAndPort := strings.Split(urls[1], ":")
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
			return nil, errors.New("error trojan port format")
		}
	} else {
		return nil, errors.New("error trojan addr and port format")
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

func NewBase64SSFetcher(url string, index int) *Base64Fetcher {
	return &Base64Fetcher{
		url:                url,
		index:              index,
		v2rayConfigFromURL: v2rayConfigFromSSURL,
	}
}

type Base64Fetcher struct {
	url                string
	index              int
	v2rayConfigFromURL func(url string) (V2rayConfigMap, error)
}

// Fetch ????????????????????????trojan????????????
func (f *Base64Fetcher) Fetch() (V2rayConfigMap, time.Duration, error) {
	http := NewHttpClient()
	duration := time.Minute

	if f.url == "" {
		f.url = "https://raw.fastgit.org/freefq/free/master/v2"
	}

	rsp, err := http.Get(f.url)
	if err != nil {
		return nil, duration, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	rsp.Body.Close()
	if err != nil {
		return nil, duration, err
	}

	dec1, err := base64Decode(string(body))
	if err != nil {
		return nil, duration, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(dec1))
	index := 0
	for scanner.Scan() {
		txt := scanner.Text()
		if info, err := f.v2rayConfigFromURL(txt); err == nil {
			if index >= f.index {
				return info, duration, nil
			} else {
				index++
			}
		} else {
			// log.Println("v2ray config from url error:", err, "with URL:", txt)
		}
	}

	return nil, duration, errors.New("no valid url")
}
