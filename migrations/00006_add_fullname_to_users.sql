-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN fullname VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN fullname;
-- +goose StatementEnd
