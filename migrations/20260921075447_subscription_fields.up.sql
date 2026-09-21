ALTER TABLE subscriptions
RENAME COLUMN billing_date TO next_billing_date;

ALTER TABLE subscriptions
ALTER COLUMN name TYPE VARCHAR(50),
ADD COLUMN cost NUMERIC(10, 2) NOT NULL DEFAULT 0,
ADD COLUMN type TEXT NOT NULL DEFAULT 'monthly',
ADD COLUMN category TEXT NOT NULL DEFAULT 'streaming',
ADD COLUMN reminder_time_in_advanced BIGINT NOT NULL DEFAULT 1,
ADD COLUMN ft_end_date TIMESTAMPTZ,
ADD COLUMN status TEXT NOT NULL DEFAULT 'active',
ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT now ();

ALTER TABLE subscriptions
ADD CONSTRAINT subscriptions_type_check CHECK (type IN ('monthly', 'yearly')),
ADD CONSTRAINT subscriptions_category_check CHECK (
    category IN ('streaming', 'music', 'productivity', 'technology')
),
ADD CONSTRAINT subscriptions_status_check CHECK (status IN ('active', 'free_trial', 'inactive')),
ADD CONSTRAINT subscriptions_cost_check CHECK (cost >= 0),
ADD CONSTRAINT subscriptions_reminder_check CHECK (reminder_time_in_advanced >= 1);

CREATE INDEX idx_subscriptions_user_id ON subscriptions (user_id);
