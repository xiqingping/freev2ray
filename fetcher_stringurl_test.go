package freev2ray

import (
	"encoding/json"
	"testing"
)

func TestStringURLFetcher(t *testing.T) {
	f := StringURLFetcher{URL: "vmess://ewogICJob3N0IjoiIiwKICAicGF0aCI6Ii9nZXR3ZWF0aGVyIiwKICAicG9ydCI6IjQ0MyIsCiAgInRscyI6InRscyIsCiAgInBzIjoiXHU2NzAwXHU2NWIwXHU1NzMwXHU1NzQwXHU1M2QxXHU5MGFlXHU0ZWY2IGNjQGJic3MubWwiLAogICJpZCI6IjUwNTBjMWI4LWQ3YzUtMTFlYi04ZGY4LTAwMDAxNzAyMjAwOCIsCiAgImFkZCI6ImFwaS5zc2ZyZWUucnUiLAogICJ2IjoiMiIsCiAgImFpZCI6IjY0IiwKICAibmV0Ijoid3MiLAogICJ0eXBlIjoibm9uZSIKfQ=="}

	cfg, d, err := f.Fetch()
	if err != nil {
		t.Fatal(err)
	}

	config, _ := json.MarshalIndent(cfg, "", "  ")

	t.Log("timeout:", d)
	t.Log("config:", string(config))
}
