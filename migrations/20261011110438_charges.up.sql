CREATE TABLE
    charges (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        subscription_id UUID NOT NULL REFERENCES subscriptions (id) ON DELETE CASCADE,
        type TEXT NOT NULL,
        amount NUMERIC(10, 2) NOT NULL,
        charged_at TIMESTAMPTZ NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now (),
        CONSTRAINT charges_type_check CHECK (type IN ('monthly', 'yearly')),
        CONSTRAINT charges_amount_check CHECK (amount >= 0),
        UNIQUE (subscription_id, charged_at)
    );
