-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS rewards(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    task_id BIGINT REFERENCES tasks(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rewards;
-- +goose StatementEnd
