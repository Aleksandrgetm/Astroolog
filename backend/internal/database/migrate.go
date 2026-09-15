package database

import (
	"astroolog/backend/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.User{}, &models.Service{}, &models.Testimonial{}, &models.Booking{}, &models.ContactRequest{}, &models.BookingOption{}, &models.CatalogRevision{}); err != nil {
		return err
	}
	if err := MigrateLocalizedContent(db); err != nil {
		return err
	}
	return MigrateCurrentCatalog(db)
}

// SeedDemo is opt-in: demo copy must not be published as real client feedback.
func SeedDemo(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&models.Testimonial{}).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return tx.Create(&[]models.Testimonial{
				{ClientName: "Анна · пример отзыва", Text: "После встречи стало проще услышать себя. Вместо длинного списка «надо» появился один понятный шаг, который я действительно хочу сделать.", Rating: 5, IsActive: true, SortOrder: 1},
				{ClientName: "Мария · пример отзыва", Text: "Для меня самым ценным оказалась спокойная атмосфера. Можно было говорить честно, без страха показаться неправильной.", Rating: 5, IsActive: true, SortOrder: 2},
				{ClientName: "Ольга · пример отзыва", Text: "Я пришла с вопросом о работе, а ушла с пониманием, что для меня важно. Теперь есть над чем подумать и с чего начать.", Rating: 5, IsActive: true, SortOrder: 3},
			}).Error
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err := MigrateLocalizedContent(db); err != nil {
		return err
	}
	return MigrateCurrentCatalog(db)
}
