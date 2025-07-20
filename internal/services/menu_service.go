package services

import (
	"recursiveDine/internal/repositories"
)

type MenuService struct {
	menuRepo *repositories.MenuRepository
}

func NewMenuService(menuRepo *repositories.MenuRepository) *MenuService {
	return &MenuService{
		menuRepo: menuRepo,
	}
}

func (s *MenuService) GetCompleteMenu() ([]repositories.MenuCategory, error) {
	return s.menuRepo.GetCompleteMenu()
}

func (s *MenuService) GetAllCategories() ([]repositories.MenuCategory, error) {
	return s.menuRepo.GetAllCategories()
}

func (s *MenuService) GetCategoryByID(id uint) (*repositories.MenuCategory, error) {
	return s.menuRepo.GetCategoryByID(id)
}

func (s *MenuService) GetAllMenuItems() ([]repositories.MenuItem, error) {
	return s.menuRepo.GetAllMenuItems()
}

func (s *MenuService) GetMenuItemByID(id uint) (*repositories.MenuItem, error) {
	return s.menuRepo.GetMenuItemByID(id)
}

// Category CRUD operations

func (s *MenuService) CreateCategory(category *repositories.MenuCategory) error {
	return s.menuRepo.CreateCategory(category)
}

func (s *MenuService) UpdateCategory(category *repositories.MenuCategory) error {
	return s.menuRepo.UpdateCategory(category)
}

func (s *MenuService) DeleteCategory(id uint) error {
	return s.menuRepo.DeleteCategory(id)
}

// Menu Item CRUD operations

func (s *MenuService) CreateMenuItem(item *repositories.MenuItem) error {
	return s.menuRepo.CreateMenuItem(item)
}

func (s *MenuService) UpdateMenuItem(item *repositories.MenuItem) error {
	return s.menuRepo.UpdateMenuItem(item)
}

func (s *MenuService) DeleteMenuItem(id uint) error {
	return s.menuRepo.DeleteMenuItem(id)
}

func (s *MenuService) UpdateMenuItemAvailability(id uint, available bool) error {
	return s.menuRepo.UpdateMenuItemAvailability(id, available)
}

func (s *MenuService) GetMenuItemsByCategory(categoryID uint) ([]repositories.MenuItem, error) {
	return s.menuRepo.GetMenuItemsByCategory(categoryID)
}

func (s *MenuService) SearchMenuItems(query string) ([]repositories.MenuItem, error) {
	return s.menuRepo.SearchMenuItems(query)
}

func (s *MenuService) GetMenuItemsByIDs(ids []uint) ([]repositories.MenuItem, error) {
	return s.menuRepo.GetMenuItemsByIDs(ids)
}
