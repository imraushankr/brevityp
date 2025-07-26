-- Brevity Migration: create_credit_table
-- Generated: 2025-07-25T17:41:06Z
-- Direction: UP

-- Add your SQL below this line

CREATE TABLE
  credits (
    id VARCHAR(20) PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (
      type IN ('free', 'trial', 'paid', 'promo', 'referral')
    ),
    amount INTEGER NOT NULL DEFAULT 0,
    remaining INTEGER NOT NULL DEFAULT 0,
    expires_at TIMESTAMP,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
  );

CREATE INDEX idx_credits_user_id ON credits (user_id);

CREATE TABLE
  credit_usages (
    id VARCHAR(20) PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    credit_id VARCHAR(20) NOT NULL,
    url_id VARCHAR(20),
    amount INTEGER NOT NULL DEFAULT 1,
    operation VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (credit_id) REFERENCES credits (id) ON DELETE CASCADE,
    FOREIGN KEY (url_id) REFERENCES urls (id) ON DELETE SET NULL
  );

CREATE INDEX idx_credit_usages_user_id ON credit_usages (user_id);

CREATE INDEX idx_credit_usages_credit_id ON credit_usages (credit_id);

CREATE INDEX idx_credit_usages_url_id ON credit_usages (url_id);