package freev2ray

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"time"

	"github.com/tidwall/gjson"
)

// https://free.v2x-nav.ml/

type V2xrayVmessFetcher struct{}

// Fetch 从网络上获取免费V2ray节点信息
func (V2xrayVmessFetcher) Fetch() (V2rayConfigMap, time.Duration, error) {
	httpClient := NewHttpClient()

	h, m, s := time.Now().Clock()
	duration := (time.Duration(23-h)*3600 + time.Duration(59-m)*60 + time.Duration(70-s)) * time.Second

	rsp, err := httpClient.Get("https://v2x-ray.com/impl/getHourFreeNodes")
	if err != nil {
		return nil, duration, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, duration, err
	}

	dec1, err := base64.RawStdEncoding.DecodeString(string(body))
	if err != nil {
		log.Println("base64 decode1:", string(body), "with error", err)
		return nil, duration, err
	}

	dec2, err := base64.URLEncoding.DecodeString(string(dec1))
	if err != nil {
		log.Println("base64 decode2:", string(body), "with error", err)
		return nil, duration, err
	}

	freeNodes := string(dec2)
	configMap, err := v2rayConfigFromVmessURL(gjson.Get(freeNodes, "2.shareLink").String())

	return configMap, duration, err
}
