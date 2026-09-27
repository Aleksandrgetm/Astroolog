package services

import (
	"astroolog/backend/internal/models"
	"astroolog/backend/internal/repositories"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"regexp"
	"strings"
	"time"
	_ "time/tzdata"
)

var ErrInvalid = errors.New("invalid request")
var ErrStalePrice = errors.New("stale price")
var ErrIdempotency = errors.New("idempotency key reused with different payload")
var riga, _ = time.LoadLocation("Europe/Riga")
var phonePattern = regexp.MustCompile(`^\+?[0-9 ().-]+$`)
var keyPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{16,128}$`)

func ValidPhone(phone string) bool {
	if phone == "" {
		return true
	}
	if len(phone) > 30 || !phonePattern.MatchString(phone) {
		return false
	}
	digits := 0
	for _, c := range phone {
		if c >= '0' && c <= '9' {
			digits++
		}
	}
	return digits >= 7 && digits <= 15
}
func ValidDate(value string, now time.Time) bool {
	d, err := time.ParseInLocation("2006-01-02", value, riga)
	if err != nil {
		return false
	}
	local := now.In(riga)
	today := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, riga)
	return !d.Before(today)
}
func localized(ru, lv, en, legacy, language string) string {
	v := ru
	if language == "lv" {
		v = lv
	}
	if language == "en" {
		v = en
	}
	if strings.TrimSpace(v) == "" {
		v = ru
	}
	if strings.TrimSpace(v) == "" {
		v = legacy
	}
	return v
}

type Input struct {
	SubmissionKey string `json:"-"`
	BonusCode     string `json:"bonus_code"`
	BonusPrice    *int64 `json:"bonus_price"`

	Language    string `json:"language" binding:"omitempty,oneof=ru lv en"`
	ServiceType string `json:"service_type"`
	Price       *int64 `json:"price"`

	Name          string `json:"name" binding:"required,max=100"`
	Email         string `json:"email" binding:"required,email,max=254"`
	Phone         string `json:"phone" binding:"omitempty,max=30"`
	Message       string `json:"message" binding:"max=3000"`
	Consent       bool   `json:"consent" binding:"required"`
	ServiceID     uint   `json:"service_id"`
	PreferredDate string `json:"preferred_date"`
}

func (i *Input) Normalize() error {
	i.Name = strings.TrimSpace(i.Name)
	i.Message = strings.TrimSpace(i.Message)
	i.Email = strings.ToLower(strings.TrimSpace(i.Email))
	i.Phone = strings.TrimSpace(i.Phone)
	if i.Language == "" {
		i.Language = "ru"
	}
	if i.Name == "" || !ValidPhone(i.Phone) {
		return ErrInvalid
	}
	return nil
}

type Requests struct{ Repo repositories.Repository }

func (s Requests) Book(i Input) error {
	if err := i.Normalize(); err != nil {
		return err
	}
	if i.ServiceID == 0 || i.ServiceType == "" {
		return ErrInvalid
	}
	if err := s.Repo.ActiveService(i.ServiceID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalid
		}
		return err
	}
	option, err := s.Repo.BookingOption(i.ServiceType)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrInvalid
	}
	if err != nil {
		return err
	}
	// Reject stale/tampered quotes; never silently store a price the visitor did not see.
	if option.ServiceID != i.ServiceID || option.IsAddon {
		return ErrInvalid
	}
	if (option.Price == nil) != (i.Price == nil) || (option.Price != nil && *option.Price != *i.Price) {
		return ErrStalePrice
	}
	var bonusCode, bonusTitle *string
	bonusFormat := ""
	var bonusPrice *int64
	total := option.Price
	if i.BonusCode != "" {
		bonus, err := s.Repo.BookingOption(i.BonusCode)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalid
		}
		if err != nil {
			return err
		}
		if !bonus.IsAddon || bonus.ServiceID != i.ServiceID {
			return ErrInvalid
		}
		if bonus.Price == nil || i.BonusPrice == nil || *bonus.Price != *i.BonusPrice {
			return ErrStalePrice
		}
		if option.Price == nil {
			return ErrInvalid
		}
		sum := *option.Price + *bonus.Price
		if sum < *option.Price {
			return ErrInvalid
		}
		total = &sum
		bonusCode = &bonus.Code
		title := localized(bonus.TitleRU, bonus.TitleLV, bonus.TitleEN, bonus.Title, i.Language)
		bonusTitle = &title
		bonusPrice = bonus.Price
		bonusFormat = localized(bonus.FormatRU, bonus.FormatLV, bonus.FormatEN, bonus.Format, i.Language)
	} else if i.BonusPrice != nil {
		return ErrInvalid
	}
	service, err := s.Repo.ServiceByID(i.ServiceID)
	if err != nil {
		return err
	}
	title := localized(service.TitleRU, service.TitleLV, service.TitleEN, service.Title, i.Language)
	optionTitle := localized(option.TitleRU, option.TitleLV, option.TitleEN, option.Title, i.Language)
	format := localized(option.FormatRU, option.FormatLV, option.FormatEN, option.Format, i.Language)
	if bonusFormat != "" {
		format += " / " + bonusFormat
	}
	now := time.Now()
	return s.Repo.Booking(&models.Booking{ConsentAt: &now, ServiceTitle: &title, OptionTitle: &optionTitle, OptionFormat: &format, Language: &i.Language, ServiceType: option.Code, Price: option.Price, Currency: "EUR", Name: i.Name, Email: i.Email, Phone: i.Phone, ServiceID: i.ServiceID, PreferredDate: nil, SubmissionKey: optionalKey(i.SubmissionKey), BonusCode: bonusCode, BonusTitle: bonusTitle, BonusPrice: bonusPrice, TotalPrice: total, Message: i.Message, Status: "new"})
}
func (s Requests) Contact(i Input) error {
	if err := i.Normalize(); err != nil {
		return err
	}
	if i.Message == "" {
		return ErrInvalid
	}
	now := time.Now()
	return s.Repo.Contact(&models.ContactRequest{ConsentAt: &now, Status: "new", Language: i.Language, Name: i.Name, Email: i.Email, Phone: i.Phone, Message: i.Message})
}

func (s Requests) Submit(i Input, booking bool, key string) error {
	if !keyPattern.MatchString(key) {
		return ErrInvalid
	}
	if err := i.Normalize(); err != nil {
		return err
	}
	raw, err := json.Marshal(i)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(raw)
	hash := hex.EncodeToString(sum[:])
	kind := "contact"
	if booking {
		kind = "bookings"
	}
	return s.Repo.DB.Transaction(func(tx *gorm.DB) error {
		claim := models.SubmissionKey{Key: key, Kind: kind, Fingerprint: hash}
		result := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&claim)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			var old models.SubmissionKey
			if err := tx.First(&old, "key = ?", key).Error; err != nil {
				return err
			}
			if old.Kind != kind || old.Fingerprint != hash {
				return ErrIdempotency
			}
			return nil
		}
		worker := Requests{Repo: repositories.Repository{DB: tx}}
		i.SubmissionKey = key
		if booking {
			return worker.Book(i)
		}
		return worker.Contact(i)
	})
}

func optionalKey(key string) *string {
	if key == "" {
		return nil
	}
	return &key
}
