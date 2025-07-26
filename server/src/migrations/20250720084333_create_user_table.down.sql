-- Brevity Migration: create_user_table
-- Generated: 2025-07-20T08:43:33Z
-- Direction: DOWN
-- Add your SQL below this line

-- Drop indexes first
DROP INDEX IF EXISTS idx_users_reset_password_token;

DROP INDEX IF EXISTS idx_users_verification_token;

DROP INDEX IF EXISTS idx_users_username;

DROP INDEX IF EXISTS idx_users_email;

DROP INDEX IF EXISTS idx_users_deleted_at;

-- Drop the users table
DROP TABLE IF EXISTS users;