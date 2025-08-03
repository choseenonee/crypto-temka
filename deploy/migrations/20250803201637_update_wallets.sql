-- +goose Up
-- +goose StatementBegin
ALTER TABLE IF EXISTS wallets
    ADD COLUMN is_outcome bool DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE IF EXISTS wallets
    DROP COLUMN is_outcome;
-- +goose StatementEnd
