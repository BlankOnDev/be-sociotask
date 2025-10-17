-- +goose Up
-- +goose StatementBegin
CREATE TYPE status_type AS ENUM ('PENDING', 'ACTIVE', 'COMPLETED');

ALTER TABLE tasks
ADD COLUMN status status_type DEFAULT 'PENDING';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks DROP COLUMN status;
DROP TYPE status_type;
-- +goose StatementEnd
