package services_test

import (
	"astroolog/backend/internal/config"
	"astroolog/backend/internal/database"
	"astroolog/backend/internal/models"
	"astroolog/backend/internal/repositories"
	"astroolog/backend/internal/services"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestDatabaseIntegration(t *testing.T) {
	name := os.Getenv("TEST_DB_NAME")
	if name == "" {
		t.Skip("requires isolated TEST_DB_NAME")
	}
	if !strings.HasSuffix(name, "_test") {
		t.Fatal("isolated test database required")
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	cfg.DBName = name
	db := database.Connect(cfg)
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	repo := repositories.Repository{DB: db}
	requests := services.Requests{Repo: repo}
	key := fmt.Sprintf("integration-%d", time.Now().UnixNano())
	input := services.Input{Name: "Integration test", Email: "integration@example.invalid", Message: "Isolated database test", Consent: true, Language: "lv"}
	defer func() {
		db.Where("key LIKE ?", key+"%").Delete(&models.SubmissionKey{})
		db.Where("email = ?", input.Email).Delete(&models.ContactRequest{})
		db.Where("email = ?", input.Email).Delete(&models.Booking{})
	}()
	var wg sync.WaitGroup
	errs := make(chan error, 6)
	for n := 0; n < 6; n++ {
		wg.Add(1)
		go func() { defer wg.Done(); errs <- requests.Submit(input, false, key) }()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	var count int64
	db.Model(&models.ContactRequest{}).Where("email=?", input.Email).Count(&count)
	if count != 1 {
		t.Fatalf("duplicates: %d", count)
	}
	changed := input
	changed.Message = "different"
	if !errors.Is(requests.Submit(changed, false, key), services.ErrIdempotency) {
		t.Fatal("conflicting reuse accepted")
	}
	opts, err := repo.BookingOptions()
	if err != nil || len(opts) == 0 {
		t.Fatal(err)
	}
	input.ServiceID = opts[0].ServiceID
	input.ServiceType = opts[0].Code
	input.Price = opts[0].Price
	if err := requests.Submit(input, true, key+"-booking"); err != nil {
		t.Fatal(err)
	}
	var booking models.Booking
	db.Where("email=?", input.Email).First(&booking)
	if booking.ServiceTitle == nil || booking.OptionTitle == nil || booking.Language == nil || *booking.Language != "lv" || booking.Currency != "EUR" {
		t.Fatal("missing snapshot")
	}
	rollback := errors.New("test rollback")
	err = db.Transaction(func(tx *gorm.DB) error {
		var service models.Service
		if err := tx.First(&service, input.ServiceID).Error; err != nil {
			return err
		}
		if err := tx.Model(&service).Updates(map[string]any{"title_ru": "Editor content", "price": 999999}).Error; err != nil {
			return err
		}
		if err := database.Migrate(tx); err != nil {
			return err
		}
		var after models.Service
		tx.First(&after, service.ID)
		if after.TitleRU != "Editor content" {
			t.Error("restart overwrites editor content")
		}
		derived, e := (repositories.Repository{DB: tx}).Service(service.Slug)
		if e != nil {
			return e
		}
		if derived.Price == nil || *derived.Price == 999999 {
			t.Error("service price not derived")
		}
		tx.Model(&service).Update("is_active", false)
		active, e := (repositories.Repository{DB: tx}).BookingOptions()
		if e != nil {
			return e
		}
		for _, o := range active {
			if o.ServiceID == service.ID {
				t.Error("inactive parent option returned")
			}
		}
		return rollback
	})
	if !errors.Is(err, rollback) {
		t.Fatal(err)
	}
	if input.Price != nil {
		wrong := *input.Price + 1
		input.Price = &wrong
		if !errors.Is(requests.Submit(input, true, key+"-stale"), services.ErrStalePrice) {
			t.Fatal("stale price accepted")
		}
	}
}
