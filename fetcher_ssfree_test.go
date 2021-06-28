package freev2ray

import (
	"encoding/json"
	"testing"
)

func TestSSFreeVmessFetcher(t *testing.T) {
	f := SSFreeVmessFetcher{}
	cfg, d, err := f.Fetch()
	if err != nil {
		t.Fatal(err)
	}

	config, _ := json.MarshalIndent(cfg, "", "  ")

	t.Log("timeout:", d)
	t.Log("config:", string(config))
}
