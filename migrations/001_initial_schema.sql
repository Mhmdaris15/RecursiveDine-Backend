-- Initial schema migration
-- Created: 2025-08-08

-- Connect to the recursive_dine database
-- \c recursive_dine;

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'customer',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create tables table
CREATE TABLE tables (
    id SERIAL PRIMARY KEY,
    number INTEGER UNIQUE NOT NULL,
    qr_code VARCHAR(255) UNIQUE NOT NULL,
    capacity INTEGER NOT NULL,
    is_available BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create menu_categories table
CREATE TABLE menu_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create menu_items table
CREATE TABLE menu_items (
    id SERIAL PRIMARY KEY,
    category_id INTEGER REFERENCES menu_categories(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    image_url VARCHAR(255),
    is_available BOOLEAN DEFAULT TRUE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create orders table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    table_id INTEGER REFERENCES tables(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    total_amount DECIMAL(10,2) NOT NULL,
    special_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create order_items table
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id INTEGER REFERENCES menu_items(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    special_request TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create payments table
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INTEGER UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
    method VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    amount DECIMAL(10,2) NOT NULL,
    qris_data TEXT,
    transaction_id VARCHAR(255),
    external_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

CREATE INDEX idx_tables_number ON tables(number);
CREATE INDEX idx_tables_qr_code ON tables(qr_code);
CREATE INDEX idx_tables_deleted_at ON tables(deleted_at);

CREATE INDEX idx_menu_categories_active ON menu_categories(is_active);
CREATE INDEX idx_menu_categories_sort ON menu_categories(sort_order);
CREATE INDEX idx_menu_categories_deleted_at ON menu_categories(deleted_at);

CREATE INDEX idx_menu_items_category ON menu_items(category_id);
CREATE INDEX idx_menu_items_available ON menu_items(is_available);
CREATE INDEX idx_menu_items_sort ON menu_items(sort_order);
CREATE INDEX idx_menu_items_deleted_at ON menu_items(deleted_at);

CREATE INDEX idx_orders_user ON orders(user_id);
CREATE INDEX idx_orders_table ON orders(table_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_deleted_at ON orders(deleted_at);

CREATE INDEX idx_order_items_order ON order_items(order_id);
CREATE INDEX idx_order_items_menu_item ON order_items(menu_item_id);

CREATE INDEX idx_payments_order ON payments(order_id);
CREATE INDEX idx_payments_transaction ON payments(transaction_id);
CREATE INDEX idx_payments_external ON payments(external_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_deleted_at ON payments(deleted_at);

-- Insert sample data with properly hashed passwords
-- admin@recursivedine.com: password 'admin123'
-- staff1@recursivedine.com: password 'password123'  
-- customer1@example.com: password 'password123'
INSERT INTO users (username, email, password, role) VALUES 
('admin', 'admin@recursivedine.com', '$2a$10$jrezlNG2pZrkppFHamKbseC5IC0WxzX/WQm5U9Bl.i7NWOCun5TMO', 'admin'),
('staff1', 'staff1@recursivedine.com', '$2a$10$yX8EQib9UrToO3ThKGZ.VO.4QFETRzrokCckN7H473STJO5sQ2viC', 'staff'),
('customer1', 'customer1@example.com', '$2a$10$yX8EQib9UrToO3ThKGZ.VO.4QFETRzrokCckN7H473STJO5sQ2viC', 'customer');

INSERT INTO tables (number, qr_code, capacity) VALUES 
(1, 'QR001', 4),
(2, 'QR002', 2),
(3, 'QR003', 6),
(4, 'QR004', 4),
(5, 'QR005', 8);

INSERT INTO menu_categories (name, description, sort_order) VALUES 
('Appetizers', 'Start your meal with our delicious appetizers', 1),
('Main Courses', 'Hearty main dishes to satisfy your hunger', 2),
('Desserts', 'Sweet endings to your meal', 3),
('Beverages', 'Refreshing drinks and beverages', 4);

INSERT INTO menu_items (category_id, name, description, price, sort_order) VALUES 
(1, 'Spring Rolls', 'Crispy spring rolls with vegetables', 8.99, 1),
(1, 'Chicken Wings', 'Spicy buffalo chicken wings', 12.99, 2),
(2, 'Grilled Salmon', 'Fresh salmon with herbs and lemon', 24.99, 1),
(2, 'Beef Steak', 'Tender beef steak with garlic butter', 28.99, 2),
(2, 'Pasta Carbonara', 'Creamy pasta with bacon and eggs', 18.99, 3),
(3, 'Chocolate Cake', 'Rich chocolate cake with frosting', 7.99, 1),
(3, 'Ice Cream', 'Vanilla ice cream with toppings', 5.99, 2),
(4, 'Coffee', 'Freshly brewed coffee', 3.99, 1),
(4, 'Fresh Juice', 'Orange or apple juice', 4.99, 2),
(4, 'Soft Drinks', 'Coca-Cola, Sprite, or Fanta', 2.99, 3);
