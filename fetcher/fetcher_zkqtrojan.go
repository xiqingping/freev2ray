package fetcher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/xiqingping/freev2ray"
)

// Documents
// https://www.youhou8.com/scripts/max/%E3%80%90%E5%85%8D%E8%B4%B9%E5%88%86%E4%BA%AB%E3%80%91%E7%A7%91%E5%AD%A6%E4%B8%8A%E7%BD%91%EF%BC%8C%E5%85%8D%E8%B4%B9%E8%8A%82%E7%82%B9%EF%BC%8C%E5%8F%AF%E8%A7%82%E7%9C%8B4K_YouTube%E8%A7%86%E9%A2%91%EF%BC%8C%E4%B8%8Agoogle%E6%9F%A5%E8%B5%84%E6%96%99%EF%BC%8CTrojan_%E8%B4%A6%E5%8F%B7%E5%88%86%E4%BA%AB

type ZKQTrojanFetcher struct {
}

func (f ZKQTrojanFetcher) Fetch() (freev2ray.V2rayConfigMap, time.Duration, error) {
	http := NewHttpClient()

	duration := time.Minute * 2

	resp, err := http.Get(fmt.Sprintf("https://zkq8.com/v2ray2_link.txt?t=%d", rand.Uint32()))
	if err != nil {
		return nil, duration, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
	}

	shost := doc.Find("body > table > tbody > tr:nth-child(2) > td:nth-child(3) > b")
	if shost == nil {
		err = errors.New("not found host")
		return nil, duration, err
	}

	sport := doc.Find("body > table > tbody > tr:nth-child(2) > td:nth-child(5) > b")
	if sport == nil {
		err = errors.New("not found port")
		return nil, duration, err
	}

	resp.Body.Close()

	resp, err = http.Get(fmt.Sprintf("https://zkq8.com//trojan_pwd1.txt?t=%d", rand.Uint32()))
	if err != nil {
		return nil, duration, err
	}

	password, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, duration, err
	}
	resp.Body.Close()

	port, _ := strconv.Atoi(sport.Text())
	addr := shost.Text()

	h, m, s := time.Now().Clock()
	duration = (time.Duration(23-h)*3600 + time.Duration(59-m)*60 + time.Duration(70-s)) * time.Second

	return freev2ray.V2rayConfigMap{
		"outbounds.0.protocol":                                 "trojan",
		"outbounds.0.settings.servers.0.address":               addr,
		"outbounds.0.settings.servers.0.port":                  port,
		"outbounds.0.settings.servers.0.password":              strings.TrimSpace(string(password)),
		"outbounds.0.settings.servers.0.level":                 0,
		"outbounds.0.streamSettings.security":                  "tls",
		"outbounds.0.streamSettings.tlsSettings.allowInsecure": true,
		"outbounds.0.streamSettings.tlsSettings.serverName":    addr,
	}, duration, nil

}
