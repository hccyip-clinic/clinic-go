package main

import (
	"fmt"
	"log"
	"net/http"

	"clinic-hcc-app/internal/config"
	"clinic-hcc-app/internal/database"
	"clinic-hcc-app/internal/handlers"
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

	router := handlers.NewRouter(db)
	mux := router.Setup()

	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Starting server on %s (database: %s)", addr, cfg.DatabasePath)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}