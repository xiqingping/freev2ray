package fetcher

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"time"

	"github.com/xiqingping/freev2ray"
)

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
	v2rayConfigFromURL func(url string) (freev2ray.V2rayConfigMap, error)
}

// Fetch 从网络上获取免费trojan节点信息
func (f *Base64Fetcher) Fetch() (freev2ray.V2rayConfigMap, time.Duration, error) {
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
