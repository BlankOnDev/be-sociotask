-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN x_id VARCHAR(255) UNIQUE DEFAULT NULL,
ADD COLUMN wallet_address VARCHAR(255) UNIQUE DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN x_id, DROP COLUMN wallet_address;
-- +goose StatementEnd