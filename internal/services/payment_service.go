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
	PaymentID     uint   `json:"payment_id"`
	QRISData      string `json:"qris_data"`
	Amount        float64 `json:"amount"`
	ExpiresAt     time.Time `json:"expires_at"`
	TransactionID string `json:"transaction_id"`
}

type PaymentVerificationRequest struct {
	TransactionID string `json:"transaction_id" binding:"required"`
	ExternalID    string `json:"external_id" binding:"required"`
	Amount        float64 `json:"amount" binding:"required"`
	Status        string `json:"status" binding:"required"`
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
