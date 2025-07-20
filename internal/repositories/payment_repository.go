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

func (r *PaymentRepository) GetAll(limit, offset int) ([]Payment, error) {
	var payments []Payment
	err := r.db.Preload("Order").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) GetAllPaymentsPaginated(offset, limit int, status, method string) ([]*Payment, int64, error) {
	var payments []*Payment
	var total int64
	
	query := r.db.Model(&Payment{})
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if method != "" {
		query = query.Where("payment_method = ?", method)
	}
	
	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	
	// Get payments with preloads
	err = query.Preload("Order").
		Preload("Order.User").
		Preload("Order.Table").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error
	
	return payments, total, err
}

func (r *PaymentRepository) GetByIDWithDetails(id uint) (*Payment, error) {
	var payment Payment
	err := r.db.Preload("Order").
		Preload("Order.User").
		Preload("Order.Table").
		Preload("Order.OrderItems").
		Preload("Order.OrderItems.MenuItem").
		First(&payment, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("payment not found")
	}
	return &payment, err
}

func (r *PaymentRepository) GetPaymentStatistics(from, to string) (map[string]interface{}, error) {
	var stats map[string]interface{}
	// This would be a complex query - for now return empty
	stats = make(map[string]interface{})
	return stats, nil
}

func (r *PaymentRepository) GetDailyRevenueByPayment(from, to string) (map[string]interface{}, error) {
	var revenue map[string]interface{}
	// This would be a complex query - for now return empty
	revenue = make(map[string]interface{})
	return revenue, nil
}

func (r *PaymentRepository) GetCashPaymentsByPeriod(cashierID uint, shiftStart, shiftEnd string) ([]Payment, error) {
	var payments []Payment
	err := r.db.Where("payment_method = ? AND created_at BETWEEN ? AND ?", PaymentMethodCash, shiftStart, shiftEnd).
		Preload("Order").
		Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) CreateReconciliation(reconciliation interface{}) error {
	// Placeholder - would need proper reconciliation model
	return nil
}

func (r *PaymentRepository) SoftDelete(id uint) error {
	return r.db.Model(&Payment{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}
