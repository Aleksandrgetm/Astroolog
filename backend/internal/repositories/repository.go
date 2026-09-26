package repositories

import (
	"astroolog/backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct{ DB *gorm.DB }

func (r Repository) Services() ([]models.Service, error) {
	v := []models.Service{}
	err := r.DB.Model(&models.Service{}).Select("services.*, (SELECT MIN(o.price) FROM booking_options o WHERE o.service_id=services.id AND o.is_active=true AND o.is_addon=false) AS price").Where("services.is_active = ?", true).Order("sort_order, id").Find(&v).Error
	return v, err
}
func (r Repository) Service(slug string) (models.Service, error) {
	var v models.Service
	err := r.DB.Model(&models.Service{}).Select("services.*, (SELECT MIN(o.price) FROM booking_options o WHERE o.service_id=services.id AND o.is_active=true AND o.is_addon=false) AS price").Where("slug = ? AND is_active = ?", slug, true).First(&v).Error
	return v, err
}
func (r Repository) ActiveService(id uint) error {
	var v models.Service
	return r.DB.Clauses(clause.Locking{Strength: "SHARE"}).Where("id = ? AND is_active = ?", id, true).First(&v).Error
}
func (r Repository) Booking(v *models.Booking) error        { return r.DB.Create(v).Error }
func (r Repository) Contact(v *models.ContactRequest) error { return r.DB.Create(v).Error }

func (r Repository) BookingOptions() ([]models.BookingOption, error) {
	options := []models.BookingOption{}
	err := r.DB.Model(&models.BookingOption{}).Select("booking_options.*").Joins("JOIN services ON services.id=booking_options.service_id").Where("booking_options.is_active=true AND services.is_active=true").Order("booking_options.sort_order, code").Find(&options).Error
	return options, err
}
func (r Repository) BookingOption(code string) (models.BookingOption, error) {
	var option models.BookingOption
	err := r.DB.Clauses(clause.Locking{Strength: "SHARE"}).Where("code = ? AND is_active = ?", code, true).First(&option).Error
	return option, err
}

func (r Repository) ServiceByID(id uint) (models.Service, error) {
	var s models.Service
	err := r.DB.First(&s, id).Error
	return s, err
}
