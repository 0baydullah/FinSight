CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,

    user_id BIGINT NOT NULL,

    name VARCHAR(100) NOT NULL,

    is_system BOOLEAN NOT NULL DEFAULT FALSE,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_categories_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_categories_user_id
ON categories(user_id);

CREATE INDEX idx_categories_deleted_at
ON categories(deleted_at);

CREATE UNIQUE INDEX idx_categories_user_name
ON categories(user_id, LOWER(name))
WHERE deleted_at IS NULL;