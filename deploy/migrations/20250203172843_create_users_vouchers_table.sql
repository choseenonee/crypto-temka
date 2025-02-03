-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_vouchers (
    user_id INTEGER REFERENCES users(id),
    voucher_id VARCHAR REFERENCES vouchers(id),
    UNIQUE (user_id, voucher_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_vouchers;
-- +goose StatementEnd
