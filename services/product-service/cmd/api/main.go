package main

import (
	"log"

	"product-service/config"
	"product-service/internal/app"
	"product-service/internal/infrastructure"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect database
	db, err := infrastructure.NewSQLite(cfg)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// Initialize Gin server and routes
	server := app.NewApp(db)

	log.Println("Server is running on :8080")

	if err := server.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
