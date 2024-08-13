-- +goose Up
-- +goose StatementBegin
-- TODO: user_id references ...
CREATE TABLE IF NOT EXISTS wallets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    token VARCHAR,
    deposit INTEGER CHECK (deposit >= 0),
    UNIQUE (user_id, token)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets;
-- +goose StatementEnd
