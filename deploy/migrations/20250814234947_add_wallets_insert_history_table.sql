-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wallets_insert_history (
    id SERIAL primary key,
    wallet_id INTEGER REFERENCES wallets(id),
    amount DOUBLE PRECISION,
    ts TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets_insert_history;
-- +goose StatementEnd
