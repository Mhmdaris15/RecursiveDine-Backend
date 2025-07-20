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
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=50"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password  string    `json:"-" gorm:"not null" validate:"required,min=6"`
	Role      UserRole  `json:"role" gorm:"not null;default:'customer'" validate:"required,oneof=customer staff admin"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Table struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Number      int       `json:"number" gorm:"uniqueIndex;not null" validate:"required,min=1"`
	QRCode      string    `json:"qr_code" gorm:"uniqueIndex;not null" validate:"required"`
	Capacity    int       `json:"capacity" gorm:"not null" validate:"required,min=1,max=20"`
	IsAvailable bool      `json:"is_available" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type MenuCategory struct {
	ID          uint        `json:"id" gorm:"primaryKey"`
	Name        string      `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Description string      `json:"description"`
	IsActive    bool        `json:"is_active" gorm:"default:true"`
	SortOrder   int         `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	MenuItems   []MenuItem  `json:"menu_items,omitempty" gorm:"foreignKey:CategoryID"`
}

type MenuItem struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	CategoryID  uint         `json:"category_id" gorm:"not null" validate:"required"`
	Name        string       `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Description string       `json:"description"`
	Price       float64      `json:"price" gorm:"not null" validate:"required,gt=0"`
	ImageURL    string       `json:"image_url"`
	IsAvailable bool         `json:"is_available" gorm:"default:true"`
	SortOrder   int          `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Category    MenuCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
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

type Order struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	UserID        uint        `json:"user_id" gorm:"not null" validate:"required"`
	TableID       uint        `json:"table_id" gorm:"not null" validate:"required"`
	Status        OrderStatus `json:"status" gorm:"not null;default:'pending'" validate:"required,oneof=pending confirmed preparing ready served cancelled"`
	TotalAmount   float64     `json:"total_amount" gorm:"not null" validate:"required,gt=0"`
	SpecialNotes  string      `json:"special_notes"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relations
	User       User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Table      Table       `json:"table,omitempty" gorm:"foreignKey:TableID"`
	OrderItems []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
	Payment    *Payment    `json:"payment,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	OrderID        uint      `json:"order_id" gorm:"not null" validate:"required"`
	MenuItemID     uint      `json:"menu_item_id" gorm:"not null" validate:"required"`
	Quantity       int       `json:"quantity" gorm:"not null" validate:"required,min=1"`
	UnitPrice      float64   `json:"unit_price" gorm:"not null" validate:"required,gt=0"`
	TotalPrice     float64   `json:"total_price" gorm:"not null" validate:"required,gt=0"`
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
	ID            uint          `json:"id" gorm:"primaryKey"`
	OrderID       uint          `json:"order_id" gorm:"uniqueIndex;not null" validate:"required"`
	Method        PaymentMethod `json:"method" gorm:"not null" validate:"required,oneof=qris cash"`
	Status        PaymentStatus `json:"status" gorm:"not null;default:'pending'" validate:"required,oneof=pending completed failed refunded"`
	Amount        float64       `json:"amount" gorm:"not null" validate:"required,gt=0"`
	QRISData      string        `json:"qris_data,omitempty"` // Encrypted QRIS payload
	TransactionID string        `json:"transaction_id,omitempty"`
	ExternalID    string        `json:"external_id,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relations
	Order Order `json:"order,omitempty" gorm:"foreignKey:OrderID"`
}
