package fetcher

import (
	"encoding/json"
	"testing"
)

func TestBase64VmessFetcher(t *testing.T) {
	f := NewBase64VmessFetcher("", 0)

	cfg, d, err := f.Fetch()
	if err != nil {
		t.Fatal(err)
	}

	config, _ := json.MarshalIndent(cfg, "", "  ")

	t.Log("timeout:", d)
	t.Log("config:", string(config))
}

func TestBase64TrojanFetcher(t *testing.T) {
	f := NewBase64TrojanFetcher("", 0)

	cfg, d, err := f.Fetch()
	if err != nil {
		t.Fatal(err)
	}

	config, _ := json.MarshalIndent(cfg, "", "  ")

	t.Log("timeout:", d)
	t.Log("config:", string(config))
}
