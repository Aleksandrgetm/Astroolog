package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Document map[string]interface{}

func (d Document) Value() (driver.Value, error) {
	if d == nil {
		return "{}", nil
	}
	b, e := json.Marshal(d)
	return string(b), e
}
func (d *Document) Scan(v interface{}) error {
	var b []byte
	switch x := v.(type) {
	case []byte:
		b = x
	case string:
		b = []byte(x)
	case nil:
		*d = Document{}
		return nil
	default:
		return fmt.Errorf("invalid JSONB")
	}
	return json.Unmarshal(b, d)
}

type AdminSession struct {
	ID        string    `gorm:"primaryKey" json:"-"`
	UserID    uint      `gorm:"index;not null"`
	CSRFHash  string    `json:"-"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time
}
type LoginAttempt struct {
	Key         string `gorm:"primaryKey"`
	Count       int
	WindowStart time.Time
}
type AdminAuditLog struct {
	Base
	UserID     uint     `json:"user_id"`
	Action     string   `json:"action"`
	EntityType string   `json:"entity_type"`
	EntityID   string   `json:"entity_id"`
	Metadata   Document `json:"metadata" gorm:"type:jsonb"`
}
type Page struct {
	Base
	Key          string            `json:"key" gorm:"uniqueIndex"`
	Slug         string            `json:"slug"`
	IsActive     bool              `json:"is_active"`
	Translations []PageTranslation `json:"translations"`
}
type PageTranslation struct {
	ID              uint   `json:"id" gorm:"primaryKey"`
	PageID          uint   `json:"page_id" gorm:"uniqueIndex:page_lang"`
	Language        string `json:"language" gorm:"uniqueIndex:page_lang"`
	Title           string `json:"title"`
	MetaTitle       string `json:"meta_title"`
	MetaDescription string `json:"meta_description"`
}
type PageSection struct {
	Base
	PageID       uint                     `json:"page_id" gorm:"index"`
	SectionKey   string                   `json:"section_key" gorm:"uniqueIndex"`
	SectionType  string                   `json:"section_type"`
	SortOrder    int                      `json:"sort_order"`
	IsActive     bool                     `json:"is_active"`
	Settings     Document                 `json:"settings" gorm:"type:jsonb"`
	Translations []PageSectionTranslation `json:"translations"`
}
type PageSectionTranslation struct {
	ID            uint     `json:"id" gorm:"primaryKey"`
	PageSectionID uint     `json:"page_section_id" gorm:"uniqueIndex:section_lang"`
	Language      string   `json:"language" gorm:"uniqueIndex:section_lang"`
	Fields        Document `json:"fields" gorm:"type:jsonb"`
}
type ProcessStep struct {
	Base
	PageSectionID uint                     `json:"page_section_id"`
	SortOrder     int                      `json:"sort_order"`
	IconKey       string                   `json:"icon_key"`
	Translations  []ProcessStepTranslation `json:"translations"`
}
type ProcessStepTranslation struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	ProcessStepID uint   `json:"process_step_id" gorm:"uniqueIndex:step_lang"`
	Language      string `json:"language" gorm:"uniqueIndex:step_lang"`
	Label         string `json:"label"`
	Title         string `json:"title"`
	Description   string `json:"description"`
}
type Direction struct {
	Base
	Slug         string                 `json:"slug" gorm:"uniqueIndex"`
	ContentKey   string                 `json:"content_key"`
	Image        string                 `json:"image"`
	IsActive     bool                   `json:"is_active"`
	SortOrder    int                    `json:"sort_order"`
	ServiceID    uint                   `json:"service_id"`
	Service      Service                `json:"-" gorm:"constraint:OnDelete:RESTRICT"`
	Options      []DirectionOption      `json:"options"`
	Translations []DirectionTranslation `json:"translations"`
}
type DirectionOption struct {
	DirectionID uint   `json:"direction_id" gorm:"primaryKey"`
	Code        string `json:"code" gorm:"primaryKey"`
}
type DirectionTranslation struct {
	ID               uint   `json:"id" gorm:"primaryKey"`
	DirectionID      uint   `json:"direction_id" gorm:"uniqueIndex:direction_lang"`
	Language         string `json:"language" gorm:"uniqueIndex:direction_lang"`
	Eyebrow          string `json:"eyebrow"`
	Title            string `json:"title"`
	ShortDescription string `json:"short_description"`
	Lead             string `json:"lead"`
	Context          string `json:"context"`
	FormatNote       string `json:"format_note"`
	Insights         string `json:"insights"`
	Questions        string `json:"questions"`
	Process          string `json:"process"`
	SEOTitle         string `json:"seo_title"`
	SEODescription   string `json:"seo_description"`
	Alt              string `json:"alt"`
}
type Review struct {
	Base
	IsActive           bool                `json:"is_active"`
	IsFeatured         bool                `json:"is_featured"`
	SortOrder          int                 `json:"sort_order"`
	HomeSortOrder      *int                `json:"home_sort_order"`
	AuthorDisplayName  string              `json:"author_display_name"`
	PublicationConsent bool                `json:"publication_consent"`
	Translations       []ReviewTranslation `json:"translations"`
}
type ReviewTranslation struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	ReviewID  uint   `json:"review_id" gorm:"uniqueIndex:review_lang"`
	Language  string `json:"language" gorm:"uniqueIndex:review_lang"`
	ShortText string `json:"short_text"`
	FullText  string `json:"full_text"`
}
type Certificate struct {
	Base
	Issuer       string                   `json:"issuer"`
	IssuedAt     *time.Time               `json:"issued_at"`
	Year         string                   `json:"year"`
	IsActive     bool                     `json:"is_active"`
	SortOrder    int                      `json:"sort_order"`
	Translations []CertificateTranslation `json:"translations"`
	Assets       []CertificateAsset       `json:"assets"`
}
type CertificateTranslation struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	CertificateID uint   `json:"certificate_id" gorm:"uniqueIndex:certificate_lang"`
	Language      string `json:"language" gorm:"uniqueIndex:certificate_lang"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Alt           string `json:"alt"`
}
type CertificateAsset struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	CertificateID uint       `json:"certificate_id" gorm:"uniqueIndex:certificate_asset_lang"`
	Language      string     `json:"language" gorm:"uniqueIndex:certificate_asset_lang"`
	MediaID       uint       `json:"media_id"`
	Media         MediaAsset `json:"media" gorm:"constraint:OnDelete:RESTRICT"`
}
type MediaAsset struct {
	Base
	Filename         string   `json:"filename"`
	OriginalFilename string   `json:"original_filename"`
	MimeType         string   `json:"mime_type"`
	Size             int64    `json:"size"`
	Width            int      `json:"width"`
	Height           int      `json:"height"`
	StoragePath      string   `json:"storage_path" gorm:"uniqueIndex"`
	CreatedBy        uint     `json:"created_by"`
	Variants         Document `json:"variants" gorm:"type:jsonb"`
}
type SiteSettings struct {
	ID                  uint      `json:"id" gorm:"primaryKey"`
	ExpertName          string    `json:"expert_name"`
	Phone               string    `json:"phone"`
	Email               string    `json:"email"`
	WhatsappPhone       string    `json:"whatsapp_phone"`
	TelegramURL         string    `json:"telegram_url"`
	ProfessionalTitleRU string    `json:"professional_title_ru"`
	ProfessionalTitleLV string    `json:"professional_title_lv"`
	ProfessionalTitleEN string    `json:"professional_title_en"`
	UpdatedAt           time.Time `json:"updated_at"`
}
type Publication struct {
	Base
	UserID     uint       `json:"user_id"`
	Status     string     `json:"status"`
	Error      string     `json:"error"`
	FinishedAt *time.Time `json:"finished_at"`
}
