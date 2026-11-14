package main

import (
	"context"
	"fmt"
	"os"
	"product-service/config"
	"product-service/internal/app"
	"product-service/internal/infrastructure/database"
)

func main() {
	if err := Start(); err != nil {
		fmt.Printf("Failed to start service, error %v", err)
		os.Exit(1)
	}
}

func Start() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := database.Connect(context.Background(), cfg.DB.URL())
	if err != nil {
		return err
	}
	defer db.Close()

	server := app.NewApp(db)
	if err := server.Run(":" + cfg.AppPort); err != nil {
		return err
	}
	return nil
}
