DROP INDEX IF EXISTS idx_subscriptions_user_id;

ALTER TABLE subscriptions
DROP CONSTRAINT IF EXISTS subscriptions_reminder_check,
DROP CONSTRAINT IF EXISTS subscriptions_cost_check,
DROP CONSTRAINT IF EXISTS subscriptions_status_check,
DROP CONSTRAINT IF EXISTS subscriptions_category_check,
DROP CONSTRAINT IF EXISTS subscriptions_type_check;

ALTER TABLE subscriptions
DROP COLUMN IF EXISTS updated_at,
DROP COLUMN IF EXISTS created_at,
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS ft_end_date,
DROP COLUMN IF EXISTS reminder_time_in_advanced,
DROP COLUMN IF EXISTS category,
DROP COLUMN IF EXISTS type,
DROP COLUMN IF EXISTS cost,
ALTER COLUMN name TYPE TEXT;

ALTER TABLE subscriptions
RENAME COLUMN next_billing_date TO billing_date;
