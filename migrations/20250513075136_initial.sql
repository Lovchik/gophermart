-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id    int GENERATED ALWAYS AS IDENTITY NOT NULL PRIMARY KEY,
    login text,
    password  text
);
-- +goose StatementEnd
