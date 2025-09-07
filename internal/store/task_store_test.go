package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migration test db error: %v", err)
	}

	_, err = db.Exec("TRUNCATE tasks CASCADE")
	if err != nil {
		t.Fatalf("truncating table: %v", err)
	}

	return db
}

func TestCreateTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresTaskStore(db)

	tests := []struct {
		name    string
		task    *Task
		wantErr bool
	}{
		{
			name: "valid task",
			task: &Task{
				Title:       "Repost X Post",
				Description: "Repost X Post for crypto",
				UserID:      1,
				RewardUSDT:  1.15,
			},
			wantErr: false,
		},
		{
			name: "invalid task",
			task: &Task{
				Title:       "Like Instagram Post",
				Description: "Like Instagram Post for crypto",
				RewardUSDT:  0.377,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdTask, err := store.CreateTask(tt.task)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.task.Title, createdTask.Title)

			// Check in the database
			retrieved, err := store.GetTaskByID(int64(createdTask.ID))
			assert.NoError(t, err)
			assert.Equal(t, createdTask.ID, retrieved.ID)
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresTaskStore(db)

	initialTask, err := store.CreateTask(&Task{
		Title:       "Test",
		Description: "Test",
		UserID:      1,
		RewardUSDT:  1,
	})
	if err != nil {
		panic("failed to add initial task for testing")
	}

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "valid id",
			id:      int64(initialTask.ID),
			wantErr: false,
		},
		{
			name:    "non-existent id",
			id:      99999, // This is random ID
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := store.GetTaskByID(tt.id)
			if tt.wantErr {
				assert.Equal(t, nil, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.id, int64(task.ID))
		})
	}
}
