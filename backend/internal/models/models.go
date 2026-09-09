package models

import "time"

type Base struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type User struct {
	Base
	Email        string `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string `json:"-" gorm:"not null"`
	Role         string `json:"role" gorm:"not null;default:client;check:role IN ('client','admin')"`
}
type Service struct {
	TitleRU            string `json:"title_ru"`
	TitleLV            string `json:"title_lv"`
	TitleEN            string `json:"title_en"`
	ShortDescriptionRU string `json:"short_description_ru"`
	ShortDescriptionLV string `json:"short_description_lv"`
	ShortDescriptionEN string `json:"short_description_en"`
	DescriptionRU      string `json:"description_ru"`
	DescriptionLV      string `json:"description_lv"`
	DescriptionEN      string `json:"description_en"`
	DurationRU         string `json:"duration_ru"`
	DurationLV         string `json:"duration_lv"`
	DurationEN         string `json:"duration_en"`
	Base
	Title            string `json:"title"`
	Slug             string `json:"slug" gorm:"uniqueIndex;not null"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	Image            string `json:"image"`
	Price            *int64 `json:"price"` // Minor currency units; nil until prices are confirmed.
	Duration         string `json:"duration"`
	IsActive         bool   `json:"is_active" gorm:"index"`
	SortOrder        int    `json:"sort_order"`
}
type Testimonial struct {
	TextRU       string `json:"text_ru"`
	TextLV       string `json:"text_lv"`
	TextEN       string `json:"text_en"`
	ClientNameRU string `json:"client_name_ru"`
	ClientNameLV string `json:"client_name_lv"`
	ClientNameEN string `json:"client_name_en"`
	Base
	ClientName string `json:"client_name"`
	Text       string `json:"text"`
	Rating     int    `json:"rating" gorm:"check:rating >= 1 AND rating <= 5"`
	IsActive   bool   `json:"is_active" gorm:"index"`
	SortOrder  int    `json:"sort_order"`
}
type Booking struct {
	Base
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	Phone         string     `json:"phone"`
	ServiceID     uint       `json:"service_id" gorm:"not null;index"`
	Service       Service    `json:"-" gorm:"constraint:OnDelete:RESTRICT"`
	PreferredDate *time.Time `json:"preferred_date" gorm:"type:date"`
	Message       string     `json:"message"`
	Status        string     `json:"status" gorm:"default:new;check:status IN ('new','confirmed','completed','cancelled')"`
}
type ContactRequest struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
