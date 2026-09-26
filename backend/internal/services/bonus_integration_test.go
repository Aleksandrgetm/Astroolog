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
	"testing"
	"time"
)

func TestBonusBookingIntegration(t *testing.T) {
	name := os.Getenv("TEST_DB_NAME")
	if name == "" {
		t.Skip("requires isolated TEST_DB_NAME")
	}
	if !strings.HasSuffix(name, "_test") {
		t.Fatal("isolated DB required")
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
	rollback := errors.New("rollback test changes")
	err = db.Transaction(func(tx *gorm.DB) error {
		repo := repositories.Repository{DB: tx}
		worker := services.Requests{Repo: repo}
		base, e := repo.BookingOption("personality")
		if e != nil {
			return e
		}
		bonus, e := repo.BookingOption("online-meeting-40")
		if e != nil {
			return e
		}
		if !bonus.IsAddon || bonus.Price == nil || *bonus.Price != 2000 {
			t.Fatal("bonus not seeded")
		}
		var legacy models.BookingOption
		if e := tx.First(&legacy, "code = ?", "reading-call-40").Error; e != nil {
			return e
		}
		if legacy.IsActive || legacy.Price == nil || *legacy.Price != 7000 {
			t.Fatal("historical variant lost")
		}
		input := services.Input{Name: "Bonus integration", Email: "bonus@example.invalid", Consent: true, ServiceID: base.ServiceID, ServiceType: base.Code, Price: base.Price, Language: "ru", PreferredDate: "1900-01-01"}
		key := fmt.Sprintf("bonus-%d", time.Now().UnixNano())
		for n, code := range []string{"personality", "finances", "relationships", "child-matrix"} {
			input.ServiceType = code
			for _, withBonus := range []bool{false, true} {
				input.BonusCode = ""
				input.BonusPrice = nil
				expected := int64(5000)
				if withBonus {
					input.BonusCode = bonus.Code
					input.BonusPrice = bonus.Price
					expected = 7000
				}
				k := fmt.Sprintf("%s-%d-%t", key, n, withBonus)
				if e := worker.Submit(input, true, k); e != nil {
					return e
				}
				if e := worker.Submit(input, true, k); e != nil {
					return e
				}
				var rows []models.Booking
				tx.Where("submission_key=?", k).Find(&rows)
				if len(rows) != 1 || rows[0].TotalPrice == nil || *rows[0].TotalPrice != expected || rows[0].PreferredDate != nil {
					t.Fatalf("invalid snapshot: %+v", rows)
				}
				if withBonus && (rows[0].BonusTitle == nil || rows[0].BonusPrice == nil || *rows[0].BonusPrice != 2000) {
					t.Fatal("missing bonus snapshot")
				}
			}
		}
		check := func(label string, want error) {
			t.Helper()
			if e := worker.Submit(input, true, key+"-"+label); !errors.Is(e, want) {
				t.Fatalf("%s: %v", label, e)
			}
		}
		input.ServiceType = ""
		check("no-base", services.ErrInvalid)
		input.ServiceType = bonus.Code
		input.Price = bonus.Price
		check("addon-as-base", services.ErrInvalid)
		input.ServiceType = base.Code
		input.Price = base.Price
		tx.Model(&models.BookingOption{}).Where("code=?", bonus.Code).Update("is_active", false)
		check("inactive", services.ErrInvalid)
		tx.Model(&models.BookingOption{}).Where("code=?", bonus.Code).Update("is_active", true)
		wrong := int64(1)
		input.BonusPrice = &wrong
		check("tampered", services.ErrStalePrice)
		input.BonusPrice = bonus.Price
		input.Price = &wrong
		check("stale-base", services.ErrStalePrice)
		input.Price = base.Price
		input.ServiceType = "introduction"
		var intro models.BookingOption
		tx.First(&intro, "code=?", "introduction")
		input.ServiceID = intro.ServiceID
		input.Price = intro.Price
		check("wrong-service", services.ErrInvalid)
		historicalTitle := legacy.TitleRU
		historicalDate := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		historical := models.Booking{Name: "Historical fixture", Email: "history@example.invalid", ServiceID: legacy.ServiceID, ServiceType: legacy.Code, Price: legacy.Price, OptionTitle: &historicalTitle, PreferredDate: &historicalDate}
		if e := tx.Create(&historical).Error; e != nil {
			return e
		}
		// Later catalog edits must not overwrite a saved historical snapshot or be reset on restart.
		tx.Model(&models.BookingOption{}).Where("code=?", bonus.Code).Update("price", 2300)
		if e := database.Migrate(tx); e != nil {
			return e
		}
		after, e := repo.BookingOption(bonus.Code)
		if e != nil {
			return e
		}
		if *after.Price != 2300 {
			t.Fatal("migration overwrote admin price")
		}
		var old models.Booking
		if e := tx.First(&old, historical.ID).Error; e != nil {
			return e
		}
		if old.Price == nil || *old.Price != 7000 || old.OptionTitle == nil || *old.OptionTitle != historicalTitle || old.PreferredDate == nil {
			t.Fatal("historical booking no longer readable")
		}
		input.ServiceID = base.ServiceID
		input.ServiceType = base.Code
		input.Price = base.Price
		if e := worker.Submit(input, true, key+"-0-true"); e != nil {
			t.Fatalf("saved retry after price edit: %v", e)
		}
		var saved models.Booking
		tx.Where("submission_key=?", key+"-0-true").First(&saved)
		if *saved.BonusPrice != 2000 || *saved.TotalPrice != 7000 {
			t.Fatal("snapshot mutated")
		}
		return rollback
	})
	if !errors.Is(err, rollback) {
		t.Fatal(err)
	}
}
