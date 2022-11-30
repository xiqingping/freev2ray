package fetcher

import (
	"encoding/json"
	"testing"
)

func TestZKQTrojanFetcher(t *testing.T) {
	f := ZKQTrojanFetcher{}

	cfg, d, err := f.Fetch()
	if err != nil {
		t.Fatal(err)
	}

	config, _ := json.MarshalIndent(cfg, "", "  ")

	t.Log("timeout:", d)
	t.Log("config:", string(config))
}
