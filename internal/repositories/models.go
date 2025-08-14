package repositories

import (
	"time"

	"gorm.io/gorm"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleStaff    UserRole = "staff"
	RoleCashier  UserRole = "cashier"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null;type:varchar(100)"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null;type:text"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;type:text"`
	Phone     string         `json:"phone" gorm:"not null;type:varchar(20)"`
	Password  string         `json:"-" gorm:"not null;type:text"`
	Role      UserRole       `json:"role" gorm:"not null;type:text;default:customer"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Table struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Number      int            `json:"number" gorm:"uniqueIndex;not null"`
	QRCode      string         `json:"qr_code" gorm:"uniqueIndex;not null"`
	Capacity    int            `json:"capacity" gorm:"not null"`
	IsAvailable bool           `json:"is_available" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type MenuCategory struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	MenuItems   []MenuItem     `json:"menu_items,omitempty" gorm:"foreignKey:CategoryID"`
}

type MenuItem struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CategoryID  uint           `json:"category_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	ImageURL    string         `json:"image_url"`
	IsAvailable bool           `json:"is_available" gorm:"default:true"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Category    MenuCategory   `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusServed    OrderStatus = "served"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderType string

const (
	OrderTypeDineIn   OrderType = "dine_in"
	OrderTypeTakeaway OrderType = "takeaway"
)

type Order struct {
	ID                      uint           `json:"id" gorm:"primaryKey"`
	UserID                  uint           `json:"user_id" gorm:"not null"`
	TableID                 uint           `json:"table_id"` // Optional for takeaway orders
	OrderType               OrderType      `json:"order_type" gorm:"not null;default:dine_in"`
	Status                  OrderStatus    `json:"status" gorm:"not null;default:pending"`
	SubtotalAmount          float64        `json:"subtotal_amount" gorm:"not null"`      // Amount before tax
	VATAmount               float64        `json:"vat_amount" gorm:"not null;default:0"` // VAT 10% in Indonesia
	TotalAmount             float64        `json:"total_amount" gorm:"not null"`         // Final amount including VAT
	CustomerName            string         `json:"customer_name" gorm:"type:varchar(255)"`
	CustomerPhone           string         `json:"customer_phone" gorm:"type:varchar(20)"` // For takeaway notifications
	CashierName             string         `json:"cashier_name" gorm:"type:varchar(255)"`
	SpecialNotes            string         `json:"special_notes"`
	EstimatedCompletionTime *time.Time     `json:"estimated_completion_time"` // For takeaway orders
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	DeletedAt               gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User       User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Table      *Table      `json:"table,omitempty" gorm:"foreignKey:TableID"` // Nullable for takeaway
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
	Payment    *Payment    `json:"payment,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	OrderID        uint      `json:"order_id" gorm:"not null"`
	MenuItemID     uint      `json:"menu_item_id" gorm:"not null"`
	Quantity       int       `json:"quantity" gorm:"not null"`
	UnitPrice      float64   `json:"unit_price" gorm:"not null"`
	TotalPrice     float64   `json:"total_price" gorm:"not null"`
	SpecialRequest string    `json:"special_request"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations
	Order    Order    `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	MenuItem MenuItem `json:"menu_item,omitempty" gorm:"foreignKey:MenuItemID"`
}

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type PaymentMethod string

const (
	PaymentMethodQRIS PaymentMethod = "qris"
	PaymentMethodCash PaymentMethod = "cash"
)

type Payment struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	OrderID       uint           `json:"order_id" gorm:"uniqueIndex;not null"`
	Method        PaymentMethod  `json:"method" gorm:"not null"`
	Status        PaymentStatus  `json:"status" gorm:"not null;default:pending"`
	Amount        float64        `json:"amount" gorm:"not null"`
	QRISData      string         `json:"qris_data,omitempty"` // Encrypted QRIS payload
	TransactionID string         `json:"transaction_id,omitempty"`
	ExternalID    string         `json:"external_id,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	Order Order `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}
