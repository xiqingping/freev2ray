package freev2ray

import (
	"log"
	"reflect"
	"time"

	"github.com/tidwall/sjson"
	v2ray "github.com/v2fly/v2ray-core/v5"
)

// V2rayConfigMap V2ray config
type V2rayConfigMap map[string]interface{}

// OutboundInfoFetcher 免费节点抓取接口
type OutboundInfoFetcher interface {
	Fetch() (V2rayConfigMap, time.Duration, error)
}

func updateConfig(defConfig string, configMap V2rayConfigMap) (string, error) {
	cfg := defConfig
	var err error
	for k, v := range configMap {
		cfg, err = sjson.Set(cfg, k, v)
		if err != nil {
			break
		}
	}
	return cfg, err
}

// v2rayConfigRunner 从网络上获取免费V2ray节点，并生成V2ray配置文件
func v2rayConfigRunner(fetcher OutboundInfoFetcher, defConfig string, ch chan<- []byte) {
	var orgInfo V2rayConfigMap

	if fetcher == nil {
		ch <- []byte(defConfig)
		return
	}

	for {
		info, d, err := fetcher.Fetch()
		if err == nil {
			if !reflect.DeepEqual(info, orgInfo) {
				if cfg, err := updateConfig(defConfig, info); err == nil {
					ch <- []byte(cfg)
					orgInfo = info
				} else {
					log.Println("Update config error", err)
				}
			}
		} else {
			log.Println("Get node info error", err)
		}
		time.Sleep(d)
	}
}

// StartV2rayConfigRunner 从网络上获取免费V2ray节点，并生成V2ray配置文件
func StartV2rayConfigRunner(fetcher OutboundInfoFetcher, defConfig string) <-chan []byte {
	ch := make(chan []byte)
	go v2rayConfigRunner(fetcher, defConfig, ch)
	return ch
}

func ServerLoop(cfgJSONCh <-chan []byte, hook func(cfgJSON []byte) []byte) {
	var server *v2ray.Instance
	var cfgJSON []byte

	for {
		cfgJSON = <-cfgJSONCh
		if hook != nil {
			cfgJSON = hook(cfgJSON)
		}

		log.Println("ConfigJSON:")
		log.Println(string(cfgJSON))

		if cfg, err := v2ray.LoadConfig("json", cfgJSON); err != nil {
			log.Println("Failed to load config", err)
		} else {
			if server != nil {
				if err = server.Close(); err != nil {
					log.Println("Close v2ray error:", err)
				} else {
					log.Println("V2ray server closed")
				}
			}

			if server, err = v2ray.New(cfg); err != nil {
				log.Println("Failed to create server", err)
			} else {
				log.Println("Server created")
				if err = server.Start(); err != nil {
					log.Println("Failed to start", err)
				} else {
					log.Println("Server started")
				}
			}
		}
	}
}
