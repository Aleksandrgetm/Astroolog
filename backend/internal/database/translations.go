package database

import (
	_ "embed"
	"encoding/json"
	"gorm.io/gorm"
)

//go:embed translations.json
var translationJSON []byte

type translatedValue struct {
	RU string `json:"ru"`
	LV string `json:"lv"`
	EN string `json:"en"`
}
type translatedRow struct {
	Slug       string                     `json:"slug"`
	ClientName string                     `json:"client_name"`
	Fields     map[string]translatedValue `json:"fields"`
}
type translationData struct {
	Services     []translatedRow `json:"services"`
	Testimonials []translatedRow `json:"testimonials"`
}

// Additive and idempotent: legacy data and nonempty translations are never overwritten.
func MigrateLocalizedContent(db *gorm.DB) error {
	var data translationData
	if err := json.Unmarshal(translationJSON, &data); err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		fields := map[string][]string{"services": {"title", "short_description", "description", "duration"}, "testimonials": {"text", "client_name"}}
		for table, columns := range fields {
			for _, field := range columns {
				if err := tx.Table(table).Where("COALESCE(TRIM("+field+"_ru), '') = ''").UpdateColumn(field+"_ru", gorm.Expr(field)).Error; err != nil {
					return err
				}
			}
		}
		// Seed translations only for exact known demo copy, not client-authored content.
		for table, rows := range map[string][]translatedRow{"services": data.Services, "testimonials": data.Testimonials} {
			for _, row := range rows {
				for field, values := range row.Fields {
					// Only allow the known content columns, even if the embedded seed is edited.
					allowed := false
					for _, column := range fields[table] {
						if column == field {
							allowed = true
						}
					}
					if !allowed {
						continue
					}
					for locale, value := range map[string]string{"lv": values.LV, "en": values.EN} {
						query := tx.Table(table).Where(field+"_ru = ?", values.RU).Where("COALESCE(TRIM(" + field + "_" + locale + "), '') = ''")
						if table == "services" {
							query = query.Where("slug = ?", row.Slug)
						} else {
							query = query.Where("client_name = ?", row.ClientName)
						}
						if err := query.UpdateColumn(field+"_"+locale, value).Error; err != nil {
							return err
						}
					}
				}
			}
		}
		return nil
	})
}
