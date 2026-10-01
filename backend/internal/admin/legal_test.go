package admin

import (
	"astroolog/backend/internal/models"
	"testing"
)

func TestTechnicalCookieContentCannotBeEdited(t *testing.T) {
	for _, key := range []string{"legal.cookies.section2", "legal.cookies.section3", "legal.cookies.section4", "legal.cookies.section5", "legal.cookies.section8"} {
		section := &models.PageSection{SectionKey: key, IsActive: true, Translations: []models.PageSectionTranslation{{Language: "ru", Fields: models.Document{}}}}
		if err := (Handler{}).validate(nil, "sections", section, map[string]any{}); err == nil {
			t.Fatalf("technical section editable: %s", key)
		}
	}
	section := &models.PageSection{SectionKey: "legal.privacy.section1", IsActive: true, Translations: []models.PageSectionTranslation{{Language: "ru", Fields: models.Document{"legal.privacy.section1.title": "Заголовок", "legal.privacy.section1.body": "Обычный текст"}}}}
	if err := (Handler{}).validate(nil, "sections", section, map[string]any{}); err != nil {
		t.Fatal(err)
	}
}
