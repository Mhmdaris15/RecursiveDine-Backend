package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"recursiveDine/internal/config"
	"recursiveDine/internal/repositories"
)

type PaymentService struct {
	paymentRepo *repositories.PaymentRepository
	orderRepo   *repositories.OrderRepository
	config      *config.Config
}

type QRISPaymentRequest struct {
	OrderID uint `json:"order_id" binding:"required"`
}

type QRISPaymentResponse struct {
	PaymentID     uint      `json:"payment_id"`
	QRISData      string    `json:"qris_data"`
	Amount        float64   `json:"amount"`
	ExpiresAt     time.Time `json:"expires_at"`
	TransactionID string    `json:"transaction_id"`
}

type PaymentVerificationRequest struct {
	TransactionID string  `json:"transaction_id" binding:"required"`
	ExternalID    string  `json:"external_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Status        string  `json:"status" binding:"required"`
}

func NewPaymentService(paymentRepo *repositories.PaymentRepository, orderRepo *repositories.OrderRepository, config *config.Config) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		config:      config,
	}
}

func (s *PaymentService) InitiateQRISPayment(req *QRISPaymentRequest) (*QRISPaymentResponse, error) {
	// Get order details
	order, err := s.orderRepo.GetByID(req.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Check if order is in valid status for payment
	if order.Status != repositories.OrderStatusPending {
		return nil, errors.New("order is not pending payment")
	}

	// Check if payment already exists
	if existingPayment, err := s.paymentRepo.GetByOrderID(req.OrderID); err == nil {
		if existingPayment.Status == repositories.PaymentStatusCompleted {
			return nil, errors.New("order already paid")
		}
		// If payment exists but not completed, we can create a new one
	}

	// Generate transaction ID
	transactionID, err := s.generateTransactionID()
	if err != nil {
		return nil, errors.New("failed to generate transaction ID")
	}

	// Generate QRIS data (simplified - in production, integrate with actual QRIS provider)
	qrisData, err := s.generateQRISData(order, transactionID)
	if err != nil {
		return nil, errors.New("failed to generate QRIS data")
	}

	// Create payment record
	payment := &repositories.Payment{
		OrderID:       req.OrderID,
		Method:        repositories.PaymentMethodQRIS,
		Status:        repositories.PaymentStatusPending,
		Amount:        order.TotalAmount,
		QRISData:      qrisData,
		TransactionID: transactionID,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, errors.New("failed to create payment record")
	}

	return &QRISPaymentResponse{
		PaymentID:     payment.ID,
		QRISData:      qrisData,
		Amount:        payment.Amount,
		ExpiresAt:     time.Now().Add(15 * time.Minute), // 15 minutes expiry
		TransactionID: transactionID,
	}, nil
}

func (s *PaymentService) VerifyPayment(req *PaymentVerificationRequest) error {
	// Get payment by transaction ID
	payment, err := s.paymentRepo.GetByTransactionID(req.TransactionID)
	if err != nil {
		return errors.New("payment not found")
	}

	// Verify amount matches
	if payment.Amount != req.Amount {
		return errors.New("amount mismatch")
	}

	// Update payment status based on verification
	var newStatus repositories.PaymentStatus
	switch req.Status {
	case "success", "completed":
		newStatus = repositories.PaymentStatusCompleted
	case "failed":
		newStatus = repositories.PaymentStatusFailed
	default:
		return errors.New("invalid payment status")
	}

	// Update payment
	payment.Status = newStatus
	payment.ExternalID = req.ExternalID
	if err := s.paymentRepo.Update(payment); err != nil {
		return errors.New("failed to update payment status")
	}

	// Update order status if payment is completed
	if newStatus == repositories.PaymentStatusCompleted {
		if err := s.orderRepo.UpdateStatus(payment.OrderID, repositories.OrderStatusConfirmed); err != nil {
			return errors.New("failed to update order status")
		}
	}

	return nil
}

func (s *PaymentService) GetPaymentStatus(paymentID uint) (*repositories.Payment, error) {
	return s.paymentRepo.GetByID(paymentID)
}

func (s *PaymentService) GetPaymentByOrderID(orderID uint) (*repositories.Payment, error) {
	return s.paymentRepo.GetByOrderID(orderID)
}

func (s *PaymentService) GetPaymentByTransactionID(transactionID string) (*repositories.Payment, error) {
	return s.paymentRepo.GetByTransactionID(transactionID)
}

func (s *PaymentService) generateTransactionID() (string, error) {
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	randomString := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("RD%d%s", timestamp, randomString), nil
}

func (s *PaymentService) generateQRISData(order *repositories.Order, transactionID string) (string, error) {
	// Simplified QRIS data generation
	// In production, integrate with actual QRIS provider like GoPay, OVO, etc.
	qrisData := fmt.Sprintf("QRIS:%s:%.2f:%s:%s",
		s.config.QRISMerchantID,
		order.TotalAmount,
		transactionID,
		s.config.QRISCallbackURL,
	)

	// In production, you would encrypt this data for security
	return qrisData, nil
}

func (s *PaymentService) RefundPayment(paymentID uint, reason string) error {
	payment, err := s.paymentRepo.GetByID(paymentID)
	if err != nil {
		return errors.New("payment not found")
	}

	if payment.Status != repositories.PaymentStatusCompleted {
		return errors.New("can only refund completed payments")
	}

	// Update payment status
	payment.Status = repositories.PaymentStatusRefunded
	if err := s.paymentRepo.Update(payment); err != nil {
		return errors.New("failed to update payment status")
	}

	// Update order status
	if err := s.orderRepo.UpdateStatus(payment.OrderID, repositories.OrderStatusCancelled); err != nil {
		return errors.New("failed to update order status")
	}

	return nil
}

func (s *PaymentService) ProcessCashPayment(orderID uint, amountPaid, changeAmount float64) error {
	// Get order details
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	// Check if order is in valid status for payment
	if order.Status != repositories.OrderStatusPending {
		return errors.New("order is not pending payment")
	}

	// Check if payment already exists
	if existingPayment, err := s.paymentRepo.GetByOrderID(orderID); err == nil {
		if existingPayment.Status == repositories.PaymentStatusCompleted {
			return errors.New("order already paid")
		}
	}

	// Validate payment amount
	if amountPaid < order.TotalAmount {
		fmt.Println("Insufficient payment amount")
		return errors.New("insufficient payment amount")
	}

	// Generate transaction ID
	transactionID, err := s.generateTransactionID()
	if err != nil {
		return errors.New("failed to generate transaction ID")
	}

	// Create payment record
	payment := &repositories.Payment{
		OrderID:       orderID,
		Method:        repositories.PaymentMethodCash,
		Status:        repositories.PaymentStatusCompleted,
		Amount:        order.TotalAmount,
		TransactionID: transactionID,
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return errors.New("failed to create payment record")
	}

	// Update order status
	if err := s.orderRepo.UpdateStatus(orderID, repositories.OrderStatusConfirmed); err != nil {
		return errors.New("failed to update order status")
	}

	return nil
}

// Admin/Cashier payment management functions

func (s *PaymentService) GetAllPayments(page, limit int, status, method string) ([]*repositories.Payment, int64, error) {
	offset := (page - 1) * limit
	return s.paymentRepo.GetAllPaymentsPaginated(offset, limit, status, method)
}

func (s *PaymentService) GetPaymentByIDAdmin(id uint) (*repositories.Payment, error) {
	return s.paymentRepo.GetByIDWithDetails(id)
}

func (s *PaymentService) UpdatePaymentStatus(id uint, status string) (*repositories.Payment, error) {
	payment, err := s.paymentRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Convert string status to PaymentStatus enum
	var paymentStatus repositories.PaymentStatus
	switch status {
	case "pending":
		paymentStatus = repositories.PaymentStatusPending
	case "completed":
		paymentStatus = repositories.PaymentStatusCompleted
	case "failed":
		paymentStatus = repositories.PaymentStatusFailed
	case "refunded":
		paymentStatus = repositories.PaymentStatusRefunded
	case "cancelled":
		paymentStatus = repositories.PaymentStatusCancelled
	default:
		return nil, fmt.Errorf("invalid payment status")
	}

	payment.Status = paymentStatus
	if err := s.paymentRepo.Update(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) ProcessRefund(paymentID uint, amount float64, reason string) (*repositories.Payment, error) {
	payment, err := s.paymentRepo.GetByID(paymentID)
	if err != nil {
		return nil, err
	}

	if payment.Status != repositories.PaymentStatusCompleted {
		return nil, fmt.Errorf("can only refund completed payments")
	}

	if amount > payment.Amount {
		return nil, fmt.Errorf("refund amount cannot exceed original payment amount")
	}

	// Create refund record
	refund := &repositories.Payment{
		OrderID:       payment.OrderID,
		Amount:        -amount, // Negative amount for refund
		Method:        payment.Method,
		Status:        repositories.PaymentStatusCompleted,
		TransactionID: fmt.Sprintf("REFUND-%s", payment.TransactionID),
	}

	if err := s.paymentRepo.Create(refund); err != nil {
		return nil, err
	}

	// Update original payment status if full refund
	if amount == payment.Amount {
		payment.Status = repositories.PaymentStatusRefunded
		if err := s.paymentRepo.Update(payment); err != nil {
			return nil, err
		}
	}

	return refund, nil
}

func (s *PaymentService) GetPaymentStatistics(from, to, method string) (map[string]interface{}, error) {
	return s.paymentRepo.GetPaymentStatistics(from, to)
}

func (s *PaymentService) GetDailyRevenueByPayment(from, to string) (map[string]interface{}, error) {
	return s.paymentRepo.GetDailyRevenueByPayment(from, to)
}

func (s *PaymentService) ReconcileCashPayments(cashierID uint, actualAmount, expectedAmount float64, shiftStart, shiftEnd, notes string) (map[string]interface{}, error) {
	// Get cash payments for the shift period
	cashPayments, err := s.paymentRepo.GetCashPaymentsByPeriod(cashierID, shiftStart, shiftEnd)
	if err != nil {
		return nil, err
	}

	var calculatedTotal float64
	for _, payment := range cashPayments {
		if payment.Amount > 0 { // Exclude refunds
			calculatedTotal += payment.Amount
		}
	}

	difference := actualAmount - calculatedTotal

	// Create reconciliation record
	reconciliation := map[string]interface{}{
		"cashier_id":          cashierID,
		"shift_start":         shiftStart,
		"shift_end":           shiftEnd,
		"expected_amount":     expectedAmount,
		"calculated_amount":   calculatedTotal,
		"actual_amount":       actualAmount,
		"difference":          difference,
		"payment_count":       len(cashPayments),
		"notes":               notes,
		"reconciliation_time": time.Now(),
	}

	// Store reconciliation record
	if err := s.paymentRepo.CreateReconciliation(reconciliation); err != nil {
		return nil, err
	}

	return reconciliation, nil
}

func (s *PaymentService) DeletePayment(id uint) error {
	return s.paymentRepo.SoftDelete(id)
}
