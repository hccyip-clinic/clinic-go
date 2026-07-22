package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"clinic-hcc-app/internal/config"
	"clinic-hcc-app/internal/database"
	"clinic-hcc-app/internal/handlers"
	"clinic-hcc-app/internal/security"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.New(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	if len(os.Args) > 1 && os.Args[1] == "reset-password" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Fprint(os.Stderr, "New password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)
		if len(password) < 12 {
			log.Fatal("password must be at least 12 characters")
		}
		hash, err := security.HashPassword(password)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := db.Exec(`UPDATE settings SET password_hash = ?, updated_at = CURRENT_TIMESTAMP WHERE id = 1`, hash); err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(os.Stderr, "Password reset.")
		return
	}

	router := handlers.NewRouter(db)
	mux := router.Setup()

	bind := cfg.BindAddress
	if bind == "" {
		bind = "127.0.0.1"
	}
	addr := fmt.Sprintf("%s:%d", bind, cfg.ServerPort)
	log.Printf("Starting server on %s (database: %s)", addr, cfg.DatabasePath)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
