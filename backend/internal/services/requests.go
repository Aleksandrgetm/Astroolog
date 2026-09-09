package services

import (
	"astroolog/backend/internal/models"
	"astroolog/backend/internal/repositories"
	"errors"
	"gorm.io/gorm"
	"regexp"
	"strings"
	"time"
)

var ErrInvalid = errors.New("invalid request")

type Input struct {
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
	if i.Name == "" || (i.Phone != "" && !regexp.MustCompile(`^\+?[0-9 ()-]{6,30}$`).MatchString(i.Phone)) {
		return ErrInvalid
	}
	return nil
}

type Requests struct{ Repo repositories.Repository }

func (s Requests) Book(i Input) error {
	if err := i.Normalize(); err != nil {
		return err
	}
	if i.ServiceID == 0 {
		return ErrInvalid
	}
	if err := s.Repo.ActiveService(i.ServiceID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalid
		}
		return err
	}
	var date *time.Time
	if i.PreferredDate != "" {
		d, err := time.Parse("2006-01-02", i.PreferredDate)
		if err != nil || d.Before(time.Now().UTC().Truncate(24*time.Hour)) {
			return ErrInvalid
		}
		date = &d
	}
	return s.Repo.Booking(&models.Booking{Name: i.Name, Email: i.Email, Phone: i.Phone, ServiceID: i.ServiceID, PreferredDate: date, Message: i.Message, Status: "new"})
}
func (s Requests) Contact(i Input) error {
	if err := i.Normalize(); err != nil {
		return err
	}
	if i.Message == "" {
		return ErrInvalid
	}
	return s.Repo.Contact(&models.ContactRequest{Name: i.Name, Email: i.Email, Phone: i.Phone, Message: i.Message})
}
