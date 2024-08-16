-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS withdrawals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    amount FLOAT8 CHECK ( amount > 0 ),
    token VARCHAR,
    status VARCHAR CHECK ( status IN ('opened', 'declined', 'verified') ),
    properties JSONB
);

CREATE TABLE IF NOT EXISTS refers (
    id SERIAL PRIMARY KEY,
    parent_id INTEGER REFERENCES users(id),
    child_id INTEGER REFERENCES users(id),
    amount FLOAT8 CHECK ( amount >= 0 ),
    token VARCHAR,
    UNIQUE (child_id, token)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS withdrawals, refers;
-- +goose StatementEnd
