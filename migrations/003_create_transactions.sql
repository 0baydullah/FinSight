CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,

    category_id BIGINT NOT NULL,

    type VARCHAR(20) NOT NULL,

    amount BIGINT NOT NULL,

    description VARCHAR(500),

    transaction_date TIMESTAMP WITH TIME ZONE NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_transactions_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_transactions_category
        FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON DELETE RESTRICT,

    CONSTRAINT chk_transactions_type
        CHECK (type IN ('income', 'expense')),

    CONSTRAINT chk_transactions_amount
        CHECK (amount > 0)
);

CREATE INDEX idx_transactions_user_id
ON transactions(user_id);

CREATE INDEX idx_transactions_category_id
ON transactions(category_id);

CREATE INDEX idx_transactions_transaction_date
ON transactions(transaction_date);

CREATE INDEX idx_transactions_deleted_at
ON transactions(deleted_at);