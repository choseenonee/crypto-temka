-- +goose Up
-- +goose StatementBegin
ALTER TABLE rates
    ADD COLUMN is_once BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE rates
    DROP COLUMN is_once;
-- +goose StatementEnd
