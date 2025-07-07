package controllers

import (
	"net/http"

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
