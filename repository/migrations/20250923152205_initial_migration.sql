-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
    id VARCHAR(50) PRIMARY KEY, -- go's uuid converted to string (32 bytes or 36 bytes with dashes)
    name VARCHAR(100) NOT NULL UNIQUE,
    price INTEGER NOT NULL,
    description VARCHAR(512) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd
