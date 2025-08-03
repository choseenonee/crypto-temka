-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS wallets
    ADD COLUMN is_outcome bool DEFAULT true,
    DROP CONSTRAINT IF EXISTS wallets_user_id_token_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS wallets
    DROP COLUMN is_outcome;
-- +goose StatementEnd
