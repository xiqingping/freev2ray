package main

import (
	_ "embed"
	"io/ioutil"
	"log"
	"os"

	"github.com/guonaihong/clop"
	"github.com/xiqingping/freev2ray"
	"github.com/xiqingping/freev2ray/fetcher"
)

//go:embed default_config.json
var defaultConfig string

type B64FetcherArgs struct {
	URL   string `clop:"short;long" usage:"the URL of file base64 encoded torjan list"`
	Index int    `clop:"short;long" usage:"the index of the trojan used"`
}

type ZKQFetcherArg struct {
}

type FreessFetcherArg struct {
	URL string `clop:"short;long" usage:"the mirror of https://free-ss.site"`
}

type SSFreeFetcherArg struct {
}

type MickysshVmessFetcherArg struct {
	Index int `clop:"short;long" usage:"the index of the vmess used"`
}

type CommandURLArg struct {
	URL string `clop:"short;long" usage:"vmess or trojan URL"`
}

type FetchArgs struct {
	B64vmess      B64FetcherArgs          `clop:"subcommand=b64vmess" usage:"vmess outbound, use base64 fetcher"`
	B64trojan     B64FetcherArgs          `clop:"subcommand=b64trojan" usage:"trojan outbound, use base64 fetcher"`
	B64SS         B64FetcherArgs          `clop:"subcommand=b64ss" usage:"ss outbound, use base64 fetcher"`
	ZKQ           ZKQFetcherArg           `clop:"subcommand=zkqtrojan" usage:"trojan outbound, use zkq fetcher"`
	Freess        FreessFetcherArg        `clop:"subcommand=freessvmess" usage:"vmess outbound, use freess fetcher(https://free-ss.site)"`
	SSFree        SSFreeFetcherArg        `clop:"subcommand=ssfreevmess" usage:"vmess outbound, use freess fetcher(https://view.ssfree.ru)"`
	Mickyssh      MickysshVmessFetcherArg `clop:"subcommand=mickysshvmess" usage:"vmess outbound, use mickyssh fetcher"`
	CmdURL        CommandURLArg           `clop:"subcommand=cmdurl" usage:"vmess or trojan outbound, use URL from command line"`
	DefaultConfig string                  `clop:"short;long" usage:"the default config file"`
}

func CreateFetcherByCmdLine() freev2ray.OutboundInfoFetcher {

	args := &FetchArgs{}
	clop.Bind(args)
	if args.DefaultConfig != "" {
		if content, err := ioutil.ReadFile(args.DefaultConfig); err == nil {
			defaultConfig = string(content)
		} else {
			log.Fatal(err)
		}
	}

	if clop.IsSetSubcommand("b64vmess") {
		return fetcher.NewBase64VmessFetcher(args.B64vmess.URL, args.B64vmess.Index)
	} else if clop.IsSetSubcommand("b64trojan") {
		return fetcher.NewBase64TrojanFetcher(args.B64trojan.URL, args.B64trojan.Index)
	} else if clop.IsSetSubcommand("b64ss") {
		return fetcher.NewBase64SSFetcher(args.B64SS.URL, args.B64SS.Index)
	} else if clop.IsSetSubcommand("zkqtrojan") {
		return fetcher.ZKQTrojanFetcher{}
	} else if clop.IsSetSubcommand("freessvmess") {
		return fetcher.FreessVmessFetcher{MirrorURL: args.Freess.URL}
	} else if clop.IsSetSubcommand("ssfreevmess") {
		return fetcher.SSFreeVmessFetcher{}
	} else if clop.IsSetSubcommand("mickysshvmess") {
		return fetcher.MickysshVmessFetcher{Index: args.Mickyssh.Index}
	} else if clop.IsSetSubcommand("cmdurl") {
		return fetcher.StringURLFetcher{URL: args.CmdURL.URL}
	} else {
		clop.Usage()
		os.Exit(1)
		return nil
	}
}
