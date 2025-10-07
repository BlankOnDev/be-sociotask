-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
    CREATE TYPE status_rewards AS ENUM('crypto_usdt_1', 'crypto_usdt_2', 'crypto_usdt_3');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS task_rewards(
	id BIGSERIAL PRIMARY KEY,
    reward_type status_rewards,
    reward_name VARCHAR(255)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS task_rewards;
DROP TYPE IF EXISTS reward_type;
-- +goose StatementEnd