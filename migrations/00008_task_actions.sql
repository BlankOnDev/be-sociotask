-- +goose Up
-- +goose StatementBegin
DO $$ BEGIN
    CREATE TYPE status_actions AS ENUM('type_1', 'type_2', 'type_3');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS task_actions(
	id BIGSERIAL PRIMARY KEY,
    type status_actions,
    name VARCHAR(255),
    description text
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS task_actions;
DROP TYPE IF EXISTS status_actions;
-- +goose StatementEnd
