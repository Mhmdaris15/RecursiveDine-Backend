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

func (r *OrderRepository) GetByOrderType(orderType OrderType, limit, offset int) ([]Order, error) {
	var orders []Order
	err := r.db.Where("order_type = ?", orderType).
		Preload("User").
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

func (r *OrderRepository) GetByStatusAndType(status OrderStatus, orderType OrderType, limit, offset int) ([]Order, error) {
	var orders []Order
	err := r.db.Where("status = ? AND order_type = ?", status, orderType).
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

func (r *OrderRepository) GetTakeawayOrdersReady() ([]Order, error) {
	var orders []Order
	err := r.db.Where("order_type = ? AND status = ?", OrderTypeTakeaway, OrderStatusReady).
		Preload("User").
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Preload("Payment").
		Order("estimated_completion_time ASC").
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

func (r *OrderRepository) GetAllOrdersPaginated(limit, offset int, status string) ([]*Order, int64, error) {
	var orders []*Order
	var total int64

	query := r.db.Model(&Order{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get orders with preloads
	err = query.Preload("User").
		Preload("Table").
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Preload("Payment").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error

	return orders, total, err
}

func (r *OrderRepository) GetByIDWithDetails(id uint) (*Order, error) {
	var order Order
	err := r.db.Preload("User").
		Preload("Table").
		Preload("OrderItems").
		Preload("OrderItems.MenuItem").
		Preload("OrderItems.MenuItem.Category").
		Preload("Payment").
		First(&order, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("order not found")
	}
	return &order, err
}

func (r *OrderRepository) SoftDelete(id uint) error {
	return r.db.Model(&Order{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *OrderRepository) GetOrderStatistics(from, to string) (map[string]interface{}, error) {
	var stats map[string]interface{}
	// This would be a complex query - for now return empty
	stats = make(map[string]interface{})
	return stats, nil
}

func (r *OrderRepository) GetDailyRevenue(from, to string) (map[string]interface{}, error) {
	var revenue map[string]interface{}
	// This would be a complex query - for now return empty
	revenue = make(map[string]interface{})
	return revenue, nil
}

func (r *OrderRepository) UpdateOrderItems(orderID uint, items []OrderItem) error {
	// Delete existing items
	err := r.db.Where("order_id = ?", orderID).Delete(&OrderItem{}).Error
	if err != nil {
		return err
	}

	// Create new items
	for i := range items {
		items[i].OrderID = orderID
		err = r.db.Create(&items[i]).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *OrderRepository) UpdateOrder(order *Order) error {
	return r.db.Save(order).Error
}
