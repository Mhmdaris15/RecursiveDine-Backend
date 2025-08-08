# RecursiveDine API - Cashier Order Testing Guide

## Features Implemented ✅

### 1. Cashier Order Workflow
- **Endpoint**: `POST /api/v1/cashier/orders`
- **Features**: 
  - Indonesian VAT calculation (10%)
  - Customer and Cashier name capture
  - Menu item selection with quantities
  - Subtotal and total amount calculation
  - Comprehensive error logging

### 2. Database Schema Updates
- **Migration 003**: Added VAT and cashier fields to orders table
  - `subtotal_amount` (DECIMAL)
  - `vat_amount` (DECIMAL) 
  - `customer_name` (VARCHAR)
  - `cashier_name` (VARCHAR)

### 3. Enhanced Logging System
- **Error logging**: All errors logged to `error.log` file
- **Request tracking**: User context, file/line information
- **Structured logging**: Consistent format with timestamps

### 4. Complete Testing Framework
- **Postman Collection**: `RecursiveDine_E2E_Testing.postman_collection.json`
- **Environment**: `RecursiveDine_E2E_Environment.postman_environment.json`
- **Test Coverage**: Authentication, seeding, cashier orders, VAT calculation, error handling

## Quick Testing Commands

### 1. Health Check
```bash
curl http://localhost:8002/health
```

### 2. Register and Login (Admin)
```bash
# Register
curl -X POST http://localhost:8002/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin User",
    "email": "admin@recursivedine.com", 
    "phone": "+62812345678",
    "password": "admin123"
  }'

# Login and save token
curl -X POST http://localhost:8002/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@recursivedine.com",
    "password": "admin123"
  }'
```

### 3. Seed Database
```bash
curl -X POST http://localhost:8002/api/v1/admin/seed \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### 4. Register Cashier and Test Order
```bash
# Register Cashier
curl -X POST http://localhost:8002/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cashier One",
    "email": "cashier1@recursivedine.com",
    "phone": "+62812345679", 
    "password": "cashier123"
  }'

# Login Cashier
curl -X POST http://localhost:8002/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "cashier1@recursivedine.com",
    "password": "cashier123"
  }'

# Create Cashier Order with VAT
curl -X POST http://localhost:8002/api/v1/cashier/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_CASHIER_TOKEN" \
  -d '{
    "table_id": 1,
    "customer_name": "John Doe",
    "cashier_name": "Cashier One",
    "special_notes": "Extra spicy, no onions",
    "items": [
      {
        "menu_item_id": 1,
        "quantity": 2,
        "special_request": "Extra sauce"
      },
      {
        "menu_item_id": 2, 
        "quantity": 1,
        "special_request": "Well done"
      }
    ]
  }'
```

## Expected Response Format

```json
{
  "id": 1,
  "user_id": 2,
  "table_id": 1,
  "status": "pending",
  "subtotal_amount": 45000.00,
  "vat_amount": 4500.00,
  "total_amount": 49500.00,
  "customer_name": "John Doe",
  "cashier_name": "Cashier One",
  "special_notes": "Extra spicy, no onions",
  "created_at": "2024-01-01 12:00:00",
  "order_items": [
    {
      "id": 1,
      "menu_item_id": 1,
      "menu_item_name": "Nasi Goreng",
      "quantity": 2,
      "unit_price": 15000.00,
      "total_price": 30000.00,
      "special_request": "Extra sauce"
    },
    {
      "id": 2, 
      "menu_item_id": 2,
      "menu_item_name": "Ayam Bakar",
      "quantity": 1,
      "unit_price": 15000.00,
      "total_price": 15000.00,
      "special_request": "Well done"
    }
  ]
}
```

## VAT Calculation Details

- **VAT Rate**: 10% (Indonesian standard)
- **Calculation**: VAT Amount = Subtotal × 0.10
- **Total**: Total Amount = Subtotal + VAT Amount
- **Currency**: Indonesian Rupiah (Rp)

## Error Logging

All errors are logged to `/logs/error.log` with:
- Timestamp
- Error level (ERROR, WARNING, INFO)
- Context information
- File and line number
- User ID and request details

## Postman Testing

Import both files into Postman:
1. `RecursiveDine_E2E_Testing.postman_collection.json` - Complete test suite
2. `RecursiveDine_E2E_Environment.postman_environment.json` - Environment variables

The collection includes:
- ✅ Authentication workflow
- ✅ Database seeding
- ✅ Menu management
- ✅ Table management  
- ✅ **Cashier order workflow with VAT**
- ✅ Regular order workflow
- ✅ Payment processing
- ✅ Admin management
- ✅ Error handling tests

## Key Features Summary

1. **✅ Cashier-based ordering system**
2. **✅ Indonesian VAT calculation (10%)**
3. **✅ Customer and cashier name capture**
4. **✅ Enhanced order table schema**
5. **✅ Comprehensive error logging** 
6. **✅ Complete Postman testing suite**
7. **✅ End-to-end workflow testing**

The system is now ready for production use with proper Indonesian tax compliance and comprehensive testing coverage!
