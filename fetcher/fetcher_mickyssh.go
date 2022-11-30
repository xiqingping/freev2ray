package fetcher

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/xiqingping/freev2ray"
)

// https://mickyssh.me/download

// MickysshVmessFetcher 从freess获取Vmess节点
type MickysshVmessFetcher struct {
	Index int
}

// Fetch 从https://mickyssh.me/download获取免费V2ray节点信息
func (f MickysshVmessFetcher) Fetch() (freev2ray.V2rayConfigMap, time.Duration, error) {
	http := NewHttpClient()
	duration := time.Minute * 5

	rsp, err := http.Get("https://mickyssh.me/download")
	if err != nil {
		return nil, duration, err
	}

	doc, err := goquery.NewDocumentFromReader(rsp.Body)
	if err != nil {
		log.Println(err)
	}

	vmess := doc.Find(fmt.Sprintf("body > div.container > div.container > div.row > div:nth-child(%d) > div > div.card-body.pt-0", f.Index+1))
	if vmess == nil {
		err = errors.New("not found port")
		return nil, duration, err
	}
	rsp.Body.Close()

	log.Println("vmess text:", strings.TrimSpace(vmess.Text()))

	if configMap, err := v2rayConfigFromVmessURL(strings.TrimSpace(vmess.Text())); err != nil {
		return nil, duration, err
	} else {
		h, m, s := time.Now().Clock()
		duration = (time.Duration(23-h)*3600 + time.Duration(59-m)*60 + time.Duration(70-s)) * time.Second
		return configMap, duration, nil
	}
}
