-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS vouchers (
  id VARCHAR primary key,
  type VARCHAR not null, -- ignore once
  properties JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS vouchers;
-- +goose StatementEnd
