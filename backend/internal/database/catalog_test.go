package database

import (
	"encoding/json"
	"testing"
)

func TestCurrentCatalog(t *testing.T) {
	var data catalogData
	if err := json.Unmarshal(catalogJSON, &data); err != nil {
		t.Fatal(err)
	}
	if len(data.Services) != 4 || len(data.Options) != 12 {
		t.Fatal("unexpected catalog size")
	}
	expected := map[string]int64{"introduction": 0, "personality": 5000, "finances": 5000, "relationships": 5000, "child-matrix": 5000, "reading-call-40": 7000, "full-reading": 25000, "full-reading-calls": 35000, "coaching-40": 5000, "support-1m": 15000, "support-3m": 35000, "support-6m": 60000}
	groups := map[string]bool{}
	for _, s := range data.Services {
		groups[s.Slug] = true
		if s.TitleRU == "" || s.TitleLV == "" || s.TitleEN == "" || s.DescriptionRU == "" || s.DescriptionLV == "" || s.DescriptionEN == "" {
			t.Fatal("missing service translation")
		}
	}
	for _, o := range data.Options {
		price, ok := expected[o.Code]
		if !ok || o.Price == nil || *o.Price != price || !groups[o.ServiceSlug] {
			t.Fatalf("invalid option %s", o.Code)
		}
		delete(expected, o.Code)
		if o.TitleRU == "" || o.TitleLV == "" || o.TitleEN == "" || o.FormatRU == "" || o.FormatLV == "" || o.FormatEN == "" {
			t.Fatal("missing option translation")
		}
	}
	if len(expected) != 0 {
		t.Fatal("missing options")
	}
}
