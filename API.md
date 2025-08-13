# RecursiveDine API Documentation

## Overview
RecursiveDine is a secure, scalable restaurant self-service ordering system built with Go and PostgreSQL. This API provides comprehensive endpoints for managing restaurants, orders, payments, and user roles with Indonesian VAT calculation support (10% tax rate).

## Base URL
```
http://localhost:8002/api/v1
```

## Key Features
- **Multi-role Authentication**: Customer, Staff, Cashier, and Admin roles
- **Advanced User Management**: Comprehensive admin controls with filtering, search, statistics, and bulk operations
- **VAT Calculation**: Automatic 10% Indonesian VAT on cashier orders
- **Real-time Updates**: WebSocket support for kitchen operations
- **Payment Processing**: QRIS and cash payment methods
- **Order Management**: Complete order lifecycle tracking
- **Table Management**: QR code-based table system
- **Rate Limiting**: Built-in API protection
- **Comprehensive Logging**: Request/response monitoring
- **Security Features**: Role-based access control, password hashing, account protection

## Authentication
Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## User Roles
- **Customer**: Can place orders and view their own orders
- **Staff**: Can manage orders and update menu availability
- **Cashier**: Can process payments and handle cash transactions
- **Admin**: Full access to all system features

---

## 1. Authentication Endpoints

### POST /auth/register
Register a new user account.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123",
  "phone": "+1234567890"
}
```

**Response (201):**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "role": "customer",
    "created_at": "2025-07-20T10:00:00Z"
  }
}
```

### POST /auth/login
Login with email or username and password.

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Alternative (Username):**
```json
{
  "username": "john_doe",
  "password": "securepassword123"
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600,
  "user": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "role": "customer"
  }
}
```

### POST /auth/refresh
Refresh access token using refresh token.

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 3600
}
```

### POST /auth/logout
Logout and invalidate tokens. Requires authentication.

**Response (200):**
```json
{
  "message": "Logged out successfully"
}
```

---

## 2. Table Management

### GET /tables/{qr_code}
Get table information by QR code (public endpoint).

**Response (200):**
```json
{
  "id": 1,
  "number": "T001",
  "capacity": 4,
  "location": "Main Hall",
  "qr_code": "QR123456",
  "status": "available"
}
```

### GET /admin/tables
Get all tables (Admin only).

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10)
- `status`: Filter by status (available, occupied, reserved, maintenance)

**Response (200):**
```json
{
  "tables": [
    {
      "id": 1,
      "number": "T001",
      "capacity": 4,
      "location": "Main Hall",
      "qr_code": "QR123456",
      "status": "available",
      "created_at": "2025-07-20T10:00:00Z"
    }
  ],
  "total": 25,
  "page": 1,
  "limit": 10,
  "total_pages": 3
}
```

### POST /admin/tables
Create a new table (Admin only).

**Request Body:**
```json
{
  "number": "T025",
  "capacity": 6,
  "location": "VIP Section"
}
```

**Response (201):**
```json
{
  "id": 25,
  "number": "T025",
  "capacity": 6,
  "location": "VIP Section",
  "qr_code": "QR789012",
  "status": "available",
  "created_at": "2025-07-20T10:00:00Z"
}
```

### PUT /admin/tables/{id}
Update table information (Admin only).

**Request Body:**
```json
{
  "number": "T025-VIP",
  "capacity": 8,
  "location": "VIP Section Premium"
}
```

### PATCH /admin/tables/{id}/status
Update table status (Admin only).

**Request Body:**
```json
{
  "status": "maintenance"
}
```

### DELETE /admin/tables/{id}
Delete a table (Admin only).

**Response (200):**
```json
{
  "message": "Table deleted successfully"
}
```

---

## 3. Menu Management

### GET /menu
Get available menu with categories and items (public endpoint).

**Response (200):**
```json
{
  "categories": [
    {
      "id": 1,
      "name": "Appetizers",
      "description": "Start your meal right",
      "items": [
        {
          "id": 1,
          "name": "Caesar Salad",
          "description": "Fresh romaine lettuce with caesar dressing",
          "price": 12.99,
          "category_id": 1,
          "image_url": "/images/caesar-salad.jpg",
          "is_available": true,
          "preparation_time": 10
        }
      ]
    }
  ]
}
```

### GET /menu/categories
Get all menu categories (public endpoint).

**Response (200):**
```json
{
  "categories": [
    {
      "id": 1,
      "name": "Appetizers",
      "description": "Start your meal right"
    }
  ]
}
```

### GET /menu/items/search
Search menu items (public endpoint).

**Query Parameters:**
- `q`: Search query
- `category_id`: Filter by category
- `min_price`: Minimum price
- `max_price`: Maximum price

**Response (200):**
```json
{
  "items": [
    {
      "id": 1,
      "name": "Caesar Salad",
      "description": "Fresh romaine lettuce with caesar dressing",
      "price": 12.99,
      "category_id": 1,
      "is_available": true
    }
  ]
}
```

### POST /admin/menu/categories
Create a new menu category (Admin only).

**Request Body:**
```json
{
  "name": "Desserts",
  "description": "Sweet endings to your meal"
}
```

### POST /admin/menu/items
Create a new menu item (Admin only).

**Request Body:**
```json
{
  "name": "Chocolate Cake",
  "description": "Rich chocolate cake with ganache",
  "price": 8.99,
  "category_id": 2,
  "image_url": "/images/chocolate-cake.jpg",
  "preparation_time": 15,
  "is_available": true
}
```

### PATCH /admin/menu/items/{id}/availability
Update menu item availability (Admin/Staff).

**Request Body:**
```json
{
  "is_available": false
}
```

---

## 4. Order Management

### POST /orders
Create a new order (Authenticated users).

**Request Body:**
```json
{
  "table_id": 1,
  "special_notes": "No onions please",
  "items": [
    {
      "menu_item_id": 1,
      "quantity": 2,
      "special_request": "Extra dressing"
    },
    {
      "menu_item_id": 5,
      "quantity": 1,
      "special_request": ""
    }
  ]
}
```

**Response (201):**
```json
{
  "id": 1,
  "user_id": 1,
  "table_id": 1,
  "status": "pending",
  "total_amount": 34.07,
  "special_notes": "No onions please",
  "created_at": "2025-08-08T10:00:00Z",
  "items": [
    {
      "id": 1,
      "menu_item_id": 1,
      "quantity": 2,
      "price": 8.99,
      "special_request": "Extra dressing",
      "menu_item": {
        "name": "Spring Rolls",
        "description": "Crispy spring rolls with vegetables"
      }
    }
  ]
}
```

### GET /orders
Get user's orders (Authenticated users).

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10)
- `status`: Filter by status

**Response (200):**
```json
{
  "orders": [
    {
      "id": 1,
      "table_id": 1,
      "status": "pending",
      "total_amount": 34.07,
      "created_at": "2025-08-08T10:00:00Z"
    }
  ]
}
```

### GET /orders/{id}
Get specific order details (Authenticated users - own orders only).

**Response (200):**
```json
{
  "id": 1,
  "user_id": 1,
  "table_id": 1,
  "status": "confirmed",
  "total_amount": 34.07,
  "special_notes": "No onions please",
  "created_at": "2025-08-08T10:00:00Z",
  "items": [
    {
      "id": 1,
      "menu_item_id": 1,
      "quantity": 2,
      "price": 8.99,
      "special_request": "Extra dressing"
    }
  ]
}
```

### GET /admin/orders
Get all orders with advanced filtering (Admin/Staff).

**Query Parameters:**
- `page`: Page number
- `limit`: Items per page
- `status`: Filter by status (pending, confirmed, preparing, ready, served, cancelled)
- `user_id`: Filter by user
- `table_id`: Filter by table

**Response (200):**
```json
{
  "orders": [
    {
      "id": 1,
      "user_id": 1,
      "table_id": 1,
      "status": "confirmed",
      "total_amount": 34.07,
      "created_at": "2025-08-08T10:00:00Z",
      "user": {
        "name": "John Doe",
        "email": "john@example.com"
      },
      "table": {
        "number": "T001"
      }
    }
  ],
  "total": 150,
  "page": 1,
  "limit": 10,
  "total_pages": 15
}
```

### POST /cashier/orders
Create a new order with VAT calculation (Cashier/Admin only).

**Request Body:**
```json
{
  "table_id": 1,
  "customer_name": "John Doe",
  "cashier_name": "Cashier One",
  "special_notes": "Extra spicy",
  "items": [
    {
      "menu_item_id": 1,
      "quantity": 2,
      "special_request": "Extra sauce"
    },
    {
      "menu_item_id": 2,
      "quantity": 1,
      "special_request": ""
    }
  ]
}
```

**Response (201):**
```json
{
  "id": 1,
  "user_id": 2,
  "table_id": 1,
  "customer_name": "John Doe",
  "cashier_name": "Cashier One",
  "status": "pending",
  "subtotal": 30.97,
  "vat_amount": 3.10,
  "total_amount": 34.07,
  "special_notes": "Extra spicy",
  "created_at": "2025-08-08T10:00:00Z",
  "items": [
    {
      "id": 1,
      "menu_item_id": 1,
      "quantity": 2,
      "price": 8.99,
      "special_request": "Extra sauce",
      "menu_item": {
        "name": "Spring Rolls",
        "description": "Crispy spring rolls with vegetables"
      }
    },
    {
      "id": 2,
      "menu_item_id": 2,
      "quantity": 1,
      "price": 12.99,
      "special_request": "",
      "menu_item": {
        "name": "Chicken Wings",
        "description": "Spicy buffalo chicken wings"
      }
    }
  ]
}
```

### PATCH /admin/orders/{id}/status
Update order status (Admin/Staff).

**Request Body:**
```json
{
  "status": "preparing"
}
```

**Valid Status Transitions:**
- pending → confirmed, cancelled
- confirmed → preparing, cancelled
- preparing → ready
- ready → served
- served (final)
- cancelled (final)

### GET /admin/orders/statistics
Get order statistics (Admin only).

**Query Parameters:**
- `from`: Start date (YYYY-MM-DD)
- `to`: End date (YYYY-MM-DD)

**Response (200):**
```json
{
  "total_orders": 150,
  "orders_by_status": {
    "pending": 5,
    "confirmed": 10,
    "preparing": 8,
    "ready": 3,
    "served": 120,
    "cancelled": 4
  },
  "average_order_value": 29.50,
  "total_revenue": 4425.00,
  "orders_per_day": [
    {
      "date": "2025-08-08",
      "count": 15,
      "revenue": 442.50
    }
  ]
}
```

---

## 5. Payment Management

### POST /payments/qris
Initiate QRIS payment (Authenticated users).

**Request Body:**
```json
{
  "order_id": 1,
  "payment_method": "qris"
}
```

**Response (200):**
```json
{
  "payment_id": 1,
  "qr_code": "data:image/png;base64,iVBORw0KGgoAAAANS...",
  "transaction_id": "RD1642680000abc123",
  "amount": 34.07,
  "expires_at": "2025-08-08T10:15:00Z"
}
```

### POST /payments/verify
Verify payment status (Authenticated users).

**Request Body:**
```json
{
  "transaction_id": "RD1642680000abc123",
  "external_id": "QRIS123456789",
  "amount": 34.07,
  "status": "success"
}
```

### POST /cashier/payments/cash
Process cash payment (Cashier/Admin).

**Request Body:**
```json
{
  "order_id": 1,
  "amount_paid": 35.00,
  "change_amount": 0.93
}
```

**Response (200):**
```json
{
  "payment_id": 2,
  "transaction_id": "RD1642680000def456",
  "amount_paid": 35.00,
  "change_amount": 0.93,
  "order_total": 34.07,
  "message": "Cash payment processed successfully"
}
```

### GET /admin/payments
Get all payments (Admin/Cashier).

**Query Parameters:**
- `page`: Page number
- `limit`: Items per page
- `status`: Filter by status (pending, completed, failed, refunded, cancelled)
- `method`: Filter by method (qris, cash)

**Response (200):**
```json
{
  "payments": [
    {
      "id": 1,
      "order_id": 1,
      "amount": 34.07,
      "method": "qris",
      "status": "completed",
      "transaction_id": "RD1642680000abc123",
      "created_at": "2025-08-08T10:00:00Z",
      "order": {
        "id": 1,
        "table": {
          "number": "T001"
        }
      }
    }
  ],
  "total": 75,
  "page": 1,
  "limit": 10,
  "total_pages": 8
}
```

### POST /admin/payments/{id}/refund
Process payment refund (Admin/Cashier).

**Request Body:**
```json
{
  "amount": 34.07,
  "reason": "Customer complaint - food quality"
}
```

**Response (200):**
```json
{
  "refund_id": 3,
  "original_payment_id": 1,
  "refund_amount": 34.07,
  "reason": "Customer complaint - food quality",
  "status": "completed",
  "processed_at": "2025-08-08T11:00:00Z"
}
```

### POST /cashier/payments/reconcile
Reconcile cash payments for shift (Cashier only).

**Request Body:**
```json
{
  "actual_cash_amount": 450.75,
  "expected_cash_amount": 455.00,
  "shift_start_time": "2025-07-20T08:00:00Z",
  "shift_end_time": "2025-07-20T16:00:00Z",
  "notes": "5 dollar bill missing, investigating"
}
```

**Response (200):**
```json
{
  "cashier_id": 2,
  "shift_start": "2025-07-20T08:00:00Z",
  "shift_end": "2025-07-20T16:00:00Z",
  "expected_amount": 455.00,
  "calculated_amount": 455.00,
  "actual_amount": 450.75,
  "difference": -4.25,
  "payment_count": 18,
  "notes": "5 dollar bill missing, investigating",
  "reconciliation_time": "2025-07-20T16:30:00Z"
}
```

### GET /admin/payments/statistics
Get payment statistics (Admin/Cashier).

**Query Parameters:**
- `from`: Start date (YYYY-MM-DD)
- `to`: End date (YYYY-MM-DD)
- `method`: Filter by payment method

**Response (200):**
```json
{
  "total_payments": 75,
  "total_revenue": 2137.50,
  "payments_by_method": {
    "qris": {
      "count": 45,
      "total_amount": 1287.30
    },
    "cash": {
      "count": 30,
      "total_amount": 850.20
    }
  },
  "payments_by_status": {
    "completed": 70,
    "pending": 2,
    "failed": 2,
    "refunded": 1
  },
  "average_payment": 28.50,
  "refund_rate": 1.33
}
```

---

## 6. User Management (Admin Only)

### GET /admin/users
Get all users with advanced filtering, pagination, and search.

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)
- `role`: Filter by role (customer, staff, cashier, admin)
- `search`: Search by name, username, or email
- `status`: Filter by status (active, inactive, all - default: all)

**Response (200):**
```json
{
  "users": [
    {
      "id": 1,
      "username": "john_doe",
      "name": "John Doe",
      "email": "john@example.com",
      "phone": "+1234567890",
      "role": "customer",
      "is_active": true,
      "created_at": "2025-08-13T10:00:00Z",
      "updated_at": "2025-08-13T10:00:00Z"
    }
  ],
  "total": 50,
  "page": 1,
  "limit": 10,
  "total_pages": 5,
  "filters": {
    "role": "customer",
    "search": "john",
    "status": "active"
  }
}
```

### GET /admin/users/{id}
Get specific user details by ID.

**Response (200):**
```json
{
  "id": 1,
  "username": "john_doe",
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+1234567890",
  "role": "customer",
  "is_active": true,
  "created_at": "2025-08-13T10:00:00Z",
  "updated_at": "2025-08-13T10:00:00Z"
}
```

### POST /admin/users
Create a new user account.

**Request Body:**
```json
{
  "username": "jane_smith",
  "name": "Jane Smith",
  "email": "jane@example.com",
  "password": "securepassword123",
  "phone": "+1234567891",
  "role": "staff"
}
```

**Response (201):**
```json
{
  "id": 2,
  "username": "jane_smith",
  "name": "Jane Smith",
  "email": "jane@example.com",
  "phone": "+1234567891",
  "role": "staff",
  "is_active": true,
  "created_at": "2025-08-13T10:00:00Z",
  "updated_at": "2025-08-13T10:00:00Z"
}
```

### PUT /admin/users/{id}
Update user information completely.

**Request Body:**
```json
{
  "username": "jane_smith_updated",
  "name": "Jane Smith Updated",
  "email": "jane.updated@example.com",
  "phone": "+1234567892",
  "role": "cashier"
}
```

**Response (200):**
```json
{
  "id": 2,
  "username": "jane_smith_updated",
  "name": "Jane Smith Updated",
  "email": "jane.updated@example.com",
  "phone": "+1234567892",
  "role": "cashier",
  "is_active": true,
  "created_at": "2025-08-13T10:00:00Z",
  "updated_at": "2025-08-13T11:00:00Z"
}
```

### DELETE /admin/users/{id}
Soft delete a user (cannot delete own account).

**Response (200):**
```json
{
  "message": "User deleted successfully"
}
```

### PATCH /admin/users/{id}/status
Activate or deactivate user account.

**Request Body:**
```json
{
  "is_active": false
}
```

**Response (200):**
```json
{
  "message": "User deactivated successfully"
}
```

### PATCH /admin/users/{id}/role
Update user role (cannot change own role).

**Request Body:**
```json
{
  "role": "cashier"
}
```

**Response (200):**
```json
{
  "id": 2,
  "username": "jane_smith",
  "name": "Jane Smith",
  "email": "jane@example.com",
  "phone": "+1234567891",
  "role": "cashier",
  "is_active": true,
  "created_at": "2025-08-13T10:00:00Z",
  "updated_at": "2025-08-13T11:30:00Z"
}
```

### PATCH /admin/users/{id}/password
Reset user password.

**Request Body:**
```json
{
  "password": "newSecurePassword123"
}
```

**Response (200):**
```json
{
  "message": "Password reset successfully"
}
```

### GET /admin/users/statistics
Get comprehensive user statistics.

**Response (200):**
```json
{
  "total_users": 150,
  "users_by_role": {
    "customer": 120,
    "staff": 20,
    "cashier": 8,
    "admin": 2
  },
  "users_by_status": {
    "active": 145,
    "inactive": 5
  },
  "recent_registrations": [
    {
      "date": "2025-08-13",
      "count": 5
    },
    {
      "date": "2025-08-12",
      "count": 3
    }
  ],
  "role_distribution_percentage": {
    "customer": 80.0,
    "staff": 13.3,
    "cashier": 5.3,
    "admin": 1.4
  }
}
```

### PATCH /admin/users/bulk
Bulk update multiple users.

**Request Body:**
```json
{
  "user_ids": [1, 2, 3, 4],
  "updates": {
    "is_active": true,
    "role": "staff"
  }
}
```

**Response (200):**
```json
{
  "updated_count": 3,
  "failed_count": 1,
  "failed_users": [
    {
      "user_id": 1,
      "error": "Cannot change admin user role"
    }
  ],
  "message": "Bulk update completed: 3 users updated, 1 failed"
}
```

**Error Responses:**
- **400**: Invalid user ID, validation errors, cannot modify own account
- **404**: User not found
- **409**: Username or email already exists (for creation/updates)

---

## 7. Kitchen Management

### WebSocket /kitchen/updates
Real-time kitchen updates for new orders.

**Connection:**
```javascript
const ws = new WebSocket('ws://localhost:8002/kitchen/updates');
ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  console.log('New order:', data);
};
```

**Message Format:**
```json
{
  "type": "new_order",
  "order": {
    "id": 1,
    "table_number": "T001",
    "items": [
      {
        "name": "Caesar Salad",
        "quantity": 2,
        "special_request": "Extra dressing"
      }
    ],
    "special_notes": "No onions please",
    "created_at": "2025-07-20T10:00:00Z"
  }
}
```

---

## Error Handling

All endpoints return consistent error responses:

**400 Bad Request:**
```json
{
  "error": "Invalid request format",
  "details": {
    "field": "email",
    "message": "Invalid email format"
  }
}
```

**401 Unauthorized:**
```json
{
  "error": "Authentication required"
}
```

**403 Forbidden:**
```json
{
  "error": "Insufficient permissions"
}
```

**404 Not Found:**
```json
{
  "error": "Resource not found"
}
```

**500 Internal Server Error:**
```json
{
  "error": "Internal server error"
}
```

---

## Rate Limiting

API requests are rate limited:
- **Public endpoints**: 100 requests per minute
- **Authenticated endpoints**: 200 requests per minute
- **Admin endpoints**: 500 requests per minute

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 200
X-RateLimit-Remaining: 199
X-RateLimit-Reset: 1642680000
```

---

## Health Check & Monitoring

### GET /health
Check API health status.

**Response (200):**
```json
{
  "status": "healthy",
  "timestamp": "2025-08-08T10:00:00Z",
  "version": "1.0.0",
  "database": "connected",
  "features": {
    "vat_calculation": "enabled",
    "payment_methods": ["qris", "cash"],
    "roles": ["customer", "staff", "cashier", "admin"]
  }
}
```

### GET /metrics
Prometheus metrics endpoint for monitoring.

---

## Swagger Documentation

Interactive API documentation is available at:
```
http://localhost:8002/swagger/index.html
```

This provides a complete interface for testing all endpoints with proper authentication and request/response examples.
