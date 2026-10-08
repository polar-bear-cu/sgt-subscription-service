ALTER TABLE subscriptions
DROP CONSTRAINT IF EXISTS subscriptions_billing_day_check,
DROP COLUMN IF EXISTS billing_day;
