package repositories

import (
	"astroolog/backend/internal/models"
	"gorm.io/gorm"
)

type Repository struct{ DB *gorm.DB }

func (r Repository) Services() ([]models.Service, error) {
	v := []models.Service{}
	err := r.DB.Where("is_active = ?", true).Order("sort_order, id").Find(&v).Error
	return v, err
}
func (r Repository) Service(slug string) (models.Service, error) {
	var v models.Service
	err := r.DB.Where("slug = ? AND is_active = ?", slug, true).First(&v).Error
	return v, err
}
func (r Repository) ActiveService(id uint) error {
	var v models.Service
	return r.DB.Where("id = ? AND is_active = ?", id, true).First(&v).Error
}
func (r Repository) Testimonials() ([]models.Testimonial, error) {
	v := []models.Testimonial{}
	err := r.DB.Where("is_active = ?", true).Order("sort_order, id").Find(&v).Error
	return v, err
}
func (r Repository) Booking(v *models.Booking) error        { return r.DB.Create(v).Error }
func (r Repository) Contact(v *models.ContactRequest) error { return r.DB.Create(v).Error }
