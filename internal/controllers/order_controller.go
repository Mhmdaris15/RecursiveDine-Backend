package controllers

import (
	"net/http"
	"strconv"

	"recursiveDine/internal/repositories"
	"recursiveDine/internal/services"
	"recursiveDine/internal/utils"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService *services.OrderService
	authService  *services.AuthService
}

func NewOrderController(orderService *services.OrderService, authService *services.AuthService) *OrderController {
	return &OrderController{
		orderService: orderService,
		authService:  authService,
	}
}

// @Summary Create new order
// @Description Create a new order for a table
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.CreateOrderRequest true "Order details"
// @Success 201 {object} repositories.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [post]
func (ctrl *OrderController) CreateOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req services.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := ctrl.orderService.CreateOrder(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// @Summary Get order by ID
// @Description Get detailed information about a specific order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} repositories.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{id} [get]
func (ctrl *OrderController) GetOrder(c *gin.Context) {
	var req struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := ctrl.orderService.GetOrderByID(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Check if user owns this order or has staff/admin role
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")

	if order.UserID != userID.(uint) && userRole != "staff" && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// @Summary Get user orders
// @Description Get all orders for the current user
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {array} repositories.Order
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [get]
func (ctrl *OrderController) GetOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userRole, _ := c.Get("user_role")

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var orders []repositories.Order
	var err error

	// If user is staff or admin, return all orders
	if userRole == "staff" || userRole == "admin" {
		orders, err = ctrl.orderService.GetAllOrders(page, limit)
	} else {
		// Return only user's orders
		orders, err = ctrl.orderService.GetOrdersByUser(userID.(uint), page, limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"page":   page,
		"limit":  limit,
	})
}

// @Summary Update order status
// @Description Update the status of an order (staff/admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param request body map[string]string true "New status"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /orders/{id}/status [patch]
func (ctrl *OrderController) UpdateOrderStatus(c *gin.Context) {
	var uriReq struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&uriReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	var status repositories.OrderStatus
	switch req.Status {
	case "pending":
		status = repositories.OrderStatusPending
	case "confirmed":
		status = repositories.OrderStatusConfirmed
	case "preparing":
		status = repositories.OrderStatusPreparing
	case "ready":
		status = repositories.OrderStatusReady
	case "served":
		status = repositories.OrderStatusServed
	case "cancelled":
		status = repositories.OrderStatusCancelled
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	if err := ctrl.orderService.UpdateOrderStatus(uriReq.ID, status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

// @Summary Get orders by status
// @Description Get all orders with a specific status (staff/admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string true "Order status"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {array} repositories.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /orders/status [get]
func (ctrl *OrderController) GetOrdersByStatus(c *gin.Context) {
	statusParam := c.Query("status")
	if statusParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status parameter is required"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Validate status
	var status repositories.OrderStatus
	switch statusParam {
	case "pending":
		status = repositories.OrderStatusPending
	case "confirmed":
		status = repositories.OrderStatusConfirmed
	case "preparing":
		status = repositories.OrderStatusPreparing
	case "ready":
		status = repositories.OrderStatusReady
	case "served":
		status = repositories.OrderStatusServed
	case "cancelled":
		status = repositories.OrderStatusCancelled
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	orders, err := ctrl.orderService.GetOrdersByStatus(status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"status": statusParam,
		"page":   page,
		"limit":  limit,
	})
}

// @Summary Get orders by order type
// @Description Get all orders filtered by order type (dine_in or takeaway)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string true "Order type (dine_in or takeaway)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {array} repositories.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /orders/type [get]
func (ctrl *OrderController) GetOrdersByType(c *gin.Context) {
	typeParam := c.Query("type")
	if typeParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type parameter is required"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Validate order type
	var orderType repositories.OrderType
	switch typeParam {
	case "dine_in":
		orderType = repositories.OrderTypeDineIn
	case "takeaway":
		orderType = repositories.OrderTypeTakeaway
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order type. Must be 'dine_in' or 'takeaway'"})
		return
	}

	orders, err := ctrl.orderService.GetOrdersByType(orderType, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"type":   typeParam,
		"page":   page,
		"limit":  limit,
	})
}

// @Summary Get takeaway orders ready for pickup
// @Description Get all takeaway orders that are ready for customer pickup
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} repositories.Order
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders/takeaway/ready [get]
func (ctrl *OrderController) GetTakeawayOrdersReady(c *gin.Context) {
	orders, err := ctrl.orderService.GetTakeawayOrdersReady()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"count":  len(orders),
	})
}

// @Summary Get orders by status and type
// @Description Get all orders filtered by both status and order type
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string true "Order status"
// @Param type query string true "Order type (dine_in or takeaway)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Success 200 {array} repositories.Order
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /orders/filter [get]
func (ctrl *OrderController) GetOrdersByStatusAndType(c *gin.Context) {
	statusParam := c.Query("status")
	typeParam := c.Query("type")

	if statusParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status parameter is required"})
		return
	}

	if typeParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type parameter is required"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Validate status
	var status repositories.OrderStatus
	switch statusParam {
	case "pending":
		status = repositories.OrderStatusPending
	case "confirmed":
		status = repositories.OrderStatusConfirmed
	case "preparing":
		status = repositories.OrderStatusPreparing
	case "ready":
		status = repositories.OrderStatusReady
	case "served":
		status = repositories.OrderStatusServed
	case "cancelled":
		status = repositories.OrderStatusCancelled
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
		return
	}

	// Validate order type
	var orderType repositories.OrderType
	switch typeParam {
	case "dine_in":
		orderType = repositories.OrderTypeDineIn
	case "takeaway":
		orderType = repositories.OrderTypeTakeaway
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order type. Must be 'dine_in' or 'takeaway'"})
		return
	}

	orders, err := ctrl.orderService.GetOrdersByStatusAndType(status, orderType, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"status": statusParam,
		"type":   typeParam,
		"page":   page,
		"limit":  limit,
	})
}
func (oc *OrderController) CreateCashierOrder(c *gin.Context) {
	var request services.CashierOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.LogError("Invalid cashier order request format", err, map[string]interface{}{
			"request_body": request,
		})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// Get cashier user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		utils.LogWarning("Cashier order attempt without authentication", map[string]interface{}{
			"customer_name": request.CustomerName,
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	cashierUserID, ok := userID.(uint)
	if !ok {
		utils.LogError("Invalid user ID type in cashier order", nil, map[string]interface{}{
			"user_id": userID,
		})
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	utils.LogInfo("Creating cashier order", map[string]interface{}{
		"cashier_user_id": cashierUserID,
		"customer_name":   request.CustomerName,
		"cashier_name":    request.CashierName,
		"table_id":        request.TableID,
		"items_count":     len(request.Items),
	})

	// Create the order through cashier service
	order, err := oc.orderService.CreateCashierOrder(cashierUserID, &request)
	if err != nil {
		utils.LogError("Failed to create cashier order", err, map[string]interface{}{
			"cashier_user_id": cashierUserID,
			"customer_name":   request.CustomerName,
			"request":         request,
		})
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.LogInfo("Cashier order created successfully", map[string]interface{}{
		"order_id":        order.ID,
		"cashier_user_id": cashierUserID,
		"customer_name":   order.CustomerName,
		"total_amount":    order.TotalAmount,
	})

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cashier order created successfully",
		"order":   order,
	})
}
