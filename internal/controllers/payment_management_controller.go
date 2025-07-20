package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "recursiveDine/internal/repositories" // Used in swagger annotations
	"recursiveDine/internal/services"
)

type PaymentManagementController struct {
	paymentService *services.PaymentService
}

func NewPaymentManagementController(paymentService *services.PaymentService) *PaymentManagementController {
	return &PaymentManagementController{
		paymentService: paymentService,
	}
}

// @Summary Get all payments
// @Description Get all payments with pagination (admin/cashier only)
// @Tags payments-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Param method query string false "Filter by payment method"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/payments [get]
func (ctrl *PaymentManagementController) GetAllPayments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	method := c.Query("method")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	payments, total, err := ctrl.paymentService.GetAllPayments(page, limit, status, method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments":    payments,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}

// @Summary Get payment by ID
// @Description Get payment details by ID (admin/cashier only)
// @Tags payments-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Payment ID"
// @Success 200 {object} repositories.Payment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/payments/{id} [get]
func (ctrl *PaymentManagementController) GetPaymentByID(c *gin.Context) {
	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	payment, err := ctrl.paymentService.GetPaymentByIDAdmin(uint(paymentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// @Summary Update payment status
// @Description Update payment status (admin/cashier only)
// @Tags payments-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Payment ID"
// @Param request body map[string]string true "Payment status"
// @Success 200 {object} repositories.Payment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/payments/{id}/status [patch]
func (ctrl *PaymentManagementController) UpdatePaymentStatus(c *gin.Context) {
	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
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
	validStatuses := []string{"pending", "completed", "failed", "refunded", "cancelled"}
	isValid := false
	for _, status := range validStatuses {
		if status == req.Status {
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment status"})
		return
	}

	payment, err := ctrl.paymentService.UpdatePaymentStatus(uint(paymentID), req.Status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// @Summary Process refund
// @Description Process payment refund (admin/cashier only)
// @Tags payments-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Payment ID"
// @Param request body map[string]interface{} true "Refund details"
// @Success 200 {object} repositories.Payment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/payments/{id}/refund [post]
func (ctrl *PaymentManagementController) ProcessRefund(c *gin.Context) {
	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required,min=0.01"`
		Reason string  `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := ctrl.paymentService.ProcessRefund(uint(paymentID), req.Amount, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// @Summary Get payment statistics
// @Description Get payment statistics (admin/cashier only)
// @Tags payments-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Param method query string false "Filter by payment method"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/payments/statistics [get]
func (ctrl *PaymentManagementController) GetPaymentStatistics(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	method := c.Query("method")

	stats, err := ctrl.paymentService.GetPaymentStatistics(from, to, method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// @Summary Get daily revenue
// @Description Get daily revenue by payment method (admin/cashier only)
// @Tags payments-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param from query string false "Start date (YYYY-MM-DD)"
// @Param to query string false "End date (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /admin/payments/revenue [get]
func (ctrl *PaymentManagementController) GetDailyRevenueByPayment(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	revenue, err := ctrl.paymentService.GetDailyRevenueByPayment(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, revenue)
}

// @Summary Reconcile cash payments
// @Description Reconcile cash payments for cashier shift (cashier only)
// @Tags payments-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "Reconciliation data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /cashier/payments/reconcile [post]
func (ctrl *PaymentManagementController) ReconcileCashPayments(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		ActualCashAmount   float64 `json:"actual_cash_amount" binding:"required,min=0"`
		ExpectedCashAmount float64 `json:"expected_cash_amount" binding:"required,min=0"`
		Notes              string  `json:"notes"`
		ShiftStartTime     string  `json:"shift_start_time" binding:"required"`
		ShiftEndTime       string  `json:"shift_end_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := ctrl.paymentService.ReconcileCashPayments(
		userID,
		req.ActualCashAmount,
		req.ExpectedCashAmount,
		req.ShiftStartTime,
		req.ShiftEndTime,
		req.Notes,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Delete payment
// @Description Soft delete payment (admin only)
// @Tags payments-management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Payment ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/payments/{id} [delete]
func (ctrl *PaymentManagementController) DeletePayment(c *gin.Context) {
	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	if err := ctrl.paymentService.DeletePayment(uint(paymentID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}
