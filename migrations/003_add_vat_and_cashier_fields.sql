-- Connect to the recursive_dine database
-- \c recursive_dine;

-- Add VAT and cashier information to orders table
ALTER TABLE orders ADD COLUMN subtotal_amount DECIMAL(10,2) NOT NULL DEFAULT 0;
ALTER TABLE orders ADD COLUMN vat_amount DECIMAL(10,2) NOT NULL DEFAULT 0;
ALTER TABLE orders ADD COLUMN customer_name VARCHAR(255);
ALTER TABLE orders ADD COLUMN cashier_name VARCHAR(255);

-- Update existing orders to move total_amount to subtotal_amount and calculate VAT
UPDATE orders SET 
    subtotal_amount = total_amount / 1.1,
    vat_amount = total_amount - (total_amount / 1.1),
    total_amount = total_amount
WHERE subtotal_amount = 0;

-- Add comment to document VAT rate
COMMENT ON COLUMN orders.vat_amount IS 'VAT amount at 10% rate for Indonesia';
