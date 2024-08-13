-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rates (
    id SERIAL PRIMARY KEY,
    title VARCHAR,
    profit INTEGER CHECK (profit > 0),
    min_lock_days INTEGER,
    commission INTEGER,
    properties JSONB
);
-- TODO: user_id references ...
CREATE TABLE IF NOT EXISTS users_rates (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    rate_id INTEGER REFERENCES rates(id),
    lock DATE,
    last_updated DATE,
    opened DATE,
    deposit INTEGER CHECK (deposit >= 0),
    earned_pool INTEGER NOT NULL DEFAULT 0,
    next_day_charge INTEGER,
    outcome_pool INTEGER NOT NULL DEFAULT 0 CHECK (outcome_pool >= 0),
    token VARCHAR
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rates, users_rates;
-- +goose StatementEnd
