package database

import (
	"astroolog/backend/internal/models"
	"fmt"
	"gorm.io/gorm"
)

// Versioned, transactional schema migration. Advisory lock serializes concurrent starts.
func migratePreAdmin(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(20260915)").Error; err != nil {
			return err
		}
		if err := tx.Exec("CREATE TABLE IF NOT EXISTS catalog_revisions (id text PRIMARY KEY)").Error; err != nil {
			return err
		}
		var applied int64
		if err := tx.Model(&models.CatalogRevision{}).Where("id = ?", "schema-pre-admin-v1").Count(&applied).Error; err != nil {
			return err
		}
		if applied > 0 {
			return nil
		}
		if tx.Migrator().HasTable("booking_options") {
			// These two known legacy formats were individual numerology readings.
			if err := tx.Exec(`UPDATE booking_options SET service_id=(SELECT id FROM services WHERE slug='personal-matrix') WHERE service_id IS NULL AND code IN ('numerology-reading','numerology-reading-call')`).Error; err != nil {
				return err
			}
			var invalid int64
			if err := tx.Raw(`SELECT count(*) FROM booking_options o LEFT JOIN services s ON s.id=o.service_id WHERE s.id IS NULL OR o.price<0`).Scan(&invalid).Error; err != nil {
				return err
			}
			if invalid > 0 {
				return fmt.Errorf("migration stopped: %d orphan/negative-price booking options require review", invalid)
			}
		}
		for _, table := range []string{"services", "bookings"} {
			if tx.Migrator().HasTable(table) && tx.Migrator().HasColumn(table, "price") {
				var n int64
				if err := tx.Table(table).Where("price < 0").Count(&n).Error; err != nil {
					return err
				}
				if n > 0 {
					return fmt.Errorf("migration stopped: negative prices in %s", table)
				}
			}
		}
		if err := tx.AutoMigrate(&models.User{}, &models.Service{}, &models.Testimonial{}, &models.BookingOption{}, &models.Booking{}, &models.ContactRequest{}, &models.SubmissionKey{}); err != nil {
			return err
		}
		if err := tx.Exec("UPDATE contact_requests SET status='new' WHERE status IS NULL OR status='' ").Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE contact_requests SET updated_at=created_at WHERE updated_at IS NULL").Error; err != nil {
			return err
		}
		if err := tx.Exec("UPDATE testimonials SET is_active=false WHERE client_name IN ('Анна · пример отзыва','Мария · пример отзыва','Ольга · пример отзыва')").Error; err != nil {
			return err
		}
		if err := MigrateCurrentCatalog(tx); err != nil {
			return err
		}
		return tx.Create(&models.CatalogRevision{ID: "schema-pre-admin-v1"}).Error
	})
}
