-- Migration: add_order_type_field
-- Created: 2025-08-14 15:34:06

-- Add order_type field to orders table to distinguish between dine-in and takeaway
ALTER TABLE orders ADD COLUMN order_type VARCHAR(20) NOT NULL DEFAULT 'dine_in';

-- Add constraint to ensure only valid order types
ALTER TABLE orders ADD CONSTRAINT chk_order_type CHECK (order_type IN ('dine_in', 'takeaway'));

-- Create index for order_type for faster queries
CREATE INDEX idx_orders_order_type ON orders(order_type);

-- Add estimated_completion_time for takeaway orders
ALTER TABLE orders ADD COLUMN estimated_completion_time TIMESTAMP;

-- Add customer_phone for takeaway orders (optional, for notifications)
ALTER TABLE orders ADD COLUMN customer_phone VARCHAR(20);
