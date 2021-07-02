package freev2ray

import (
	"encoding/json"
	"testing"
)

func TestNewSSFreeVmessFetcher(t *testing.T) {
	f := NewSSFreeVmessFetcher()
	cfg, d, err := f.Fetch()
	if err != nil {
		t.Fatal(err)
	}

	config, _ := json.MarshalIndent(cfg, "", "  ")

	t.Log("timeout:", d)
	t.Log("config:", string(config))
}


func TestNewFreev2rayVmessFetcher(t *testing.T) {
	f := NewFreev2rayVmessFetcher()
	cfg, d, err := f.Fetch()
	if err != nil {
		t.Fatal(err)
	}

	config, _ := json.MarshalIndent(cfg, "", "  ")

	t.Log("timeout:", d)
	t.Log("config:", string(config))
}
