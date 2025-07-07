package repositories

import (
	"errors"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(payment *Payment) error {
	return r.db.Create(payment).Error
}

func (r *PaymentRepository) GetByID(id uint) (*Payment, error) {
	var payment Payment
	err := r.db.Preload("Order").First(&payment, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("payment not found")
	}
	return &payment, err
}

func (r *PaymentRepository) GetByOrderID(orderID uint) (*Payment, error) {
	var payment Payment
	err := r.db.Where("order_id = ?", orderID).Preload("Order").First(&payment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("payment not found")
	}
	return &payment, err
}

func (r *PaymentRepository) GetByTransactionID(transactionID string) (*Payment, error) {
	var payment Payment
	err := r.db.Where("transaction_id = ?", transactionID).Preload("Order").First(&payment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("payment not found")
	}
	return &payment, err
}

func (r *PaymentRepository) GetByExternalID(externalID string) (*Payment, error) {
	var payment Payment
	err := r.db.Where("external_id = ?", externalID).Preload("Order").First(&payment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("payment not found")
	}
	return &payment, err
}

func (r *PaymentRepository) GetByStatus(status PaymentStatus, limit, offset int) ([]Payment, error) {
	var payments []Payment
	err := r.db.Where("status = ?", status).
		Preload("Order").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) Update(payment *Payment) error {
	return r.db.Save(payment).Error
}

func (r *PaymentRepository) UpdateStatus(paymentID uint, status PaymentStatus) error {
	return r.db.Model(&Payment{}).Where("id = ?", paymentID).Update("status", status).Error
}

func (r *PaymentRepository) Delete(id uint) error {
	return r.db.Delete(&Payment{}, id).Error
}

func (r *PaymentRepository) GetPendingPayments() ([]Payment, error) {
	var payments []Payment
	err := r.db.Where("status = ?", PaymentStatusPending).
		Preload("Order").
		Order("created_at ASC").
		Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) IsTransactionIDExists(transactionID string) (bool, error) {
	var count int64
	err := r.db.Model(&Payment{}).Where("transaction_id = ?", transactionID).Count(&count).Error
	return count > 0, err
}

func (r *PaymentRepository) IsExternalIDExists(externalID string) (bool, error) {
	var count int64
	err := r.db.Model(&Payment{}).Where("external_id = ?", externalID).Count(&count).Error
	return count > 0, err
}
