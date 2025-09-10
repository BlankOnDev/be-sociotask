-- +goose Up
-- +goose StatementBegin
-- Drop existing foreign key constraint
ALTER TABLE tasks DROP CONSTRAINT tasks_user_id_fkey;

-- Add new foreign key constraint with CASCADE delete
ALTER TABLE tasks 
ADD CONSTRAINT tasks_user_id_fkey 
FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop CASCADE constraint
ALTER TABLE tasks DROP CONSTRAINT tasks_user_id_fkey;

-- Add back original constraint without CASCADE
ALTER TABLE tasks 
ADD CONSTRAINT tasks_user_id_fkey 
FOREIGN KEY (user_id) REFERENCES users (id);
-- +goose StatementEnd