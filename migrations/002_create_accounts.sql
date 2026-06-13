CREATE TYPE account_type   AS ENUM ('savings', 'checking');
CREATE TYPE account_status AS ENUM ('active', 'blocked', 'closed');

CREATE TABLE accounts (
    id             UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id    UUID           NOT NULL REFERENCES customers(id),
    account_number VARCHAR(20)    NOT NULL UNIQUE,
    type           account_type   NOT NULL,
    currency       CHAR(3)        NOT NULL DEFAULT 'USD',
    balance        NUMERIC(18,2)  NOT NULL DEFAULT 0,
    status         account_status NOT NULL DEFAULT 'active',
    created_at     TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ    NOT NULL DEFAULT now(),

    CONSTRAINT balance_non_negative CHECK (balance >= 0)
);

CREATE INDEX idx_accounts_customer ON accounts(customer_id);
CREATE INDEX idx_accounts_number   ON accounts(account_number);
