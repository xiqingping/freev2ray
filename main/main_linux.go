package main

import (
	"context"
	"os"
	"strings"

	"github.com/coreos/go-iptables/iptables"
	"github.com/thecodeteam/goodbye"
	"github.com/tidwall/sjson"
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

func OSHookConfig(cfgJSON []byte) []byte {
	if cfg, err := sjson.Set(string(cfgJSON), "outbounds.0.streamSettings.sockopt.mark", 255); err != nil {
		return cfgJSON
	} else {
		return []byte(cfg)
	}
}

func main() {
	cfgJSONCh := startV2rayConfigRunner()

	ctx := context.Background()
	defer goodbye.Exit(ctx, -1)
	goodbye.Notify(ctx)
	initIptables(iptabesConfigs)
	goodbye.Register(func(ctx context.Context, s os.Signal) {
		deinitIptables(iptabesConfigs)
	})

	serverLoop(cfgJSONCh)
}
