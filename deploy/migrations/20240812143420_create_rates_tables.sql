-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rates (
    id SERIAL PRIMARY KEY,
    title VARCHAR,
    profit FLOAT8 CHECK (profit > 0),
    min_lock_days INTEGER CHECK ( min_lock_days % 7 = 0 ),
    commission INTEGER,
    properties JSONB
);

CREATE TABLE IF NOT EXISTS users_rates (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    rate_id INTEGER REFERENCES rates(id),
    lock DATE,
    last_updated DATE,
    opened DATE,
    deposit FLOAT8 CHECK (deposit >= 0),
    earned_pool FLOAT8 NOT NULL DEFAULT 0,
    next_day_charge FLOAT8 CHECK ( next_day_charge <= (earned_pool + outcome_pool + deposit) ),
    outcome_pool FLOAT8 NOT NULL DEFAULT 0 CHECK (outcome_pool >= 0),
    token VARCHAR
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rates, users_rates;
-- +goose StatementEnd
