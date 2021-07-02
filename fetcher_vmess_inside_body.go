package freev2ray

import (
	"errors"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// #intro > div > div > footer > ul:nth-child(1) > li:nth-child(2) > button

//VmessInsideBodyFetcher 从某个网址的Body中获取Vmess节点
type VmessInsideBodyFetcher struct {
	URL       string
	ParseArg  interface{}
	ParseFunc func(*goquery.Document, interface{}) (string, time.Duration, error)
}

// Fetch 从网址的Body中获取Vmess节点
func (f VmessInsideBodyFetcher) Fetch() (V2rayConfigMap, time.Duration, error) {
	http := NewHttpClient()
	duration := time.Minute * 5

	rsp, err := http.Get(f.URL)
	if err != nil {
		return nil, duration, err
	}

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		log.Println(err)
	}
	url, d, err := f.ParseFunc(doc, f.ParseArg)
	rsp.Body.Close()

	if err != nil {
		return nil, duration, err
	}

	if configMap, err := v2rayConfigFromVmessURL(url); err == nil {
		return configMap, d, nil
	} else if configMap, err = v2rayConfigFromVmessURL(url); err == nil {
		return configMap, d, nil
	} else {
		return nil, duration, err
	}
}

// NewSSFreeVmessFetcher https://view.ssfree.ru vmess Fetcher
func NewSSFreeVmessFetcher() VmessInsideBodyFetcher {
	return VmessInsideBodyFetcher{
		URL:      "https://view.ssfree.ru",
		ParseArg: nil,
		ParseFunc: func(doc *goquery.Document, arg interface{}) (string, time.Duration, error) {
			selector := doc.Find("#btn")
			if selector == nil {
				return "", 0, errors.New("vmess selector not exsit")
			}

			vmess, ok := selector.Attr("data-clipboard-text")
			if !ok {
				return "", 0, errors.New("vmess url not exsit")
			}

			h, m, s := time.Now().Clock()
			d := (time.Duration(23-h)*3600 + time.Duration(59-m)*60 + time.Duration(70-s)) * time.Second

			if d > time.Second*3600*12 {
				d -= time.Second * 3600 * 12
			}
			return vmess, d, nil
		},
	}
}

// NewFreev2rayVmessFetcher https://view.freev2ray.org vmess Fetcher
func NewFreev2rayVmessFetcher() VmessInsideBodyFetcher {
	return VmessInsideBodyFetcher{
		URL:      "https://view.freev2ray.org",
		ParseArg: nil,
		ParseFunc: func(doc *goquery.Document, arg interface{}) (string, time.Duration, error) {
			selector := doc.Find("#intro > div > div > footer > ul:nth-child(1) > li:nth-child(2) > button")
			if selector == nil {
				return "", 0, errors.New("vmess selector not exsit")
			}

			vmess, ok := selector.Attr("data-clipboard-text")
			if !ok {
				return "", 0, errors.New("vmess url not exsit")
			}

			h, m, s := time.Now().Clock()
			d := (time.Duration(23-h)*3600 + time.Duration(59-m)*60 + time.Duration(70-s)) * time.Second

			if d > time.Second*3600*12 {
				d -= time.Second * 3600 * 12
			}
			return vmess, d, nil
		},
	}
}
