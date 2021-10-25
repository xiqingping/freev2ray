package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/xiqingping/freev2ray"
)

func main() {
	f, _ := os.Executable()
	os.Chdir(filepath.Dir(f))

	for idx, protocol := range gjson.Get(defaultConfig, "inbounds.#.protocol").Array() {
		if protocol.String() == "dokodemo-door" {
			defaultConfig, _ = sjson.Delete(defaultConfig, fmt.Sprintf("inbounds.%d", idx))
			break
		}
	}

	fetcher := CreateFetcherByCmdLine()
	freev2ray.ServerLoop(freev2ray.StartV2rayConfigRunner(fetcher, defaultConfig), nil)
}
