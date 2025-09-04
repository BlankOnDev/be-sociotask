-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks
ADD COLUMN user_id BIGINT REFERENCES users (id),
ADD COLUMN reward_usdt FLOAT DEFAULT 0 NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks DROP COLUMN user_id, DROP COLUMN reward_usdt;
-- +goose StatementEnd