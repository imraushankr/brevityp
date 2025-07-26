-- Brevity Migration: create_url_table
-- Generated: 2025-07-25T17:34:27Z
-- Direction: UP
-- Add your SQL below this line
-- CREATE TABLE
--   urls (
--     id VARCHAR(20) PRIMARY KEY,
--     original_url TEXT NOT NULL,
--     short_code VARCHAR(10) NOT NULL UNIQUE,
--     user_id VARCHAR(20),
--     created_by_ip VARCHAR(45),
--     title VARCHAR(100),
--     description VARCHAR(255),
--     clicks INTEGER DEFAULT 0,
--     expires_at TIMESTAMP
--     WITH
--       TIME ZONE,
--       is_active BOOLEAN DEFAULT true,
--       created_at TIMESTAMP
--     WITH
--       TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
--       updated_at TIMESTAMP
--     WITH
--       TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
--       deleted_at TIMESTAMP
--     WITH
--       TIME ZONE,
--       FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
--   );
CREATE TABLE
  urls (
    id VARCHAR(20) PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) NOT NULL UNIQUE,
    user_id VARCHAR(20) NULL, -- Explicitly marked as NULLable
    created_by_ip VARCHAR(45),
    title VARCHAR(100),
    description VARCHAR(255),
    clicks INTEGER DEFAULT 0,
    expires_at TIMESTAMP
    WITH
      TIME ZONE,
      is_active BOOLEAN DEFAULT true,
      created_at TIMESTAMP
    WITH
      TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP
    WITH
      TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
      deleted_at TIMESTAMP
    WITH
      TIME ZONE,
      FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
  );

CREATE INDEX idx_urls_user_id ON urls (user_id);

CREATE INDEX idx_urls_deleted_at ON urls (deleted_at);

CREATE TABLE
  url_clicks (
    id VARCHAR(20) PRIMARY KEY,
    url_id VARCHAR(20) NOT NULL,
    ip_address VARCHAR(45),
    referrer TEXT,
    user_agent TEXT,
    country VARCHAR(2),
    city TEXT,
    device VARCHAR(20),
    os VARCHAR(20),
    browser VARCHAR(20),
    created_at TIMESTAMP
    WITH
      TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
      FOREIGN KEY (url_id) REFERENCES urls (id)
  );

CREATE INDEX idx_url_clicks_url_id ON url_clicks (url_id);