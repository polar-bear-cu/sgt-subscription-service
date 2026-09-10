CREATE TABLE
    subscriptions (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id UUID NOT NULL,
        name TEXT NOT NULL
    );