-- Connect to the recursive_dine database
-- \c recursive_dine;

-- Add name and phone fields to users table
ALTER TABLE users 
ADD COLUMN IF NOT EXISTS name VARCHAR(100),
ADD COLUMN IF NOT EXISTS phone VARCHAR(20);

-- Create index for phone field
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);

-- Make the new fields not null with defaults for existing records
UPDATE users SET 
    name = CASE 
        WHEN username = 'admin' THEN 'Administrator'
        WHEN username = 'staff1' THEN 'Staff Member'
        WHEN username = 'customer1' THEN 'Customer User'
        ELSE CONCAT(UPPER(SUBSTRING(username, 1, 1)), SUBSTRING(username, 2))
    END
WHERE name IS NULL;

UPDATE users SET 
    phone = CASE 
        WHEN username = 'admin' THEN '+1234567001'
        WHEN username = 'staff1' THEN '+1234567002'
        WHEN username = 'customer1' THEN '+1234567003'
        ELSE '+1234567000'
    END
WHERE phone IS NULL;

-- Now make the fields non-nullable
ALTER TABLE users 
ALTER COLUMN name SET NOT NULL,
ALTER COLUMN phone SET NOT NULL;
