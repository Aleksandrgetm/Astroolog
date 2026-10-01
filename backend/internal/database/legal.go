package database

import (
	"astroolog/backend/internal/models"
	_ "embed"
	"encoding/json"
	"gorm.io/gorm"
)

//go:embed legal_seed.json
var legalSeed []byte

func legalDefinitions() CMSSeedData {
	var data CMSSeedData
	if err := json.Unmarshal(legalSeed, &data); err != nil {
		panic(err)
	}
	return data
}

// Add legal sections once, preserving all existing editorial content and unrelated pages.
func migrateLegal(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(20260915)").Error; err != nil {
			return err
		}
		const revision = "legal-documents-20261001-v2"
		var count int64
		if err := tx.Model(&models.CatalogRevision{}).Where("id=?", revision).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return nil
		}
		for _, definition := range legalDefinitions().Pages {
			key := definition.Key
			var page models.Page
			if err := tx.Where("key=?", key).Attrs(models.Page{Key: key, Slug: "/" + key, IsActive: true}).FirstOrCreate(&page).Error; err != nil {
				return err
			}
			for _, tr := range definition.Translations {
				tr.PageID = page.ID
				var existing models.PageTranslation
				if err := tx.Where("page_id=? AND language=?", page.ID, tr.Language).Attrs(tr).FirstOrCreate(&existing).Error; err != nil {
					return err
				}
			}
			for i, d := range legalDefinitions().Sections {
				if d.Page != key {
					continue
				}
				section := models.PageSection{PageID: page.ID, SectionKey: d.SectionKey, SectionType: d.SectionType, SortOrder: i, IsActive: true, Settings: d.Settings, Translations: d.Translations}
				var existing int64
				if err := tx.Model(&models.PageSection{}).Where("section_key=?", d.SectionKey).Count(&existing).Error; err != nil {
					return err
				}
				if existing == 0 {
					if err := tx.Create(&section).Error; err != nil {
						return err
					}
				}
			}
		}
		return tx.Create(&models.CatalogRevision{ID: revision}).Error
	})
}
