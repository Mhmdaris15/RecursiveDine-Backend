package services

import (
	"errors"
	"fmt"

	"recursiveDine/internal/repositories"
)

type OrderService struct {
	orderRepo *repositories.OrderRepository
	menuRepo  *repositories.MenuRepository
}

type CreateOrderRequest struct {
	TableID      uint                     `json:"table_id" binding:"required"`
	SpecialNotes string                   `json:"special_notes"`
	Items        []CreateOrderItemRequest `json:"items" binding:"required,dive"`
}

type CreateOrderItemRequest struct {
	MenuItemID     uint   `json:"menu_item_id" binding:"required"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	SpecialRequest string `json:"special_request"`
}

func NewOrderService(orderRepo *repositories.OrderRepository, menuRepo *repositories.MenuRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		menuRepo:  menuRepo,
	}
}

func (s *OrderService) CreateOrder(userID uint, req *CreateOrderRequest) (*repositories.Order, error) {
	// Validate menu items
	menuItemIDs := make([]uint, len(req.Items))
	for i, item := range req.Items {
		menuItemIDs[i] = item.MenuItemID
	}

	menuItems, err := s.menuRepo.GetMenuItemsByIDs(menuItemIDs)
	if err != nil {
		return nil, errors.New("failed to fetch menu items")
	}

	if len(menuItems) != len(req.Items) {
		return nil, errors.New("some menu items are not available")
	}

	// Create menu items map for easy lookup
	menuItemMap := make(map[uint]*repositories.MenuItem)
	for i := range menuItems {
		menuItemMap[menuItems[i].ID] = &menuItems[i]
	}

	// Calculate total amount and create order items
	var totalAmount float64
	orderItems := make([]repositories.OrderItem, 0, len(req.Items))

	for _, item := range req.Items {
		menuItem, exists := menuItemMap[item.MenuItemID]
		if !exists {
			return nil, fmt.Errorf("menu item with ID %d not found", item.MenuItemID)
		}

		if !menuItem.IsAvailable {
			return nil, fmt.Errorf("menu item '%s' is not available", menuItem.Name)
		}

		totalPrice := menuItem.Price * float64(item.Quantity)
		totalAmount += totalPrice

		orderItem := repositories.OrderItem{
			MenuItemID:     item.MenuItemID,
			Quantity:       item.Quantity,
			UnitPrice:      menuItem.Price,
			TotalPrice:     totalPrice,
			SpecialRequest: item.SpecialRequest,
		}

		orderItems = append(orderItems, orderItem)
	}

	// Create order
	order := &repositories.Order{
		UserID:       userID,
		TableID:      req.TableID,
		Status:       repositories.OrderStatusPending,
		TotalAmount:  totalAmount,
		SpecialNotes: req.SpecialNotes,
		OrderItems:   orderItems,
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, errors.New("failed to create order")
	}

	// Return order with all relations
	return s.orderRepo.GetByID(order.ID)
}

func (s *OrderService) GetOrderByID(orderID uint) (*repositories.Order, error) {
	return s.orderRepo.GetByID(orderID)
}

func (s *OrderService) GetOrdersByUser(userID uint, page, limit int) ([]repositories.Order, error) {
	offset := (page - 1) * limit
	return s.orderRepo.GetByUserID(userID, limit, offset)
}

func (s *OrderService) GetOrdersByTable(tableID uint, page, limit int) ([]repositories.Order, error) {
	offset := (page - 1) * limit
	return s.orderRepo.GetByTableID(tableID, limit, offset)
}

func (s *OrderService) GetOrdersByStatus(status repositories.OrderStatus, page, limit int) ([]repositories.Order, error) {
	offset := (page - 1) * limit
	return s.orderRepo.GetByStatus(status, limit, offset)
}

func (s *OrderService) GetAllOrders(page, limit int) ([]repositories.Order, error) {
	offset := (page - 1) * limit
	return s.orderRepo.GetAll(limit, offset)
}

func (s *OrderService) UpdateOrderStatus(orderID uint, status repositories.OrderStatus) error {
	// Validate order exists
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return errors.New("order not found")
	}

	// Validate status transition
	if err := s.validateStatusTransition(order.Status, status); err != nil {
		return err
	}

	return s.orderRepo.UpdateStatus(orderID, status)
}

func (s *OrderService) GetActiveOrders() ([]repositories.Order, error) {
	return s.orderRepo.GetActiveOrders()
}

func (s *OrderService) GetKitchenOrders() ([]repositories.Order, error) {
	return s.orderRepo.GetKitchenOrders()
}

func (s *OrderService) validateStatusTransition(currentStatus, newStatus repositories.OrderStatus) error {
	validTransitions := map[repositories.OrderStatus][]repositories.OrderStatus{
		repositories.OrderStatusPending: {
			repositories.OrderStatusConfirmed,
			repositories.OrderStatusCancelled,
		},
		repositories.OrderStatusConfirmed: {
			repositories.OrderStatusPreparing,
			repositories.OrderStatusCancelled,
		},
		repositories.OrderStatusPreparing: {
			repositories.OrderStatusReady,
		},
		repositories.OrderStatusReady: {
			repositories.OrderStatusServed,
		},
		repositories.OrderStatusServed: {
			// Final status
		},
		repositories.OrderStatusCancelled: {
			// Final status
		},
	}

	validNextStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return errors.New("invalid current status")
	}

	for _, validStatus := range validNextStatuses {
		if validStatus == newStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
}

// Admin/Staff order management functions

func (s *OrderService) GetAllOrdersAdmin(page, limit int, status string) ([]*repositories.Order, int64, error) {
	offset := (page - 1) * limit
	return s.orderRepo.GetAllOrdersPaginated(offset, limit, status)
}

func (s *OrderService) GetOrderByIDAdmin(id uint) (*repositories.Order, error) {
	return s.orderRepo.GetByIDWithDetails(id)
}

func (s *OrderService) UpdateOrderStatusAdmin(id uint, status string) (*repositories.Order, error) {
	// Convert string status to OrderStatus enum
	var orderStatus repositories.OrderStatus
	switch status {
	case "pending":
		orderStatus = repositories.OrderStatusPending
	case "confirmed":
		orderStatus = repositories.OrderStatusConfirmed
	case "preparing":
		orderStatus = repositories.OrderStatusPreparing
	case "ready":
		orderStatus = repositories.OrderStatusReady
	case "served":
		orderStatus = repositories.OrderStatusServed
	case "cancelled":
		orderStatus = repositories.OrderStatusCancelled
	default:
		return nil, errors.New("invalid order status")
	}

	if err := s.orderRepo.UpdateStatus(id, orderStatus); err != nil {
		return nil, err
	}

	return s.orderRepo.GetByID(id)
}

func (s *OrderService) DeleteOrder(id uint) error {
	return s.orderRepo.SoftDelete(id)
}

func (s *OrderService) GetOrderStatistics(from, to string) (map[string]interface{}, error) {
	return s.orderRepo.GetOrderStatistics(from, to)
}

func (s *OrderService) GetDailyRevenue(from, to string) (map[string]interface{}, error) {
	return s.orderRepo.GetDailyRevenue(from, to)
}

func (s *OrderService) UpdateOrderItems(orderID uint, items []repositories.OrderItem) (*repositories.Order, error) {
	// Check if order is in pending status
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return nil, err
	}

	if order.Status != repositories.OrderStatusPending {
		return nil, fmt.Errorf("can only update items for pending orders")
	}

	// Calculate new total
	var total float64
	for i, item := range items {
		menuItem, err := s.menuRepo.GetMenuItemByID(item.MenuItemID)
		if err != nil {
			return nil, fmt.Errorf("menu item not found: %d", item.MenuItemID)
		}
		total += menuItem.Price * float64(item.Quantity)
		items[i].UnitPrice = menuItem.Price
		items[i].TotalPrice = menuItem.Price * float64(item.Quantity)
		items[i].OrderID = orderID
	}

	// Update order items and total
	if err := s.orderRepo.UpdateOrderItems(orderID, items); err != nil {
		return nil, err
	}

	// Update total amount
	order.TotalAmount = total
	if err := s.orderRepo.UpdateOrder(order); err != nil {
		return nil, err
	}

	return s.orderRepo.GetByIDWithDetails(orderID)
}
