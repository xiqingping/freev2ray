package main

import (
	"github.com/xiqingping/freev2ray"
)

func main() {
	fetcher := CreateFetcherByCmdLine()
	freev2ray.ServerLoop(freev2ray.StartV2rayConfigRunner(fetcher, defaultConfig), func (cfgJSON []byte) []byte {
		return cfgJSON
	})
}
