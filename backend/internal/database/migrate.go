package database

import (
	"astroolog/backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.User{}, &models.Service{}, &models.Testimonial{}, &models.Booking{}, &models.ContactRequest{}); err != nil {
		return err
	}
	return MigrateLocalizedContent(db)
}

// SeedDemo is opt-in: demo copy must not be published as real client feedback.
func SeedDemo(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		services := []models.Service{
			{Title: "Познакомься с собой", Slug: "first-step", ShortDescription: "Бережное знакомство с собой и ответы на первые вопросы.", Description: "Вводная встреча для тех, кто хочет познакомиться с подходом. Обсудим ваш запрос, обратим внимание на сильные стороны и определим, с чего начать. После встречи у вас будет краткий конспект и один небольшой шаг для самостоятельной работы.", Image: "/images/journal-placeholder.webp", Duration: "30 минут · онлайн", IsActive: true, SortOrder: 1},
			{Title: "Моя матрица возможностей", Slug: "personal-matrix", ShortDescription: "Ваши сильные стороны, внутренние опоры и возможности для роста.", Description: "Индивидуальная встреча о ваших качествах, привычных сценариях и желаниях. Используем нумерологический разбор как повод для размышления, а не как предсказание. Вы получите запись встречи по согласованию и письменные вопросы для самостоятельной работы.", Image: "/images/journal-placeholder.webp", Duration: "90 минут · онлайн", IsActive: true, SortOrder: 2},
			{Title: "Точка роста", Slug: "growth-point", ShortDescription: "Один важный запрос. Больше ясности. Конкретный следующий шаг.", Description: "Сфокусированная консультация о работе, деньгах, отношениях или родительстве. Разберём текущую ситуацию, отделим ваши желания от ожиданий окружающих и составим реалистичный план действий. В детских запросах обсуждаем поддержку ребёнка без ярлыков и прогнозов.", Image: "/images/relationships-placeholder.webp", Duration: "60 минут · онлайн", IsActive: true, SortOrder: 3},
			{Title: "Путь к себе", Slug: "path-to-self", ShortDescription: "Личное сопровождение, чтобы перенести понимание в повседневную жизнь.", Description: "Четыре индивидуальные встречи в течение месяца. На первой определим запрос и критерии прогресса. Между встречами — небольшие практики и наблюдения. Вместе будем уточнять план и поддерживать изменения в комфортном для вас темпе.", Image: "/images/expert-placeholder.webp", Duration: "4 встречи · 1 месяц", IsActive: true, SortOrder: 4},
		}
		if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "slug"}}, DoNothing: true}).Create(&services).Error; err != nil {
			return err
		}
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
	return MigrateLocalizedContent(db)
}
