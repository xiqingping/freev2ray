package fetcher

import (
	"encoding/json"
	"testing"
)

func TestMickysshVmessFetcher(t *testing.T) {
	f := MickysshVmessFetcher{Index: 0}

	cfg, d, err := f.Fetch()
	if err != nil {
		t.Fatal(err)
	}

	config, _ := json.MarshalIndent(cfg, "", "  ")

	t.Log("timeout:", d)
	t.Log("config:", string(config))
}
