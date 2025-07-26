-- Brevity Migration: create_subscription_table
-- Generated: 2025-07-25T17:43:56Z
-- Direction: UP

-- Add your SQL below this line

CREATE TABLE
  subscriptions (
    id VARCHAR(20) PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    plan VARCHAR(20) NOT NULL CHECK (plan IN ('free', 'basic', 'pro', 'enterprise')),
    stripe_id VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    starts_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    renews_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
  );

CREATE INDEX idx_subscriptions_user_id ON subscriptions (user_id);

CREATE INDEX idx_subscriptions_stripe_id ON subscriptions (stripe_id);

CREATE TABLE
  payments (
    id VARCHAR(20) PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    subscription_id VARCHAR(20),
    amount INTEGER NOT NULL,
    currency VARCHAR(3) DEFAULT 'usd',
    stripe_id VARCHAR(255),
    status VARCHAR(20),
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    paid_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (subscription_id) REFERENCES subscriptions (id) ON DELETE SET NULL
  );

CREATE INDEX idx_payments_user_id ON payments (user_id);

CREATE INDEX idx_payments_subscription_id ON payments (subscription_id);

CREATE INDEX idx_payments_stripe_id ON payments (stripe_id);