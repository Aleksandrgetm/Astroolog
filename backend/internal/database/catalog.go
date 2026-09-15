package database

import (
	"astroolog/backend/internal/models"
	_ "embed"
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:embed catalog_20260913.json
var catalogJSON []byte

type catalogData struct {
	Services []models.Service `json:"services"`
	Options  []struct {
		models.BookingOption
		ServiceSlug string `json:"service_slug"`
	} `json:"options"`
}

// One-time catalog revision. Historical service IDs and booking price snapshots stay intact.
func MigrateCurrentCatalog(db *gorm.DB) error {
	var data catalogData
	if err := json.Unmarshal(catalogJSON, &data); err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		revision := models.CatalogRevision{ID: "catalog-2026-09-13"}
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&revision)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return migrateCatalogDescriptions(tx, data)
		}
		if err := tx.Model(&models.Service{}).Where("slug = ?", "growth-point").Update("is_active", false).Error; err != nil {
			return err
		}
		ids := map[string]uint{}
		for _, service := range data.Services {
			if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "slug"}}, DoUpdates: clause.AssignmentColumns([]string{"title", "title_ru", "title_lv", "title_en", "short_description", "short_description_ru", "short_description_lv", "short_description_en", "description", "description_ru", "description_lv", "description_en", "duration", "duration_ru", "duration_lv", "duration_en", "image", "price", "is_active", "sort_order"})}).Create(&service).Error; err != nil {
				return err
			}
			ids[service.Slug] = service.ID
		}
		if err := tx.Model(&models.BookingOption{}).Where("code IN ?", []string{"numerology-reading", "numerology-reading-call"}).Update("is_active", false).Error; err != nil {
			return err
		}
		for _, item := range data.Options {
			item.BookingOption.ServiceID = ids[item.ServiceSlug]
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&item.BookingOption).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Apply approved copy independently without changing configured prices or historical bookings.
func migrateCatalogDescriptions(tx *gorm.DB, data catalogData) error {
	result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.CatalogRevision{ID: "catalog-copy-2026-09-13"})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil
	}
	for _, service := range data.Services {
		if err := tx.Model(&models.Service{}).Where("slug = ?", service.Slug).Updates(map[string]interface{}{
			"description": service.Description, "description_ru": service.DescriptionRU, "description_lv": service.DescriptionLV, "description_en": service.DescriptionEN,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}
