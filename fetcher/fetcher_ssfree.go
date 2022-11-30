package fetcher

import (
	"errors"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/xiqingping/freev2ray"
)

// https://view.ssfree.ru/

// SSFreeVmessFetcher 从freess获取Vmess节点
type SSFreeVmessFetcher struct {
}

// Fetch 从https://view.ssfree.ru获取免费V2ray节点信息
func (f SSFreeVmessFetcher) Fetch() (freev2ray.V2rayConfigMap, time.Duration, error) {
	http := NewHttpClient()
	duration := time.Minute * 5

	rsp, err := http.Get("https://view.ssfree.ru")
	if err != nil {
		return nil, duration, err
	}

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		log.Println(err)
	}

	vmessSelector := doc.Find("#btn")
	if vmessSelector == nil {
		err = errors.New("not found port")
		return nil, duration, err
	}
	rsp.Body.Close()

	vmess, ok := vmessSelector.Attr("data-clipboard-text")
	if !ok {
		return nil, duration, errors.New("vmess url not exsit")
	}

	if configMap, err := v2rayConfigFromVmessURL(vmess); err != nil {
		return nil, duration, err
	} else {
		h, m, s := time.Now().Clock()
		duration = (time.Duration(23-h)*3600 + time.Duration(59-m)*60 + time.Duration(70-s)) * time.Second

		if duration > time.Second*3600*12 {
			duration -= time.Second * 3600 * 12
		}
		return configMap, duration, nil
	}
}
