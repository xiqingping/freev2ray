package freev2ray

import (
	"time"
)

//StringURLFetcher 从freess获取Vmess节点
type StringURLFetcher struct {
	URL string
}

// Fetch 从https://mickyssh.me/download获取免费V2ray节点信息
func (f StringURLFetcher) Fetch() (V2rayConfigMap, time.Duration, error) {
	duration := time.Hour * 10000

	if configMap, err := v2rayConfigFromVmessURL(f.URL); err == nil {
		return configMap, duration, nil
	} else if configMap, err := v2rayConfigFromTrojanURL(f.URL); err == nil {
		return configMap, duration, nil
	} else {
		return V2rayConfigMap{}, duration, nil
	}
}
