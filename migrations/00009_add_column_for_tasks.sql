-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks
-- ADD COLUMN user_id BIGINT,
ADD COLUMN reward_id INT,
ADD COLUMN due_date TIMESTAMPTZ,
ADD COLUMN max_participant VARCHAR(50),
ADD COLUMN task_image TEXT,
ADD COLUMN action_id INT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- ALTER TABLE tasks DROP COLUMN user_id;
ALTER TABLE tasks DROP COLUMN reward_id;
ALTER TABLE tasks DROP COLUMN due_date;
ALTER TABLE tasks DROP COLUMN max_participant;
ALTER TABLE tasks DROP COLUMN task_image;
ALTER TABLE tasks DROP COLUMN action_id;
-- +goose StatementEnd