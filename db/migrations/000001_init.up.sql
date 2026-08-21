CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    username   VARCHAR(255) UNIQUE NOT NULL,
    email      VARCHAR(255) UNIQUE NOT NULL,
    password   VARCHAR(255)        NOT NULL,
    created_at TIMESTAMPTZ         NOT NULL DEFAULT NOW()
);


CREATE TABLE wallets
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT UNIQUE NOT NULL,
    balance    BIGINT        NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_wallet_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE,

    CONSTRAINT wallet_balance_non_negative
        CHECK (balance >= 0)
);


CREATE TABLE transactions
(
    id            BIGSERIAL PRIMARY KEY,
    wallet_id     BIGINT      NOT NULL,
    type          VARCHAR(50) NOT NULL,
    amount        BIGINT      NOT NULL,
    balance_after BIGINT      NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_transaction_wallet
        FOREIGN KEY (wallet_id)
            REFERENCES wallets (id)
            ON DELETE CASCADE,

    CONSTRAINT transaction_amount_not_zero
        CHECK (amount <> 0),

    CONSTRAINT transaction_balance_non_negative
        CHECK (balance_after >= 0),

    CONSTRAINT transaction_type_valid
        CHECK (type IN ('deposit', 'withdrawal', 'transfer'))
);


CREATE TABLE transfers
(
    id                 BIGSERIAL PRIMARY KEY,
    sender_wallet_id   BIGINT      NOT NULL,
    receiver_wallet_id BIGINT      NOT NULL,
    amount             BIGINT      NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_transfer_sender
        FOREIGN KEY (sender_wallet_id)
            REFERENCES wallets (id)
            ON DELETE RESTRICT,

    CONSTRAINT fk_transfer_receiver
        FOREIGN KEY (receiver_wallet_id)
            REFERENCES wallets (id)
            ON DELETE RESTRICT,

    CONSTRAINT transfer_amount_positive
        CHECK (amount > 0),

    CONSTRAINT transfer_different_wallets
        CHECK (sender_wallet_id <> receiver_wallet_id)
);


-- Indexes

CREATE INDEX idx_transactions_wallet_id
    ON transactions (wallet_id);

CREATE INDEX idx_transactions_wallet_created_at
    ON transactions (wallet_id, created_at DESC);

CREATE INDEX idx_transfers_sender_wallet_id
    ON transfers (sender_wallet_id);

CREATE INDEX idx_transfers_receiver_wallet_id
    ON transfers (receiver_wallet_id);

CREATE INDEX idx_transfers_created_at
    ON transfers (created_at DESC);