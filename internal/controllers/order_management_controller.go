package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"recursiveDine/internal/repositories"
	"recursiveDine/internal/services"
)

type OrderManagementController struct {
	orderService *services.OrderService
}

func NewOrderManagementController(orderService *services.OrderService) *OrderManagementController {
	return &OrderManagementController{
		orderService: orderService,
	}
}

// @Summary Get all orders
// @Description Get all orders with pagination (staff/admin only)
// @Tags orders-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/orders [get]
func (ctrl *OrderManagementController) GetAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	orders, total, err := ctrl.orderService.GetAllOrdersAdmin(page, limit, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":      orders,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}

// @Summary Get order by ID
// @Description Get order details by ID (staff/admin only)
// @Tags orders-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} repositories.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/orders/{id} [get]
func (ctrl *OrderManagementController) GetOrderByID(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := ctrl.orderService.GetOrderByIDAdmin(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// @Summary Update order status
// @Description Update order status (staff/admin only)
// @Tags orders-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param request body map[string]string true "Order status"
// @Success 200 {object} repositories.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/orders/{id}/status [patch]
func (ctrl *OrderManagementController) UpdateOrderStatus(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	validStatuses := []string{"pending", "confirmed", "preparing", "ready", "delivered", "cancelled"}
	isValid := false
	for _, status := range validStatuses {
		if status == req.Status {
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order status"})
		return
	}

	order, err := ctrl.orderService.UpdateOrderStatusAdmin(uint(orderID), req.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// @Summary Delete order
// @Description Soft delete order (admin only)
// @Tags orders-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/orders/{id} [delete]
func (ctrl *OrderManagementController) DeleteOrder(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if err := ctrl.orderService.DeleteOrder(uint(orderID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}

// @Summary Get order statistics
// @Description Get order statistics (admin only)
// @Tags orders-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/orders/statistics [get]
func (ctrl *OrderManagementController) GetOrderStatistics(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	stats, err := ctrl.orderService.GetOrderStatistics(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary Get daily revenue
// @Description Get daily revenue report (admin only)
// @Tags orders-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/orders/revenue [get]
func (ctrl *OrderManagementController) GetDailyRevenue(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	revenue, err := ctrl.orderService.GetDailyRevenue(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, revenue)
}

// @Summary Update order items
// @Description Update order items (staff/admin only, only for pending orders)
// @Tags orders-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param request body []repositories.OrderItem true "Order items"
// @Success 200 {object} repositories.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/orders/{id}/items [put]
func (ctrl *OrderManagementController) UpdateOrderItems(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var items []repositories.OrderItem
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := ctrl.orderService.UpdateOrderItems(uint(orderID), items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
