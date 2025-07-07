package repositories

import (
	"errors"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(order *Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepository) GetByID(id uint) (*Order, error) {
	var order Order
	err := r.db.Preload("User").
		Preload("Table").
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Preload("Payment").
		First(&order, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("order not found")
	}
	return &order, err
}

func (r *OrderRepository) GetByUserID(userID uint, limit, offset int) ([]Order, error) {
	var orders []Order
	err := r.db.Where("user_id = ?", userID).
		Preload("Table").
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Preload("Payment").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) GetByTableID(tableID uint, limit, offset int) ([]Order, error) {
	var orders []Order
	err := r.db.Where("table_id = ?", tableID).
		Preload("User").
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Preload("Payment").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) GetByStatus(status OrderStatus, limit, offset int) ([]Order, error) {
	var orders []Order
	err := r.db.Where("status = ?", status).
		Preload("User").
		Preload("Table").
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Preload("Payment").
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) GetAll(limit, offset int) ([]Order, error) {
	var orders []Order
	err := r.db.Preload("User").
		Preload("Table").
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Preload("Payment").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) Update(order *Order) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) UpdateStatus(orderID uint, status OrderStatus) error {
	return r.db.Model(&Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *OrderRepository) Delete(id uint) error {
	return r.db.Delete(&Order{}, id).Error
}

func (r *OrderRepository) GetActiveOrders() ([]Order, error) {
	var orders []Order
	err := r.db.Where("status IN ?", []OrderStatus{OrderStatusPending, OrderStatusConfirmed, OrderStatusPreparing}).
		Preload("User").
		Preload("Table").
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Order("created_at ASC").
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) GetKitchenOrders() ([]Order, error) {
	var orders []Order
	err := r.db.Where("status IN ?", []OrderStatus{OrderStatusConfirmed, OrderStatusPreparing}).
		Preload("Table").
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Order("created_at ASC").
		Find(&orders).Error
	return orders, err
}

// OrderItem operations
func (r *OrderRepository) CreateOrderItem(item *OrderItem) error {
	return r.db.Create(item).Error
}

func (r *OrderRepository) GetOrderItemsByOrderID(orderID uint) ([]OrderItem, error) {
	var items []OrderItem
	err := r.db.Where("order_id = ?", orderID).
		Preload("MenuItem").
		Find(&items).Error
	return items, err
}

func (r *OrderRepository) UpdateOrderItem(item *OrderItem) error {
	return r.db.Save(item).Error
}

func (r *OrderRepository) DeleteOrderItem(id uint) error {
	return r.db.Delete(&OrderItem{}, id).Error
}
