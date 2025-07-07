package controllers

import (
	"net/http"

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
