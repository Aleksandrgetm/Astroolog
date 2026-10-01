package database

import (
	"astroolog/backend/internal/config"
	"astroolog/backend/internal/models"
	"os"
	"strings"
	"testing"
)

func TestLegalDefinitions(t *testing.T) {
	data := legalDefinitions()
	counts := map[string]int{}
	for _, s := range data.Sections {
		counts[s.Page]++
		if !s.Critical || len(s.Translations) != 3 {
			t.Fatalf("incomplete legal section %s", s.SectionKey)
		}
		for _, tr := range s.Translations {
			for _, f := range s.Fields {
				text, ok := tr.Fields[f.Key].(string)
				if !ok || strings.TrimSpace(text) == "" {
					t.Fatalf("missing field %s %s", f.Key, tr.Language)
				}
				if strings.Contains(text, "[") || strings.Contains(text, "<script") {
					t.Fatalf("unresolved content %s", f.Key)
				}
			}
		}
	}
	for key, want := range map[string]int{"privacy": 14, "terms": 14, "cookies": 11} {
		if counts[key] != want {
			t.Fatalf("%s: got %d", key, counts[key])
		}
	}
}

func TestLegalMigrationPreservesEdits(t *testing.T) {
	name := os.Getenv("TEST_DB_NAME")
	if name == "" {
		t.Skip("isolated TEST_DB_NAME required")
	}
	if !strings.HasSuffix(name, "_test") {
		t.Fatal("isolated database only")
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	cfg.DBName = name
	db := Connect(cfg)
	if err = Migrate(db); err != nil {
		t.Fatal(err)
	}
	var section models.PageSection
	if err = db.Preload("Translations").First(&section, "section_key=?", "legal.privacy.section1").Error; err != nil {
		t.Fatal(err)
	}
	tr := section.Translations[0]
	original := tr.Fields
	edited := models.Document{"legal.privacy.section1.title": "Preserved editorial title", "legal.privacy.section1.body": "Preserved editorial text"}
	if err = db.Model(&tr).Update("fields", edited).Error; err != nil {
		t.Fatal(err)
	}
	defer db.Model(&tr).Update("fields", original)
	if err = Migrate(db); err != nil {
		t.Fatal(err)
	}
	var fresh models.PageSectionTranslation
	if err = db.First(&fresh, tr.ID).Error; err != nil {
		t.Fatal(err)
	}
	if fresh.Fields["legal.privacy.section1.body"] != "Preserved editorial text" {
		t.Fatal("migration overwrote edited legal content")
	}
}
