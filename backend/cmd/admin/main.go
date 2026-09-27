package main

import (
	"astroolog/backend/internal/config"
	"astroolog/backend/internal/database"
	"astroolog/backend/internal/models"
	"bufio"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
	"log"
	"net/mail"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 || os.Args[1] != "create" {
		log.Fatal("Usage: go run ./cmd/admin create")
	}
	cfg, e := config.Load()
	if e != nil {
		log.Fatal(e)
	}
	db := database.Connect(cfg)
	if e = database.Migrate(db); e != nil {
		log.Fatal(e)
	}
	fmt.Print("Email: ")
	email, e := bufio.NewReader(os.Stdin).ReadString('\n')
	if e != nil {
		log.Fatal(e)
	}
	email = strings.ToLower(strings.TrimSpace(email))
	a, e := mail.ParseAddress(email)
	if e != nil || a.Address != email {
		log.Fatal("Invalid email")
	}
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		log.Fatal("Password must be entered in a local interactive terminal")
	}
	fmt.Print("Password (12–72 bytes): ")
	password, e := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if e != nil || len(password) < 12 || len(password) > 72 {
		log.Fatal("Password must contain 12–72 bytes")
	}
	fmt.Print("Repeat password: ")
	repeat, e := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if e != nil || string(password) != string(repeat) {
		log.Fatal("Passwords do not match")
	}
	digest, e := bcrypt.GenerateFromPassword(password, 12)
	clear(password)
	clear(repeat)
	if e != nil {
		log.Fatal(e)
	}
	user := models.User{Email: email, PasswordHash: string(digest), Role: "admin", IsActive: true}
	if e = db.Create(&user).Error; e != nil {
		log.Fatal("Administrator was not created; check for an existing account.")
	}
	fmt.Println("Administrator created. Open /admin/login.")
}
