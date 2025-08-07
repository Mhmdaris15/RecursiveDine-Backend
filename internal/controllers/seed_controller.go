package controllers

import (
	"net/http"

	"recursiveDine/internal/services"

	"github.com/gin-gonic/gin"
)

type SeedController struct {
	seedService *services.SeedService
}

func NewSeedController(seedService *services.SeedService) *SeedController {
	return &SeedController{
		seedService: seedService,
	}
}

// SeedDatabase seeds the database with sample data
// @Summary Seed database with sample data
// @Description Seeds the database with users, tables, menu categories, and menu items
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} services.SeedResponse
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/seed [post]
func (sc *SeedController) SeedDatabase(c *gin.Context) {
	// Check if user is admin
	// userRole, exists := c.Get("userRole")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	// 	return
	// }

	// if userRole != "admin" {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admin role required."})
	// 	return
	// }

	result, err := sc.seedService.SeedAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ClearDatabase clears all data from the database
// @Summary Clear all data from database
// @Description Removes all data from the database (DANGER: Use with caution!)
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/clear [delete]
func (sc *SeedController) ClearDatabase(c *gin.Context) {
	// Check if user is admin
	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admin role required."})
		return
	}

	// Additional safety check - require confirmation header
	confirmation := c.GetHeader("X-Confirm-Clear")
	if confirmation != "YES-CLEAR-ALL-DATA" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing confirmation header. Add header 'X-Confirm-Clear: YES-CLEAR-ALL-DATA' to confirm.",
		})
		return
	}

	err := sc.seedService.ClearAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Database cleared successfully",
		"warning": "All data has been permanently deleted",
	})
}

// GetSeedStatus shows current database status
// @Summary Get database seed status
// @Description Shows the current count of records in each table
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/seed/status [get]
func (sc *SeedController) GetSeedStatus(c *gin.Context) {
	// Check if user is admin
	userRole, exists := c.Get("userRole")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admin role required."})
		return
	}

	// This would be implemented in the seed service if needed
	c.JSON(http.StatusOK, gin.H{
		"message": "Seed status endpoint - implement if needed",
		"hint":    "Could show counts of users, tables, categories, menu items, etc.",
	})
}
