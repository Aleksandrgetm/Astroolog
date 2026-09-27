package database

import (
	"astroolog/backend/internal/models"
	_ "embed"
	"encoding/json"
	"gorm.io/gorm"
)

//go:embed cms_seed.json
var cmsSeed []byte

type SectionDefinition struct {
	SectionKey  string          `json:"section_key"`
	Page        string          `json:"page"`
	SectionType string          `json:"section_type"`
	Label       string          `json:"label"`
	Critical    bool            `json:"critical"`
	Settings    models.Document `json:"settings"`
	Fields      []struct {
		Key   string `json:"key"`
		Label string `json:"label"`
	} `json:"fields"`
	Translations []models.PageSectionTranslation `json:"translations"`
}
type CMSSeedData struct {
	Pages      []models.Page       `json:"pages"`
	Sections   []SectionDefinition `json:"sections"`
	Directions []struct {
		models.Direction
		ServiceSlug string `json:"service_slug"`
	} `json:"directions"`
	Reviews      []models.Review `json:"reviews"`
	Certificates []struct {
		models.Certificate
		Images []struct {
			Language string `json:"language"`
			Path     string `json:"path"`
		} `json:"images"`
	} `json:"certificates"`
	Steps []models.ProcessStep `json:"steps"`
}

func CMSDefinitions() CMSSeedData {
	var data CMSSeedData
	if err := json.Unmarshal(cmsSeed, &data); err != nil {
		panic(err)
	}
	return data
}
func migrateCMS(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if e := tx.Exec("SELECT pg_advisory_xact_lock(20260915)").Error; e != nil {
			return e
		}
		const revision = "admin-cms-20260926-v1"
		var n int64
		if e := tx.Model(&models.CatalogRevision{}).Where("id=?", revision).Count(&n).Error; e != nil {
			return e
		}
		if n > 0 {
			return nil
		}
		if e := tx.AutoMigrate(&models.User{}, &models.Service{}, &models.BookingOption{}, &models.Booking{}, &models.ContactRequest{}, &models.AdminSession{}, &models.LoginAttempt{}, &models.AdminAuditLog{}, &models.Page{}, &models.PageTranslation{}, &models.PageSection{}, &models.PageSectionTranslation{}, &models.ProcessStep{}, &models.ProcessStepTranslation{}, &models.Direction{}, &models.DirectionOption{}, &models.DirectionTranslation{}, &models.Review{}, &models.ReviewTranslation{}, &models.MediaAsset{}, &models.Certificate{}, &models.CertificateTranslation{}, &models.CertificateAsset{}, &models.SiteSettings{}, &models.Publication{}); e != nil {
			return e
		}
		// Preserve old operational status separately; migrate only the new admin workflow.
		if e := tx.Exec("UPDATE bookings SET workflow_status=CASE WHEN status='completed' THEN 'completed' WHEN status='cancelled' THEN 'archived' WHEN status='confirmed' THEN 'in_progress' ELSE 'new' END").Error; e != nil {
			return e
		}
		data := CMSDefinitions()
		pages := map[string]uint{}
		for _, p := range data.Pages {
			if e := tx.Create(&p).Error; e != nil {
				return e
			}
			pages[p.Key] = p.ID
		}
		for i, d := range data.Sections {
			s := models.PageSection{PageID: pages[d.Page], SectionKey: d.SectionKey, SectionType: d.SectionType, SortOrder: i, IsActive: true, Settings: d.Settings, Translations: d.Translations}
			if order, ok := map[string]int{"home.hero": 0, "home.supporting": 1, "home.pain": 2, "home.numerology": 3, "home.directions": 4, "home.process": 5, "home.services": 6, "home.about-preview": 7, "home.testimonials": 8, "global.cta": 9}[s.SectionKey]; ok {
				s.SortOrder = order
			}
			if e := tx.Create(&s).Error; e != nil {
				return e
			}
			if d.SectionKey == "home.process" {
				for _, step := range data.Steps {
					step.PageSectionID = s.ID
					if e := tx.Create(&step).Error; e != nil {
						return e
					}
				}
			}
		}
		for _, d := range data.Directions {
			var service models.Service
			if e := tx.First(&service, "slug=?", d.ServiceSlug).Error; e != nil {
				return e
			}
			d.ServiceID = service.ID
			if e := tx.Create(&d.Direction).Error; e != nil {
				return e
			}
		}
		for _, r := range data.Reviews {
			if e := tx.Create(&r).Error; e != nil {
				return e
			}
		}
		for _, c := range data.Certificates {
			for _, img := range c.Images {
				m := models.MediaAsset{Filename: img.Path, OriginalFilename: img.Path, MimeType: "image/webp", Width: 1600, Height: 1131, StoragePath: img.Path}
				if e := tx.Create(&m).Error; e != nil {
					return e
				}
				c.Assets = append(c.Assets, models.CertificateAsset{Language: img.Language, MediaID: m.ID})
			}
			if e := tx.Create(&c.Certificate).Error; e != nil {
				return e
			}
		}
		for _, path := range []string{"/images/optimized/expert-1672.webp", "/images/journal-placeholder.webp", "/images/relationships-placeholder.webp", "/images/child-placeholder.webp", "/images/Hero-new.png"} {
			media := models.MediaAsset{StoragePath: path, Filename: path, OriginalFilename: path, MimeType: "image/webp"}
			if e := tx.Where("storage_path=?", path).FirstOrCreate(&media).Error; e != nil {
				return e
			}
		}
		settings := models.SiteSettings{ID: 1, ExpertName: "Елена Захарова", Phone: "+371 29 580 232", Email: "jelenabobrovska@gmail.com", WhatsappPhone: "37129580232", TelegramURL: "https://t.me/AljonaZ", ProfessionalTitleRU: "НУМЕРОЛОГ · ЖЕНСКИЙ КОУЧ", ProfessionalTitleLV: "NUMEROLOĢE · SIEVIEŠU KOUČS", ProfessionalTitleEN: "NUMEROLOGIST · WOMEN’S COACH"}
		if e := tx.Create(&settings).Error; e != nil {
			return e
		}
		return tx.Create(&models.CatalogRevision{ID: revision}).Error
	})
}

func Migrate(db *gorm.DB) error {
	if e := migrateBonus(db); e != nil {
		return e
	}
	if e := migrateCMS(db); e != nil {
		return e
	}
	return migrateCMSMedia(db)
}

//go:embed cms_media_seed.json
var cmsMediaSeed []byte

// Initial asset metadata is versioned separately so an existing CMS installation
// gains accurate dimensions and responsive variants without touching edited content.
func migrateCMSMedia(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if e := tx.Exec("SELECT pg_advisory_xact_lock(20260915)").Error; e != nil {
			return e
		}
		const revision = "admin-media-metadata-20260926-v1"
		var n int64
		if e := tx.Model(&models.CatalogRevision{}).Where("id=?", revision).Count(&n).Error; e != nil {
			return e
		}
		if n > 0 {
			return nil
		}
		var assets []models.MediaAsset
		if e := json.Unmarshal(cmsMediaSeed, &assets); e != nil {
			return e
		}
		for _, asset := range assets {
			var existing models.MediaAsset
			if e := tx.Where("storage_path=?", asset.StoragePath).FirstOrCreate(&existing, models.MediaAsset{StoragePath: asset.StoragePath}).Error; e != nil {
				return e
			}
			if e := tx.Model(&existing).Updates(map[string]any{"filename": asset.Filename, "original_filename": asset.OriginalFilename, "mime_type": asset.MimeType, "size": asset.Size, "width": asset.Width, "height": asset.Height, "variants": asset.Variants}).Error; e != nil {
				return e
			}
		}
		return tx.Create(&models.CatalogRevision{ID: revision}).Error
	})
}
