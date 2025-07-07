package main

import (
	"log"
	"fmt"
	"golang.org/x/crypto/bcrypt"
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

	// Seed sample data
	if err := seedUsers(db); err != nil {
		log.Fatal("Failed to seed users:", err)
	}

	if err := seedTables(db); err != nil {
		log.Fatal("Failed to seed tables:", err)
	}

	if err := seedMenu(db); err != nil {
		log.Fatal("Failed to seed menu:", err)
	}

	fmt.Println("Database seeding completed successfully!")
}

func seedUsers(db *gorm.DB) error {
	users := []repositories.User{
		{
			Username: "admin",
			Email:    "admin@recursiveDine.com",
			Password: hashPassword("admin123"),
			Role:     repositories.RoleAdmin,
			IsActive: true,
		},
		{
			Username: "staff1",
			Email:    "staff1@recursiveDine.com",
			Password: hashPassword("staff123"),
			Role:     repositories.RoleStaff,
			IsActive: true,
		},
		{
			Username: "customer1",
			Email:    "customer1@example.com",
			Password: hashPassword("customer123"),
			Role:     repositories.RoleCustomer,
			IsActive: true,
		},
	}

	for _, user := range users {
		var existingUser repositories.User
		if err := db.Where("username = ?", user.Username).First(&existingUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&user).Error; err != nil {
					return err
				}
				fmt.Printf("Created user: %s\n", user.Username)
			} else {
				return err
			}
		}
	}

	return nil
}

func seedTables(db *gorm.DB) error {
	tables := []repositories.Table{
		{Number: 1, QRCode: "QR001", Capacity: 4, IsAvailable: true},
		{Number: 2, QRCode: "QR002", Capacity: 2, IsAvailable: true},
		{Number: 3, QRCode: "QR003", Capacity: 6, IsAvailable: true},
		{Number: 4, QRCode: "QR004", Capacity: 4, IsAvailable: true},
		{Number: 5, QRCode: "QR005", Capacity: 8, IsAvailable: true},
	}

	for _, table := range tables {
		var existingTable repositories.Table
		if err := db.Where("number = ?", table.Number).First(&existingTable).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&table).Error; err != nil {
					return err
				}
				fmt.Printf("Created table: %d\n", table.Number)
			} else {
				return err
			}
		}
	}

	return nil
}

func seedMenu(db *gorm.DB) error {
	// Create categories
	categories := []repositories.MenuCategory{
		{Name: "Appetizers", Description: "Start your meal with our delicious appetizers", SortOrder: 1, IsActive: true},
		{Name: "Main Courses", Description: "Hearty main dishes to satisfy your hunger", SortOrder: 2, IsActive: true},
		{Name: "Desserts", Description: "Sweet endings to your meal", SortOrder: 3, IsActive: true},
		{Name: "Beverages", Description: "Refreshing drinks and beverages", SortOrder: 4, IsActive: true},
	}

	for _, category := range categories {
		var existingCategory repositories.MenuCategory
		if err := db.Where("name = ?", category.Name).First(&existingCategory).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&category).Error; err != nil {
					return err
				}
				fmt.Printf("Created category: %s\n", category.Name)
			} else {
				return err
			}
		}
	}

	// Get category IDs
	var appetizers, mainCourses, desserts, beverages repositories.MenuCategory
	db.Where("name = ?", "Appetizers").First(&appetizers)
	db.Where("name = ?", "Main Courses").First(&mainCourses)
	db.Where("name = ?", "Desserts").First(&desserts)
	db.Where("name = ?", "Beverages").First(&beverages)

	// Create menu items
	items := []repositories.MenuItem{
		{CategoryID: appetizers.ID, Name: "Spring Rolls", Description: "Crispy spring rolls with vegetables", Price: 8.99, SortOrder: 1, IsAvailable: true},
		{CategoryID: appetizers.ID, Name: "Chicken Wings", Description: "Spicy buffalo chicken wings", Price: 12.99, SortOrder: 2, IsAvailable: true},
		{CategoryID: mainCourses.ID, Name: "Grilled Salmon", Description: "Fresh salmon with herbs and lemon", Price: 24.99, SortOrder: 1, IsAvailable: true},
		{CategoryID: mainCourses.ID, Name: "Beef Steak", Description: "Tender beef steak with garlic butter", Price: 28.99, SortOrder: 2, IsAvailable: true},
		{CategoryID: mainCourses.ID, Name: "Pasta Carbonara", Description: "Creamy pasta with bacon and eggs", Price: 18.99, SortOrder: 3, IsAvailable: true},
		{CategoryID: desserts.ID, Name: "Chocolate Cake", Description: "Rich chocolate cake with frosting", Price: 7.99, SortOrder: 1, IsAvailable: true},
		{CategoryID: desserts.ID, Name: "Ice Cream", Description: "Vanilla ice cream with toppings", Price: 5.99, SortOrder: 2, IsAvailable: true},
		{CategoryID: beverages.ID, Name: "Coffee", Description: "Freshly brewed coffee", Price: 3.99, SortOrder: 1, IsAvailable: true},
		{CategoryID: beverages.ID, Name: "Fresh Juice", Description: "Orange or apple juice", Price: 4.99, SortOrder: 2, IsAvailable: true},
		{CategoryID: beverages.ID, Name: "Soft Drinks", Description: "Coca-Cola, Sprite, or Fanta", Price: 2.99, SortOrder: 3, IsAvailable: true},
	}

	for _, item := range items {
		var existingItem repositories.MenuItem
		if err := db.Where("name = ? AND category_id = ?", item.Name, item.CategoryID).First(&existingItem).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&item).Error; err != nil {
					return err
				}
				fmt.Printf("Created menu item: %s\n", item.Name)
			} else {
				return err
			}
		}
	}

	return nil
}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}
