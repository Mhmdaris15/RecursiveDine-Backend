-- Database initialization script for RecursiveDine
-- This script creates the required database and sets up basic structure

-- Connect to postgres database to create our application database
\c postgres;

-- Create the application database if it doesn't exist
SELECT 'CREATE DATABASE recursive_dine'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'recursive_dine')\gexec

-- Grant all privileges to postgres user
GRANT ALL PRIVILEGES ON DATABASE recursive_dine TO postgres;

-- Connect to the newly created database
\c recursive_dine;

-- Create extensions that might be needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Log successful database creation
\echo 'Database recursive_dine created successfully'
