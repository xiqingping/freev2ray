package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/coreos/go-iptables/iptables"
	"github.com/thecodeteam/goodbye"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/xiqingping/freev2ray"
)

// IptablesConfig iptables配置
type IptablesConfig struct {
	Table  string
	Chain  string
	Policy string
}

func initIptables(configs []IptablesConfig) {
	tables, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
	if err != nil {
		return
	}

	for _, cfg := range configs {
		tables.NewChain(cfg.Table, cfg.Chain)
		tables.AppendUnique(cfg.Table, cfg.Chain, strings.Split(cfg.Policy, " ")...)
	}
}

func deinitIptables(configs []IptablesConfig) {
	tables, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
	if err != nil {
		return
	}

	for _, cfg := range configs {
		tables.Delete(cfg.Table, cfg.Chain, strings.Split(cfg.Policy, " ")...)
	}
}

var iptabesConfigs = []IptablesConfig{
	{Table: "nat", Chain: "V2RAY", Policy: "-p tcp -j RETURN -m mark --mark 0xff"},
	{Table: "nat", Chain: "V2RAY", Policy: "-d 10.0.0.0/8 -j RETURN"},
	{Table: "nat", Chain: "V2RAY", Policy: "-d 127.0.0.0/8 -j RETURN"},
	{Table: "nat", Chain: "V2RAY", Policy: "-d 169.254.0.0/16 -j RETURN"},
	{Table: "nat", Chain: "V2RAY", Policy: "-d 172.16.0.0/12 -j RETURN"},
	{Table: "nat", Chain: "V2RAY", Policy: "-d 192.168.0.0/16 -j RETURN"},
	{Table: "nat", Chain: "V2RAY", Policy: "-d 224.0.0.0/4 -j RETURN"},
	{Table: "nat", Chain: "V2RAY", Policy: "-d 240.0.0.0/4 -j RETURN"},
	{Table: "nat", Chain: "V2RAY", Policy: "-p tcp -j REDIRECT --to-ports 12345"},
	{Table: "nat", Chain: "PREROUTING", Policy: "-p tcp -j V2RAY"},
	{Table: "nat", Chain: "OUTPUT", Policy: "-p tcp -m mark ! --mark 0xff -j V2RAY"},
}

func main() {
	ctx := context.Background()
	defer goodbye.Exit(ctx, -1)
	goodbye.Notify(ctx)

	if os.Geteuid() == 0 {
		defaultConfig, _ = sjson.Set(defaultConfig, "outbounds.0.streamSettings.sockopt.mark", 255)
		initIptables(iptabesConfigs)
		goodbye.Register(func(ctx context.Context, s os.Signal) {
			deinitIptables(iptabesConfigs)
		})
	} else {
		defaultConfig, _ = sjson.Delete(defaultConfig, "outbounds.1.streamSettings")
		for idx, protocol := range gjson.Get(defaultConfig, "inbounds.#.protocol").Array() {
			if protocol.String() == "dokodemo-door" {
				defaultConfig, _ = sjson.Delete(defaultConfig, fmt.Sprintf("inbounds.%d", idx))
				break
			}
		}

	}

	fetcher := CreateFetcherByCmdLine()
	freev2ray.ServerLoop(freev2ray.StartV2rayConfigRunner(fetcher, defaultConfig), nil)
}
