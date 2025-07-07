package main

import (
	"fmt"
	"log"
	"os"

	"recursiveDine/internal/config"
	"recursiveDine/internal/repositories"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate models
	err = db.AutoMigrate(
		&repositories.User{},
		&repositories.Table{},
		&repositories.MenuCategory{},
		&repositories.MenuItem{},
		&repositories.Order{},
		&repositories.OrderItem{},
		&repositories.Payment{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migration completed successfully!")
	os.Exit(0)
}
