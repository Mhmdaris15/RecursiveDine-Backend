package controllers

import (
	"net/http"
	"strconv"

	"recursiveDine/internal/repositories"
	"recursiveDine/internal/services"

	"github.com/gin-gonic/gin"
)

type TableController struct {
	tableService *services.TableService
}

func NewTableController(tableService *services.TableService) *TableController {
	return &TableController{
		tableService: tableService,
	}
}

// @Summary Get table by QR code
// @Description Get table information by scanning QR code
// @Tags tables
// @Accept json
// @Produce json
// @Param qr_code path string true "Table QR Code"
// @Success 200 {object} repositories.Table
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /tables/{qr_code} [get]
func (ctrl *TableController) GetTableByQRCode(c *gin.Context) {
	qrCode := c.Param("qr_code")
	if qrCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QR code is required"})
		return
	}

	table, err := ctrl.tableService.GetTableByQRCode(qrCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

// @Summary Get table by ID
// @Description Get table information by ID
// @Tags tables
// @Accept json
// @Produce json
// @Param id path int true "Table ID"
// @Success 200 {object} repositories.Table
// @Failure 404 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /tables/{id} [get]
func (ctrl *TableController) GetTableByID(c *gin.Context) {
	var req struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	table, err := ctrl.tableService.GetTableByID(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

// @Summary Get all available tables
// @Description Get all available tables
// @Tags tables
// @Accept json
// @Produce json
// @Success 200 {array} repositories.Table
// @Failure 500 {object} map[string]string
// @Router /tables [get]
func (ctrl *TableController) GetAllAvailableTables(c *gin.Context) {
	tables, err := ctrl.tableService.GetAllAvailableTables()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tables)
}

// @Summary Get all tables
// @Description Get all tables with pagination (admin/staff only)
// @Tags tables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/tables [get]
func (ctrl *TableController) GetAllTables(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	tables, err := ctrl.tableService.GetAllTables(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tables": tables,
		"page":   page,
		"limit":  limit,
	})
}

// @Summary Create table
// @Description Create a new table (admin only)
// @Tags tables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body repositories.Table true "Table data"
// @Success 201 {object} repositories.Table
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /admin/tables [post]
func (ctrl *TableController) CreateTable(c *gin.Context) {
	var table repositories.Table
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.tableService.CreateTable(&table); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, table)
}

// @Summary Update table
// @Description Update table details (admin only)
// @Tags tables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Table ID"
// @Param request body repositories.Table true "Table data"
// @Success 200 {object} repositories.Table
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/tables/{id} [put]
func (ctrl *TableController) UpdateTable(c *gin.Context) {
	tableID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table ID"})
		return
	}

	var table repositories.Table
	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	table.ID = uint(tableID)
	if err := ctrl.tableService.UpdateTable(&table); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

// @Summary Delete table
// @Description Soft delete table (admin only)
// @Tags tables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Table ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/tables/{id} [delete]
func (ctrl *TableController) DeleteTable(c *gin.Context) {
	tableID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table ID"})
		return
	}

	if err := ctrl.tableService.DeleteTable(uint(tableID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Table deleted successfully"})
}

// @Summary Update table availability
// @Description Update table availability status (staff/admin only)
// @Tags tables
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Table ID"
// @Param request body map[string]bool true "Availability status"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/tables/{id}/availability [patch]
func (ctrl *TableController) UpdateTableAvailability(c *gin.Context) {
	tableID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table ID"})
		return
	}

	var req struct {
		IsAvailable bool `json:"is_available"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.tableService.SetTableAvailability(uint(tableID), req.IsAvailable); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	status := "unavailable"
	if req.IsAvailable {
		status = "available"
	}

	c.JSON(http.StatusOK, gin.H{"message": "Table marked as " + status})
}
