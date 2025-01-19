-- +goose Up
-- +goose StatementBegin
-- Step 1: Drop the existing CHECK constraint
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_status_check;

-- Step 2: Add the new CHECK constraint
ALTER TABLE users ADD CONSTRAINT users_status_check
    CHECK (status IN ('opened', 'declined', 'verified', 'pending'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_status_check;

ALTER TABLE users ADD CONSTRAINT users_status_check
    CHECK (status IN ('opened', 'declined', 'verified'));
-- +goose StatementEnd
