-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR UNIQUE,
    phone_number VARCHAR,
    hashed_password VARCHAR,
    status VARCHAR CHECK (status IN ('opened', 'declined', 'verified')),
    refer_id INTEGER REFERENCES users(id),
    properties JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
