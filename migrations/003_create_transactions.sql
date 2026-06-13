CREATE TYPE tx_type   AS ENUM ('deposit', 'withdrawal', 'transfer_in', 'transfer_out');
CREATE TYPE tx_status AS ENUM ('success', 'failed', 'reversed');

CREATE TABLE transactions (
    id            UUID          PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id    UUID          NOT NULL REFERENCES accounts(id),
    type          tx_type       NOT NULL,
    amount        NUMERIC(18,2) NOT NULL CHECK (amount > 0),
    balance_after NUMERIC(18,2) NOT NULL,
    reference     VARCHAR(64)   NOT NULL UNIQUE,
    related_tx_id UUID          REFERENCES transactions(id),
    status        tx_status     NOT NULL DEFAULT 'success',
    created_at    TIMESTAMPTZ   NOT NULL DEFAULT now()
);

CREATE INDEX idx_transactions_account   ON transactions(account_id, created_at DESC);

CREATE INDEX idx_transactions_ref       ON transactions(reference);

CREATE INDEX idx_transactions_type_date ON transactions(account_id, type, created_at DESC);
