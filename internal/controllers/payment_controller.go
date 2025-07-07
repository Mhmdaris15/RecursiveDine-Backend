package controllers

import (
	"net/http"

	"recursiveDine/internal/services"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentService *services.PaymentService
}

func NewPaymentController(paymentService *services.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

// @Summary Initiate QRIS payment
// @Description Create a QRIS payment for an order
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.QRISPaymentRequest true "Payment request"
// @Success 201 {object} services.QRISPaymentResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/qris [post]
func (ctrl *PaymentController) InitiateQRISPayment(c *gin.Context) {
	var req services.QRISPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.paymentService.InitiateQRISPayment(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Verify payment
// @Description Verify payment completion from payment provider
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.PaymentVerificationRequest true "Payment verification data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /payments/verify [post]
func (ctrl *PaymentController) VerifyPayment(c *gin.Context) {
	var req services.PaymentVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.paymentService.VerifyPayment(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment verified successfully"})
}

// @Summary Get payment status
// @Description Get the current status of a payment
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payment_id path int true "Payment ID"
// @Success 200 {object} repositories.Payment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /payments/status/{payment_id} [get]
func (ctrl *PaymentController) GetPaymentStatus(c *gin.Context) {
	var req struct {
		PaymentID uint `uri:"payment_id" binding:"required"`
	}

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := ctrl.paymentService.GetPaymentStatus(req.PaymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// @Summary Get payment by order ID
// @Description Get payment information for a specific order
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_id query int true "Order ID"
// @Success 200 {object} repositories.Payment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /payments/order [get]
func (ctrl *PaymentController) GetPaymentByOrderID(c *gin.Context) {
	var req struct {
		OrderID uint `form:"order_id" binding:"required"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := ctrl.paymentService.GetPaymentByOrderID(req.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// @Summary Payment webhook
// @Description Handle payment status updates from payment provider
// @Tags payments
// @Accept json
// @Produce json
// @Param request body services.PaymentVerificationRequest true "Payment webhook data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /payments/webhook [post]
func (ctrl *PaymentController) PaymentWebhook(c *gin.Context) {
	var req services.PaymentVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In production, you should verify the webhook signature
	// to ensure it's coming from the legitimate payment provider

	if err := ctrl.paymentService.VerifyPayment(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}
