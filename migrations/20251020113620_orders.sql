-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders
(
    id    int GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    user_id int,
    number  varchar,
    status varchar,
    accrual float,
    uploaded_at timestamptz DEFAULT now()
);
ALTER TABLE orders
    ADD CONSTRAINT FK_orders_users FOREIGN KEY (user_id)
        REFERENCES users (id);

-- +goose StatementEnd
