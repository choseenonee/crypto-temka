-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cases (
    id SERIAL PRIMARY KEY,
    properties JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cases;
-- +goose StatementEnd
