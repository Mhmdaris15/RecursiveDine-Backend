package controllers

import (
	"net/http"
	"strconv"

	"recursiveDine/internal/repositories"
	"recursiveDine/internal/services"

	"github.com/gin-gonic/gin"
)

type MenuController struct {
	menuService *services.MenuService
}

func NewMenuController(menuService *services.MenuService) *MenuController {
	return &MenuController{
		menuService: menuService,
	}
}

// @Summary Get complete menu
// @Description Get all menu categories with their items
// @Tags menu
// @Accept json
// @Produce json
// @Success 200 {array} repositories.MenuCategory
// @Failure 500 {object} map[string]string
// @Router /menu [get]
func (ctrl *MenuController) GetMenu(c *gin.Context) {
	menu, err := ctrl.menuService.GetCompleteMenu()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, menu)
}

// @Summary Get menu categories
// @Description Get all menu categories
// @Tags menu
// @Accept json
// @Produce json
// @Success 200 {array} repositories.MenuCategory
// @Failure 500 {object} map[string]string
// @Router /menu/categories [get]
func (ctrl *MenuController) GetCategories(c *gin.Context) {
	categories, err := ctrl.menuService.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// @Summary Get menu items by category
// @Description Get all menu items for a specific category
// @Tags menu
// @Accept json
// @Produce json
// @Param category_id query int true "Category ID"
// @Success 200 {array} repositories.MenuItem
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu/items [get]
func (ctrl *MenuController) GetMenuItemsByCategory(c *gin.Context) {
	var req struct {
		CategoryID uint `form:"category_id" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items, err := ctrl.menuService.GetMenuItemsByCategory(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// @Summary Search menu items
// @Description Search menu items by name or description
// @Tags menu
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {array} repositories.MenuItem
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /menu/items/search [get]
func (ctrl *MenuController) SearchItems(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	items, err := ctrl.menuService.SearchMenuItems(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

// @Summary Get menu item by ID
// @Description Get detailed information about a specific menu item
// @Tags menu
// @Accept json
// @Produce json
// @Param id path int true "Menu Item ID"
// @Success 200 {object} repositories.MenuItem
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /menu/items/{id} [get]
func (ctrl *MenuController) GetMenuItemByID(c *gin.Context) {
	var req struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := ctrl.menuService.GetMenuItemByID(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// CRUD operations for categories (Admin only)

// @Summary Create menu category
// @Description Create a new menu category (admin only)
// @Tags menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body repositories.MenuCategory true "Category data"
// @Success 201 {object} repositories.MenuCategory
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/menu/categories [post]
func (ctrl *MenuController) CreateCategory(c *gin.Context) {
	var category repositories.MenuCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.menuService.CreateCategory(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// @Summary Update menu category
// @Description Update menu category (admin only)
// @Tags menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Param request body repositories.MenuCategory true "Category data"
// @Success 200 {object} repositories.MenuCategory
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/menu/categories/{id} [put]
func (ctrl *MenuController) UpdateCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category repositories.MenuCategory
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.ID = uint(categoryID)
	if err := ctrl.menuService.UpdateCategory(&category); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// @Summary Delete menu category
// @Description Soft delete menu category (admin only)
// @Tags menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/menu/categories/{id} [delete]
func (ctrl *MenuController) DeleteCategory(c *gin.Context) {
	categoryID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	if err := ctrl.menuService.DeleteCategory(uint(categoryID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// CRUD operations for menu items (Admin only)

// @Summary Create menu item
// @Description Create a new menu item (admin only)
// @Tags menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body repositories.MenuItem true "Menu item data"
// @Success 201 {object} repositories.MenuItem
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/menu/items [post]
func (ctrl *MenuController) CreateMenuItem(c *gin.Context) {
	var item repositories.MenuItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.menuService.CreateMenuItem(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// @Summary Update menu item
// @Description Update menu item (admin only)
// @Tags menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Menu item ID"
// @Param request body repositories.MenuItem true "Menu item data"
// @Success 200 {object} repositories.MenuItem
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/menu/items/{id} [put]
func (ctrl *MenuController) UpdateMenuItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item ID"})
		return
	}

	var item repositories.MenuItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.ID = uint(itemID)
	if err := ctrl.menuService.UpdateMenuItem(&item); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary Delete menu item
// @Description Soft delete menu item (admin only)
// @Tags menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Menu item ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/menu/items/{id} [delete]
func (ctrl *MenuController) DeleteMenuItem(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item ID"})
		return
	}

	if err := ctrl.menuService.DeleteMenuItem(uint(itemID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu item deleted successfully"})
}

// @Summary Update menu item availability
// @Description Update menu item availability (admin/staff only)
// @Tags menu
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Menu item ID"
// @Param request body map[string]bool true "Availability status"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/menu/items/{id}/availability [patch]
func (ctrl *MenuController) UpdateMenuItemAvailability(c *gin.Context) {
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item ID"})
		return
	}

	var req struct {
		IsAvailable bool `json:"is_available"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.menuService.UpdateMenuItemAvailability(uint(itemID), req.IsAvailable); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	status := "unavailable"
	if req.IsAvailable {
		status = "available"
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu item marked as " + status})
}
