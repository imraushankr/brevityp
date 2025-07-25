-- Brevity Migration: create_subscription_table
-- Generated: 2025-07-25T17:43:56Z
-- Direction: DOWN

-- Add your SQL below this line

DROP TABLE IF EXISTS payments;

DROP TABLE IF EXISTS subscriptions;

DROP TYPE IF EXISTS subscription_plan;