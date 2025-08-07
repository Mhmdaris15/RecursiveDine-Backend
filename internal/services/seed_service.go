package services

import (
	"errors"
	"fmt"

	"recursiveDine/internal/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SeedService struct {
	db           *gorm.DB
	userRepo     *repositories.UserRepository
	tableRepo    *repositories.TableRepository
	menuRepo     *repositories.MenuRepository
}

type SeedResponse struct {
	Message string                 `json:"message"`
	Results map[string]interface{} `json:"results"`
}

func NewSeedService(db *gorm.DB, userRepo *repositories.UserRepository, tableRepo *repositories.TableRepository, menuRepo *repositories.MenuRepository) *SeedService {
	return &SeedService{
		db:        db,
		userRepo:  userRepo,
		tableRepo: tableRepo,
		menuRepo:  menuRepo,
	}
}

func (s *SeedService) SeedAll() (*SeedResponse, error) {
	results := make(map[string]interface{})
	
	// Seed users
	userCount, err := s.SeedUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to seed users: %v", err)
	}
	results["users"] = map[string]interface{}{
		"created": userCount,
		"message": "Users seeded successfully",
	}

	// Seed tables
	tableCount, err := s.SeedTables()
	if err != nil {
		return nil, fmt.Errorf("failed to seed tables: %v", err)
	}
	results["tables"] = map[string]interface{}{
		"created": tableCount,
		"message": "Tables seeded successfully",
	}

	// Seed menu categories
	categoryCount, err := s.SeedMenuCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to seed menu categories: %v", err)
	}
	results["categories"] = map[string]interface{}{
		"created": categoryCount,
		"message": "Menu categories seeded successfully",
	}

	// Seed menu items
	itemCount, err := s.SeedMenuItems()
	if err != nil {
		return nil, fmt.Errorf("failed to seed menu items: %v", err)
	}
	results["menu_items"] = map[string]interface{}{
		"created": itemCount,
		"message": "Menu items seeded successfully",
	}

	return &SeedResponse{
		Message: "Database seeded successfully",
		Results: results,
	}, nil
}

func (s *SeedService) SeedUsers() (int, error) {
	// Check if users already exist
	var count int64
	if err := s.db.Model(&repositories.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	
	if count > 0 {
		return 0, nil // Users already exist, skip seeding
	}

	// Create password hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	users := []repositories.User{
		{
			Name:     "Administrator",
			Username: "admin",
			Email:    "admin@recursivedine.com",
			Phone:    "+1234567001",
			Password: string(hashedPassword),
			Role:     repositories.RoleAdmin,
			IsActive: true,
		},
		{
			Name:     "Staff Member",
			Username: "staff1",
			Email:    "staff1@recursivedine.com",
			Phone:    "+1234567002",
			Password: string(hashedPassword),
			Role:     repositories.RoleStaff,
			IsActive: true,
		},
		{
			Name:     "Kitchen Staff",
			Username: "kitchen1",
			Email:    "kitchen1@recursivedine.com",
			Phone:    "+1234567003",
			Password: string(hashedPassword),
			Role:     repositories.RoleStaff,
			IsActive: true,
		},
		{
			Name:     "Cashier",
			Username: "cashier1",
			Email:    "cashier1@recursivedine.com",
			Phone:    "+1234567004",
			Password: string(hashedPassword),
			Role:     repositories.RoleCashier,
			IsActive: true,
		},
		{
			Name:     "John Doe",
			Username: "customer1",
			Email:    "customer1@example.com",
			Phone:    "+1234567101",
			Password: string(hashedPassword),
			Role:     repositories.RoleCustomer,
			IsActive: true,
		},
		{
			Name:     "Jane Smith",
			Username: "customer2",
			Email:    "customer2@example.com",
			Phone:    "+1234567102",
			Password: string(hashedPassword),
			Role:     repositories.RoleCustomer,
			IsActive: true,
		},
	}

	for i := range users {
		if err := s.userRepo.Create(&users[i]); err != nil {
			return i, err
		}
	}

	return len(users), nil
}

func (s *SeedService) SeedTables() (int, error) {
	// Check if tables already exist
	var count int64
	if err := s.db.Model(&repositories.Table{}).Count(&count).Error; err != nil {
		return 0, err
	}
	
	if count > 0 {
		return 0, nil // Tables already exist, skip seeding
	}

	tables := []repositories.Table{
		{Number: 1, QRCode: "QR_TABLE_001", Capacity: 2, IsAvailable: true},
		{Number: 2, QRCode: "QR_TABLE_002", Capacity: 4, IsAvailable: true},
		{Number: 3, QRCode: "QR_TABLE_003", Capacity: 4, IsAvailable: true},
		{Number: 4, QRCode: "QR_TABLE_004", Capacity: 6, IsAvailable: true},
		{Number: 5, QRCode: "QR_TABLE_005", Capacity: 8, IsAvailable: true},
		{Number: 6, QRCode: "QR_TABLE_006", Capacity: 2, IsAvailable: true},
		{Number: 7, QRCode: "QR_TABLE_007", Capacity: 4, IsAvailable: true},
		{Number: 8, QRCode: "QR_TABLE_008", Capacity: 6, IsAvailable: true},
		{Number: 9, QRCode: "QR_TABLE_009", Capacity: 10, IsAvailable: true},
		{Number: 10, QRCode: "QR_TABLE_010", Capacity: 12, IsAvailable: true},
	}

	for i := range tables {
		if err := s.db.Create(&tables[i]).Error; err != nil {
			return i, err
		}
	}

	return len(tables), nil
}

func (s *SeedService) SeedMenuCategories() (int, error) {
	// Check if categories already exist
	var count int64
	if err := s.db.Model(&repositories.MenuCategory{}).Count(&count).Error; err != nil {
		return 0, err
	}
	
	if count > 0 {
		return 0, nil // Categories already exist, skip seeding
	}

	categories := []repositories.MenuCategory{
		{
			Name:        "Appetizers",
			Description: "Start your meal with our delicious appetizers",
			IsActive:    true,
			SortOrder:   1,
		},
		{
			Name:        "Main Courses",
			Description: "Hearty main dishes to satisfy your hunger",
			IsActive:    true,
			SortOrder:   2,
		},
		{
			Name:        "Desserts",
			Description: "Sweet endings to your meal",
			IsActive:    true,
			SortOrder:   3,
		},
		{
			Name:        "Beverages",
			Description: "Refreshing drinks and beverages",
			IsActive:    true,
			SortOrder:   4,
		},
		{
			Name:        "Salads",
			Description: "Fresh and healthy salad options",
			IsActive:    true,
			SortOrder:   5,
		},
	}

	for i := range categories {
		if err := s.db.Create(&categories[i]).Error; err != nil {
			return i, err
		}
	}

	return len(categories), nil
}

func (s *SeedService) SeedMenuItems() (int, error) {
	// Check if menu items already exist
	var count int64
	if err := s.db.Model(&repositories.MenuItem{}).Count(&count).Error; err != nil {
		return 0, err
	}
	
	if count > 0 {
		return 0, nil // Items already exist, skip seeding
	}

	// Get categories first
	var categories []repositories.MenuCategory
	if err := s.db.Order("sort_order").Find(&categories).Error; err != nil {
		return 0, err
	}

	if len(categories) == 0 {
		return 0, errors.New("no categories found, please seed categories first")
	}

	// Create a map for easy category lookup
	categoryMap := make(map[string]uint)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat.ID
	}

	menuItems := []repositories.MenuItem{
		// Appetizers
		{
			CategoryID:  categoryMap["Appetizers"],
			Name:        "Spring Rolls",
			Description: "Crispy spring rolls with fresh vegetables and sweet chili sauce",
			Price:       8.99,
			ImageURL:    "https://example.com/images/spring-rolls.jpg",
			IsAvailable: true,
			SortOrder:   1,
		},
		{
			CategoryID:  categoryMap["Appetizers"],
			Name:        "Chicken Wings",
			Description: "Spicy buffalo chicken wings with blue cheese dip",
			Price:       12.99,
			ImageURL:    "https://example.com/images/chicken-wings.jpg",
			IsAvailable: true,
			SortOrder:   2,
		},
		{
			CategoryID:  categoryMap["Appetizers"],
			Name:        "Mozzarella Sticks",
			Description: "Golden fried mozzarella sticks with marinara sauce",
			Price:       9.99,
			ImageURL:    "https://example.com/images/mozzarella-sticks.jpg",
			IsAvailable: true,
			SortOrder:   3,
		},

		// Main Courses
		{
			CategoryID:  categoryMap["Main Courses"],
			Name:        "Grilled Salmon",
			Description: "Fresh Atlantic salmon grilled with herbs and lemon",
			Price:       24.99,
			ImageURL:    "https://example.com/images/grilled-salmon.jpg",
			IsAvailable: true,
			SortOrder:   1,
		},
		{
			CategoryID:  categoryMap["Main Courses"],
			Name:        "Beef Steak",
			Description: "Tender 8oz ribeye steak with garlic butter",
			Price:       28.99,
			ImageURL:    "https://example.com/images/beef-steak.jpg",
			IsAvailable: true,
			SortOrder:   2,
		},
		{
			CategoryID:  categoryMap["Main Courses"],
			Name:        "Pasta Carbonara",
			Description: "Creamy pasta with bacon, eggs, and parmesan cheese",
			Price:       18.99,
			ImageURL:    "https://example.com/images/pasta-carbonara.jpg",
			IsAvailable: true,
			SortOrder:   3,
		},
		{
			CategoryID:  categoryMap["Main Courses"],
			Name:        "Chicken Parmesan",
			Description: "Breaded chicken breast with marinara sauce and melted cheese",
			Price:       22.99,
			ImageURL:    "https://example.com/images/chicken-parmesan.jpg",
			IsAvailable: true,
			SortOrder:   4,
		},

		// Desserts
		{
			CategoryID:  categoryMap["Desserts"],
			Name:        "Chocolate Cake",
			Description: "Rich chocolate cake with chocolate frosting and berries",
			Price:       7.99,
			ImageURL:    "https://example.com/images/chocolate-cake.jpg",
			IsAvailable: true,
			SortOrder:   1,
		},
		{
			CategoryID:  categoryMap["Desserts"],
			Name:        "Ice Cream",
			Description: "Vanilla ice cream with your choice of toppings",
			Price:       5.99,
			ImageURL:    "https://example.com/images/ice-cream.jpg",
			IsAvailable: true,
			SortOrder:   2,
		},
		{
			CategoryID:  categoryMap["Desserts"],
			Name:        "Tiramisu",
			Description: "Classic Italian tiramisu with coffee and mascarpone",
			Price:       8.99,
			ImageURL:    "https://example.com/images/tiramisu.jpg",
			IsAvailable: true,
			SortOrder:   3,
		},

		// Beverages
		{
			CategoryID:  categoryMap["Beverages"],
			Name:        "Coffee",
			Description: "Freshly brewed premium coffee",
			Price:       3.99,
			ImageURL:    "https://example.com/images/coffee.jpg",
			IsAvailable: true,
			SortOrder:   1,
		},
		{
			CategoryID:  categoryMap["Beverages"],
			Name:        "Fresh Juice",
			Description: "Freshly squeezed orange or apple juice",
			Price:       4.99,
			ImageURL:    "https://example.com/images/fresh-juice.jpg",
			IsAvailable: true,
			SortOrder:   2,
		},
		{
			CategoryID:  categoryMap["Beverages"],
			Name:        "Soft Drinks",
			Description: "Coca-Cola, Sprite, Fanta, or Pepsi",
			Price:       2.99,
			ImageURL:    "https://example.com/images/soft-drinks.jpg",
			IsAvailable: true,
			SortOrder:   3,
		},
		{
			CategoryID:  categoryMap["Beverages"],
			Name:        "Craft Beer",
			Description: "Local craft beer selection",
			Price:       6.99,
			ImageURL:    "https://example.com/images/craft-beer.jpg",
			IsAvailable: true,
			SortOrder:   4,
		},

		// Salads
		{
			CategoryID:  categoryMap["Salads"],
			Name:        "Caesar Salad",
			Description: "Crisp romaine lettuce with Caesar dressing and croutons",
			Price:       11.99,
			ImageURL:    "https://example.com/images/caesar-salad.jpg",
			IsAvailable: true,
			SortOrder:   1,
		},
		{
			CategoryID:  categoryMap["Salads"],
			Name:        "Greek Salad",
			Description: "Fresh vegetables with feta cheese and olive oil",
			Price:       12.99,
			ImageURL:    "https://example.com/images/greek-salad.jpg",
			IsAvailable: true,
			SortOrder:   2,
		},
	}

	for i := range menuItems {
		if err := s.db.Create(&menuItems[i]).Error; err != nil {
			return i, err
		}
	}

	return len(menuItems), nil
}

func (s *SeedService) ClearAll() error {
	// Clear in reverse order to respect foreign key constraints
	if err := s.db.Exec("DELETE FROM order_items").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM payments").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM orders").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM menu_items").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM menu_categories").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM tables").Error; err != nil {
		return err
	}
	if err := s.db.Exec("DELETE FROM users").Error; err != nil {
		return err
	}

	// Reset sequences
	if err := s.db.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE tables_id_seq RESTART WITH 1").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE menu_categories_id_seq RESTART WITH 1").Error; err != nil {
		return err
	}
	if err := s.db.Exec("ALTER SEQUENCE menu_items_id_seq RESTART WITH 1").Error; err != nil {
		return err
	}

	return nil
}
