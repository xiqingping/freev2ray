package main

import (
	"bytes"
	_ "embed"
	"io/ioutil"
	"log"
	"os"

	"github.com/guonaihong/clop"
	"github.com/xiqingping/freev2ray"

	core "github.com/v2fly/v2ray-core/v4"
)

//go:embed default_config.json
var defaultConfig []byte

type B64FetcherArgs struct {
	Url   string `clop:"short;long" usage:"the url of file base64 encoded torjan list"`
	Index int    `clop:"short;long" usage:"the index of the trojan used"`
}

type ZKQFetcherArg struct {
}

type FreessFetcherArg struct {
	Url string `clop:"short;long" usage:"the mirror of https://free-ss.site"`
}
type FetchArgs struct {
	B64vmess      B64FetcherArgs   `clop:"subcommand=b64vmess" usage:"vmess outbound, use base64 fetcher"`
	B64trojan     B64FetcherArgs   `clop:"subcommand=b64trojan" usage:"trojan outbound, use base64 fetcher"`
	ZKQ           ZKQFetcherArg    `clop:"subcommand=zkqtrojan" usage:"trojan outbound, use zkq fetcher"`
	Freess        FreessFetcherArg `clop:"subcommand=freessvmess" usage:"vmess outbound, use freess fetcher"`
	DefaultConfig string           `clop:"short;long" usage:"the default config file"`
}

func startV2rayConfigRunner() <-chan []byte {
	var cfgJSONCh <-chan []byte
	args := &FetchArgs{}
	clop.Bind(args)
	if args.DefaultConfig != "" {
		if content, err := ioutil.ReadFile(args.DefaultConfig); err == nil {
			defaultConfig = content
		} else {
			log.Fatal(err)
		}
	}

	if clop.IsSetSubcommand("b64vmess") {
		cfgJSONCh = freev2ray.StartV2rayConfigRunner(freev2ray.NewBase64VmessFetcher(args.B64vmess.Url, args.B64vmess.Index), string(defaultConfig))
	} else if clop.IsSetSubcommand("b64trojan") {
		cfgJSONCh = freev2ray.StartV2rayConfigRunner(freev2ray.NewBase64TrojanFetcher(args.B64trojan.Url, args.B64trojan.Index), string(defaultConfig))
	} else if clop.IsSetSubcommand("zkqtrojan") {
		cfgJSONCh = freev2ray.StartV2rayConfigRunner(freev2ray.ZKQTrojanFetcher{}, string(defaultConfig))
	} else if clop.IsSetSubcommand("freessvmess") {
		cfgJSONCh = freev2ray.StartV2rayConfigRunner(freev2ray.FreessVmessFetcher{args.Freess.Url}, string(defaultConfig))
	} else {
		clop.Usage()
		os.Exit(1)
	}

	return cfgJSONCh
}

func serverLoop(cfgJSONCh <-chan []byte) {
	var server *core.Instance
	var cfgJSON []byte

	for {
		cfgJSON = <-cfgJSONCh
		log.Println("ConfigJSON:")
		log.Println(string(cfgJSON))

		if cfg, err := core.LoadConfig("json", "", bytes.NewReader(OSHookConfig(cfgJSON))); err != nil {
			log.Println("Failed to load config", err)
		} else {
			if server != nil {
				if err = server.Close(); err != nil {
					log.Println("Close v2ray error:", err)
				} else {
					log.Println("V2ray server closed")
				}
			}

			if server, err = core.New(cfg); err != nil {
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
