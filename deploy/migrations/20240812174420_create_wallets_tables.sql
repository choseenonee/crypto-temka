-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wallets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    token VARCHAR,
    deposit FLOAT8 CHECK (deposit >= 0),
    outcome FLOAT8 CHECK ( outcome >= 0 ),
    UNIQUE (user_id, token)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallets;
-- +goose StatementEnd
