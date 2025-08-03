-- +goose Up
-- +goose StatementBegin
ALTER TABLE withdrawals
    ADD COLUMN wallet_id INTEGER REFERENCES wallets(id),
    DROP COLUMN token;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE withdrawals
    ADD COLUMN token VARCHAR,
    DROP COLUMN wallet_id;
-- +goose StatementEnd
