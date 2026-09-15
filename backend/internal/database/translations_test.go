package database

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDemoTranslationsAreCompleteAndPreserveNames(t *testing.T) {
	var data translationData
	if err := json.Unmarshal(translationJSON, &data); err != nil {
		t.Fatal(err)
	}
	if len(data.Services) != 3 || len(data.Testimonials) != 3 {
		t.Fatal("incomplete demo seed")
	}
	for _, service := range data.Services {
		for _, field := range []string{"title", "short_description", "description", "duration"} {
			value, ok := service.Fields[field]
			if !ok || strings.TrimSpace(value.RU) == "" || strings.TrimSpace(value.LV) == "" || strings.TrimSpace(value.EN) == "" {
				t.Fatalf("missing %s for %s", field, service.Slug)
			}
		}
	}
	for _, review := range data.Testimonials {
		name := strings.Split(review.ClientName, " · ")[0]
		for _, field := range []string{"client_name", "text"} {
			value := review.Fields[field]
			if value.RU == "" || value.LV == "" || value.EN == "" {
				t.Fatalf("missing testimonial %s", field)
			}
		}
		for _, value := range []string{review.Fields["client_name"].RU, review.Fields["client_name"].LV, review.Fields["client_name"].EN} {
			if !strings.HasPrefix(value, name+" · ") {
				t.Fatalf("proper name changed: %s", value)
			}
		}
	}
}
