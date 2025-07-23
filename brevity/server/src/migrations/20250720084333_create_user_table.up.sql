-- Brevity Migration: create_user_table
-- Generated: 2025-07-20T08:43:33Z
-- Direction: UP
-- Add your SQL below this line
CREATE TABLE
	users (
		id VARCHAR(20) PRIMARY KEY,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		username VARCHAR(30) UNIQUE NOT NULL,
		role VARCHAR(20) NOT NULL CHECK (role IN ('admin', 'user')),
		email VARCHAR(255) UNIQUE NOT NULL,
		phone VARCHAR(15),
		avatar TEXT,
		password VARCHAR(255) NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		is_verified BOOLEAN DEFAULT FALSE,
		verification_token VARCHAR(255),
		verification_expires_at TIMESTAMP,
		reset_password_token VARCHAR(255),
		reset_password_expires_at TIMESTAMP,
		last_login_at TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP
	);

-- Create indexes for soft deletes and performance
CREATE INDEX idx_users_deleted_at ON users (deleted_at);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_username ON users (username);

CREATE INDEX idx_users_verification_token ON users (verification_token);

CREATE INDEX idx_users_reset_password_token ON users (reset_password_token);