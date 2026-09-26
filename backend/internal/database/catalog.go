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

// Bootstrap only an empty catalog. Restarts never update editorial content.
func MigrateCurrentCatalog(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.Service{}).Count(&count).Error; err != nil {
		return err
	}
	if count != 0 {
		return nil
	}
	var data catalogData
	if err := json.Unmarshal(catalogJSON, &data); err != nil {
		return err
	}
	ids := map[string]uint{}
	for _, service := range data.Services {
		if err := db.Create(&service).Error; err != nil {
			return err
		}
		ids[service.Slug] = service.ID
	}
	for _, item := range data.Options {
		item.ServiceID = ids[item.ServiceSlug]
		if err := db.Create(&item.BookingOption).Error; err != nil {
			return err
		}
	}
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&models.CatalogRevision{ID: "catalog-bootstrap-v2"}).Error
}
