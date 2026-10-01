package main

import (
	"astroolog/backend/internal/models"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"io"
	"net/mail"
	"strings"
	"time"
)

var errAdminNotFound = errors.New("Administrator not found for this email.")

func findAdmin(db *gorm.DB, email string) (models.User, error) {
	var user models.User
	err := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Where("email = ? AND role = ?", email, "admin").First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errAdminNotFound
	}
	if err != nil {
		return user, errors.New("Could not find administrator.")
	}
	return user, nil
}

func passwordDigest(password, repeat []byte) ([]byte, error) {
	if len(password) < 12 || len(password) > 72 {
		return nil, errors.New("Password must contain 12–72 bytes")
	}
	if !bytes.Equal(password, repeat) {
		return nil, errors.New("Passwords do not match")
	}
	return bcrypt.GenerateFromPassword(password, 12)
}

func resetPassword(db *gorm.DB, user models.User, password, repeat []byte) error {
	digest, err := passwordDigest(password, repeat)
	if err != nil {
		return err
	}
	defer clear(digest)
	// Even an externally configured verbose GORM logger must not print a hash.
	err = db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Transaction(func(tx *gorm.DB) error {
		var current models.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND email = ? AND role = ?", user.ID, user.Email, "admin").First(&current).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.User{}).Where("id = ?", current.ID).Updates(map[string]any{"password_hash": string(digest), "updated_at": time.Now()}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", current.ID).Delete(&models.AdminSession{}).Error; err != nil {
			return err
		}
		// CLI has no authenticated web actor. user_id identifies the affected account;
		// metadata explicitly records the local CLI source, never credentials.
		return tx.Create(&models.AdminAuditLog{UserID: current.ID, Action: "password_reset", EntityType: "users", EntityID: fmt.Sprint(current.ID), Metadata: models.Document{"source": "local_cli"}}).Error
	})
	if err != nil {
		return errors.New("Password was not updated; no changes were saved.")
	}
	return nil
}

func resetInteractive(db *gorm.DB, input io.Reader, output io.Writer, terminal bool, readPassword func() ([]byte, error)) error {
	fmt.Fprint(output, "Email: ")
	email, err := bufio.NewReader(input).ReadString('\n')
	if err != nil {
		return errors.New("Could not read email.")
	}
	email = strings.ToLower(strings.TrimSpace(email))
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return errors.New("Invalid email")
	}
	user, err := findAdmin(db, email)
	if err != nil {
		return err
	}
	if !terminal {
		return errors.New("Password must be entered in a local interactive terminal")
	}
	fmt.Fprint(output, "New password (12–72 bytes): ")
	password, err := readPassword()
	fmt.Fprintln(output)
	defer clear(password)
	if err != nil {
		return errors.New("Could not read password.")
	}
	fmt.Fprint(output, "Repeat password: ")
	repeat, err := readPassword()
	fmt.Fprintln(output)
	defer clear(repeat)
	if err != nil {
		return errors.New("Could not read password confirmation.")
	}
	if err := resetPassword(db, user, password, repeat); err != nil {
		return err
	}
	fmt.Fprintln(output, "Password updated successfully.")
	return nil
}
