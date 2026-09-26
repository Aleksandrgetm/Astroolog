package database

import (
	"astroolog/backend/internal/models"
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// The previous revision remains unchanged; this additive revision also runs on existing databases.
func Migrate(db *gorm.DB) error {
	if err := migratePreAdmin(db); err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(20260915)").Error; err != nil {
			return err
		}
		const revision = "booking-bonus-20260926-v1"
		var count int64
		if err := tx.Model(&models.CatalogRevision{}).Where("id = ?", revision).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return nil
		}
		if err := tx.AutoMigrate(&models.BookingOption{}, &models.Booking{}); err != nil {
			return err
		}
		if err := tx.Model(&models.BookingOption{}).Where("code IN ?", []string{"reading-call-40", "numerology-reading-call"}).Update("is_active", false).Error; err != nil {
			return err
		}
		var service models.Service
		if err := tx.Where("slug = ?", "personal-matrix").First(&service).Error; err != nil {
			return err
		}
		price := int64(2000)
		bonus := models.BookingOption{Code: "online-meeting-40", ServiceID: service.ID, IsAddon: true, IsActive: true, SortOrder: 6, Price: &price,
			Title: "Бонус + онлайн-встреча 40 минут", TitleRU: "Бонус + онлайн-встреча 40 минут", TitleLV: "Papildinājums + 40 minūšu tiešsaistes tikšanās", TitleEN: "Bonus + 40-minute online meeting",
			Format: "Онлайн", FormatRU: "Онлайн", FormatLV: "Tiešsaistē", FormatEN: "Online"}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&bonus).Error; err != nil {
			return err
		}
		// Only replace known bootstrap copy. Preserve any independently edited content.
		var original catalogData
		if err := json.Unmarshal(catalogJSON, &original); err != nil {
			return err
		}
		replacements := map[string]string{
			"description":    "Личность — характер, особенности, сильные стороны, внутренние ресурсы и то, как лучше раскрывать свой потенциал. Также можно выбрать отдельный разбор финансов, отношений или детской матрицы. По желанию к основному разбору можно добавить бонус — онлайн-встречу 40 минут, чтобы обсудить результаты, вопросы и первые практические ориентиры.",
			"description_lv": "Personība — raksturs, īpatnības, stiprās puses, iekšējie resursi un sava potenciāla atklāšana. Var izvēlēties arī finanšu, attiecību vai bērna matricas analīzi. Izvēlētajai analīzei pēc vēlēšanās var pievienot 40 minūšu tiešsaistes tikšanos, lai pārrunātu rezultātus, jautājumus un pirmos praktiskos soļus.",
			"description_en": "Personality — character, individual traits, strengths, inner resources and fulfilling your potential. You can also choose a reading of finances, relationships or a child’s matrix. Optionally add a 40-minute online meeting to your reading to discuss the results, questions and practical first steps.",
		}
		replacements["description_ru"] = replacements["description"]
		for _, old := range original.Services {
			if old.Slug == service.Slug {
				values := map[string]string{"description": old.Description, "description_ru": old.DescriptionRU, "description_lv": old.DescriptionLV, "description_en": old.DescriptionEN}
				for column, value := range replacements {
					if err := tx.Model(&models.Service{}).Where("id = ? AND "+column+" = ?", service.ID, values[column]).Update(column, value).Error; err != nil {
						return err
					}
				}
			}
		}
		return tx.Create(&models.CatalogRevision{ID: revision}).Error
	})
}
