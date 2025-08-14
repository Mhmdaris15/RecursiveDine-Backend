package services

import (
	"errors"
	"fmt"
	"time"

	"recursiveDine/internal/repositories"
)

type OrderService struct {
	orderRepo *repositories.OrderRepository
	menuRepo  *repositories.MenuRepository
}

type CreateOrderRequest struct {
	TableID                 *uint                    `json:"table_id"` // Optional for takeaway
	OrderType               repositories.OrderType   `json:"order_type" binding:"required"`
	CustomerPhone           string                   `json:"customer_phone"` // Required for takeaway
	SpecialNotes            string                   `json:"special_notes"`
	EstimatedCompletionTime *time.Time               `json:"estimated_completion_time"` // For takeaway orders
	Items                   []CreateOrderItemRequest `json:"items" binding:"required,dive"`
}

type CreateOrderItemRequest struct {
	MenuItemID     uint   `json:"menu_item_id" binding:"required"`
	Quantity       int    `json:"quantity" binding:"required,min=1"`
	SpecialRequest string `json:"special_request"`
}

type CashierOrderRequest struct {
	TableID                 *uint                    `json:"table_id"` // Optional for takeaway
	OrderType               repositories.OrderType   `json:"order_type" binding:"required"`
	CustomerName            string                   `json:"customer_name" binding:"required"`
	CustomerPhone           string                   `json:"customer_phone"` // Required for takeaway
	CashierName             string                   `json:"cashier_name" binding:"required"`
	SpecialNotes            string                   `json:"special_notes"`
	EstimatedCompletionTime *time.Time               `json:"estimated_completion_time"` // For takeaway orders
	Items                   []CreateOrderItemRequest `json:"items" binding:"required,dive"`
}

type OrderResponse struct {
	ID                      uint                     `json:"id"`
	UserID                  uint                     `json:"user_id"`
	TableID                 uint                     `json:"table_id,omitempty"` // Omit if null for takeaway
	OrderType               repositories.OrderType   `json:"order_type"`
	Status                  repositories.OrderStatus `json:"status"`
	SubtotalAmount          float64                  `json:"subtotal_amount"`
	VATAmount               float64                  `json:"vat_amount"`
	TotalAmount             float64                  `json:"total_amount"`
	CustomerName            string                   `json:"customer_name,omitempty"`
	CustomerPhone           string                   `json:"customer_phone,omitempty"`
	CashierName             string                   `json:"cashier_name,omitempty"`
	SpecialNotes            string                   `json:"special_notes"`
	EstimatedCompletionTime *string                  `json:"estimated_completion_time,omitempty"`
	CreatedAt               string                   `json:"created_at"`
	OrderItems              []OrderItemResponse      `json:"order_items"`
}

type OrderItemResponse struct {
	ID             uint    `json:"id"`
	MenuItemID     uint    `json:"menu_item_id"`
	MenuItemName   string  `json:"menu_item_name"`
	Quantity       int     `json:"quantity"`
	UnitPrice      float64 `json:"unit_price"`
	TotalPrice     float64 `json:"total_price"`
	SpecialRequest string  `json:"special_request"`
}

func NewOrderService(orderRepo *repositories.OrderRepository, menuRepo *repositories.MenuRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
		menuRepo:  menuRepo,
	}
}

func (s *OrderService) CreateOrder(userID uint, req *CreateOrderRequest) (*repositories.Order, error) {
	// Validate order type and required fields
	if err := s.validateOrderRequest(req); err != nil {
		return nil, err
	}

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
	var subtotal float64
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
		subtotal += totalPrice

		orderItem := repositories.OrderItem{
			MenuItemID:     item.MenuItemID,
			Quantity:       item.Quantity,
			UnitPrice:      menuItem.Price,
			TotalPrice:     totalPrice,
			SpecialRequest: item.SpecialRequest,
		}

		orderItems = append(orderItems, orderItem)
	}

	// Calculate VAT and total
	const vatRate = 0.10
	vatAmount := subtotal * vatRate
	totalAmount := subtotal + vatAmount

	// Create order
	order := &repositories.Order{
		UserID:                  userID,
		OrderType:               req.OrderType,
		Status:                  repositories.OrderStatusPending,
		SubtotalAmount:          subtotal,
		VATAmount:               vatAmount,
		TotalAmount:             totalAmount,
		CustomerPhone:           req.CustomerPhone,
		SpecialNotes:            req.SpecialNotes,
		EstimatedCompletionTime: req.EstimatedCompletionTime,
		OrderItems:              orderItems,
	}

	// Set TableID only if provided (for dine-in orders)
	if req.TableID != nil {
		order.TableID = *req.TableID
	}

	if err := s.orderRepo.Create(order); err != nil {
		return nil, errors.New("failed to create order")
	}

	// Return order with all relations
	return s.orderRepo.GetByID(order.ID)
}

func (s *OrderService) validateOrderRequest(req *CreateOrderRequest) error {
	switch req.OrderType {
	case repositories.OrderTypeDineIn:
		if req.TableID == nil {
			return errors.New("table_id is required for dine-in orders")
		}
	case repositories.OrderTypeTakeaway:
		if req.CustomerPhone == "" {
			return errors.New("customer_phone is required for takeaway orders")
		}
		// Set estimated completion time if not provided (default 30 minutes)
		if req.EstimatedCompletionTime == nil {
			estimatedTime := time.Now().Add(30 * time.Minute)
			req.EstimatedCompletionTime = &estimatedTime
		}
	default:
		return errors.New("invalid order type. Must be 'dine_in' or 'takeaway'")
	}
	return nil
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
	return s.orderRepo.GetAllOrdersPaginated(limit, offset, status)
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

// CreateCashierOrder creates an order through cashier with VAT calculation
func (s *OrderService) CreateCashierOrder(cashierUserID uint, req *CashierOrderRequest) (*OrderResponse, error) {
	// Validate order type and required fields
	if err := s.validateCashierOrderRequest(req); err != nil {
		return nil, err
	}

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
	for _, item := range menuItems {
		menuItemMap[item.ID] = &item
	}

	// Calculate subtotal
	var subtotal float64
	for _, orderItem := range req.Items {
		menuItem := menuItemMap[orderItem.MenuItemID]
		if !menuItem.IsAvailable {
			return nil, fmt.Errorf("menu item %s is not available", menuItem.Name)
		}
		subtotal += menuItem.Price * float64(orderItem.Quantity)
	}

	// Calculate VAT (10% for Indonesia)
	const vatRate = 0.10
	vatAmount := subtotal * vatRate
	totalAmount := subtotal + vatAmount

	createdOrder := &repositories.Order{
		UserID:                  cashierUserID,
		OrderType:               req.OrderType,
		Status:                  repositories.OrderStatusPending,
		SubtotalAmount:          subtotal,
		VATAmount:               vatAmount,
		TotalAmount:             totalAmount,
		CustomerName:            req.CustomerName,
		CustomerPhone:           req.CustomerPhone,
		CashierName:             req.CashierName,
		SpecialNotes:            req.SpecialNotes,
		EstimatedCompletionTime: req.EstimatedCompletionTime,
	}

	// Set TableID only if provided (for dine-in orders)
	if req.TableID != nil {
		createdOrder.TableID = *req.TableID
	}

	err = s.orderRepo.Create(createdOrder)
	if err != nil {
		return nil, errors.New("failed to create order")
	}

	// Create order items
	for _, item := range req.Items {
		menuItem := menuItemMap[item.MenuItemID]
		totalPrice := menuItem.Price * float64(item.Quantity)

		orderItem := &repositories.OrderItem{
			OrderID:        createdOrder.ID,
			MenuItemID:     item.MenuItemID,
			Quantity:       item.Quantity,
			UnitPrice:      menuItem.Price,
			TotalPrice:     totalPrice,
			SpecialRequest: item.SpecialRequest,
		}

		err = s.orderRepo.CreateOrderItem(orderItem)
		if err != nil {
			return nil, errors.New("failed to create order items")
		}
	}

	// Fetch the complete order with items for response
	completeOrder, err := s.orderRepo.GetByID(createdOrder.ID)
	if err != nil {
		return nil, errors.New("failed to fetch created order")
	}

	// Build response
	orderItemResponses := make([]OrderItemResponse, len(completeOrder.OrderItems))
	for i, item := range completeOrder.OrderItems {
		orderItemResponses[i] = OrderItemResponse{
			ID:             item.ID,
			MenuItemID:     item.MenuItemID,
			MenuItemName:   item.MenuItem.Name,
			Quantity:       item.Quantity,
			UnitPrice:      item.UnitPrice,
			TotalPrice:     item.TotalPrice,
			SpecialRequest: item.SpecialRequest,
		}
	}

	response := &OrderResponse{
		ID:             completeOrder.ID,
		UserID:         completeOrder.UserID,
		TableID:        completeOrder.TableID,
		OrderType:      completeOrder.OrderType,
		Status:         completeOrder.Status,
		SubtotalAmount: completeOrder.SubtotalAmount,
		VATAmount:      completeOrder.VATAmount,
		TotalAmount:    completeOrder.TotalAmount,
		CustomerName:   completeOrder.CustomerName,
		CustomerPhone:  completeOrder.CustomerPhone,
		CashierName:    completeOrder.CashierName,
		SpecialNotes:   completeOrder.SpecialNotes,
		CreatedAt:      completeOrder.CreatedAt.Format("2006-01-02 15:04:05"),
		OrderItems:     orderItemResponses,
	}

	// Add estimated completion time if present
	if completeOrder.EstimatedCompletionTime != nil {
		estimatedTime := completeOrder.EstimatedCompletionTime.Format("2006-01-02 15:04:05")
		response.EstimatedCompletionTime = &estimatedTime
	}

	return response, nil
}

func (s *OrderService) validateCashierOrderRequest(req *CashierOrderRequest) error {
	switch req.OrderType {
	case repositories.OrderTypeDineIn:
		if req.TableID == nil {
			return errors.New("table_id is required for dine-in orders")
		}
	case repositories.OrderTypeTakeaway:
		if req.CustomerPhone == "" {
			return errors.New("customer_phone is required for takeaway orders")
		}
		// Set estimated completion time if not provided (default 30 minutes)
		if req.EstimatedCompletionTime == nil {
			estimatedTime := time.Now().Add(30 * time.Minute)
			req.EstimatedCompletionTime = &estimatedTime
		}
	default:
		return errors.New("invalid order type. Must be 'dine_in' or 'takeaway'")
	}
	return nil
}

// GetOrdersByType returns orders filtered by order type
func (s *OrderService) GetOrdersByType(orderType repositories.OrderType, page, limit int) ([]repositories.Order, error) {
	offset := (page - 1) * limit
	return s.orderRepo.GetByOrderType(orderType, limit, offset)
}

// GetTakeawayOrdersReady returns all ready takeaway orders
func (s *OrderService) GetTakeawayOrdersReady() ([]repositories.Order, error) {
	return s.orderRepo.GetTakeawayOrdersReady()
}

// GetOrdersByStatusAndType returns orders filtered by both status and type
func (s *OrderService) GetOrdersByStatusAndType(status repositories.OrderStatus, orderType repositories.OrderType, page, limit int) ([]repositories.Order, error) {
	offset := (page - 1) * limit
	return s.orderRepo.GetByStatusAndType(status, orderType, limit, offset)
}
