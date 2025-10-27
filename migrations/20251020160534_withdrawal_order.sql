-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS withdrawal_orders
(
    id    int GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    user_id int,
    accrual float,
    order_number varchar,
    processed_at timestamptz DEFAULT now()
);
ALTER TABLE withdrawal_orders
    ADD CONSTRAINT FK_withdrawal_orders_users FOREIGN KEY (user_id)
        REFERENCES users (id);

-- +goose StatementEnd
