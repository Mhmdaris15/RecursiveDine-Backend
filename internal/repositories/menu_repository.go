package repositories

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

// Category operations
func (r *MenuRepository) CreateCategory(category *MenuCategory) error {
	return r.db.Create(category).Error
}

func (r *MenuRepository) GetCategoryByID(id uint) (*MenuCategory, error) {
	var category MenuCategory
	err := r.db.Preload("MenuItems").First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("category not found")
	}
	return &category, err
}

func (r *MenuRepository) GetAllCategories() ([]MenuCategory, error) {
	var categories []MenuCategory
	err := r.db.Where("is_active = ?", true).
		Order("sort_order ASC").
		Preload("MenuItems", "is_available = ?", true).
		Find(&categories).Error
	return categories, err
}

func (r *MenuRepository) UpdateCategory(category *MenuCategory) error {
	return r.db.Save(category).Error
}

func (r *MenuRepository) DeleteCategory(id uint) error {
	return r.db.Delete(&MenuCategory{}, id).Error
}

// MenuItem operations
func (r *MenuRepository) CreateMenuItem(item *MenuItem) error {
	return r.db.Create(item).Error
}

func (r *MenuRepository) GetMenuItemByID(id uint) (*MenuItem, error) {
	var item MenuItem
	err := r.db.Preload("Category").First(&item, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("menu item not found")
	}
	return &item, err
}

func (r *MenuRepository) GetAllMenuItems() ([]MenuItem, error) {
	var items []MenuItem
	err := r.db.Where("is_available = ?", true).
		Order("sort_order ASC").
		Preload("Category").
		Find(&items).Error
	return items, err
}

func (r *MenuRepository) GetMenuItemsByCategory(categoryID uint) ([]MenuItem, error) {
	var items []MenuItem
	err := r.db.Where("category_id = ? AND is_available = ?", categoryID, true).
		Order("sort_order ASC").
		Preload("Category").
		Find(&items).Error
	return items, err
}

func (r *MenuRepository) SearchMenuItems(query string) ([]MenuItem, error) {
	var items []MenuItem
	searchQuery := "%" + strings.ToLower(query) + "%"
	err := r.db.Where("is_available = ? AND (LOWER(name) LIKE ? OR LOWER(description) LIKE ?)", 
		true, searchQuery, searchQuery).
		Order("sort_order ASC").
		Preload("Category").
		Find(&items).Error
	return items, err
}

func (r *MenuRepository) UpdateMenuItem(item *MenuItem) error {
	return r.db.Save(item).Error
}

func (r *MenuRepository) DeleteMenuItem(id uint) error {
	return r.db.Delete(&MenuItem{}, id).Error
}

func (r *MenuRepository) GetMenuItemsByIDs(ids []uint) ([]MenuItem, error) {
	var items []MenuItem
	err := r.db.Where("id IN ? AND is_available = ?", ids, true).
		Preload("Category").
		Find(&items).Error
	return items, err
}

func (r *MenuRepository) GetCompleteMenu() ([]MenuCategory, error) {
	var categories []MenuCategory
	err := r.db.Where("is_active = ?", true).
		Order("sort_order ASC").
		Preload("MenuItems", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_available = ?", true).Order("sort_order ASC")
		}).
		Find(&categories).Error
	return categories, err
}
