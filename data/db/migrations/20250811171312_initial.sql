-- +goose Up
-- +goose StatementBegin
CREATE TABLE product(
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    name VARCHAR(128) NOT NULL,
    description VARCHAR(512) NOT NULL,
    PRICE INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

--CREATE TABLE product_image(
--    product_id VARCHAR(50) PRIMARY KEY,
--    max_count INTEGER NOT NULL, -- depends on user status; default = 10
--    FOREIGN KEY (product_id) REFERENCES product(id) ON DELETE CASCADE
--);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE product_image;
DROP TABLE product;
-- +goose StatementEnd
