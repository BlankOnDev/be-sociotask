package store

import (
	"database/sql"
	"fmt"
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

	_, err = db.Exec("TRUNCATE tasks, users CASCADE")
	if err != nil {
		t.Fatalf("truncating table: %v", err)
	}

	return db
}

func TestCreateTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	taskStore := NewPostgresTaskStore(db)
	userStore := NewPostgresUserStore(db)

	initialUser := &User{
		Username: "test-create-task",
		Email:    "test-create-task@gmail.com",
		Bio:      "test-create-task",
	}
	initialUser.PasswordHash.Set("password123")

	initialUser, err := userStore.CreateUser(initialUser)
	require.NoError(t, err)

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
				UserID:      initialUser.ID,
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
			createdTask, err := taskStore.CreateTask(tt.task)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.task.Title, createdTask.Title)

			// Check in the database
			retrieved, err := taskStore.GetTaskByID(int64(createdTask.ID))
			assert.NoError(t, err)

			// Check if task is created properly
			assert.Equal(t, createdTask.ID, retrieved.ID)
			assert.Equal(t, createdTask.Title, retrieved.Title)
			assert.Equal(t, createdTask.Description, retrieved.Description)
			assert.Equal(t, createdTask.RewardUSDT, retrieved.RewardUSDT)
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	taskStore := NewPostgresTaskStore(db)
	userStore := NewPostgresUserStore(db)

	// Create initial user for happy path
	initialUser := &User{
		Username: "test-get-task-by-id",
		Email:    "test-get-task-by-id@gmail.com",
		Bio:      "test-get-task-by-id",
	}
	initialUser.PasswordHash.Set("password123")
	user, err := userStore.CreateUser(initialUser)
	if err != nil {
		t.Fatalf("failed to create initial user: %v", err)
	}

	// Create initial task for happy path
	initialTask, err := taskStore.CreateTask(&Task{
		Title:       "Test",
		Description: "Test",
		UserID:      user.ID,
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
			task, err := taskStore.GetTaskByID(tt.id)
			if tt.wantErr {
				assert.Equal(t, nil, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.id, int64(task.ID))
		})
	}
}

func TestGetAllTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	taskStore := NewPostgresTaskStore(db)
	userStore := NewPostgresUserStore(db)

	user := &User{
		Username: "test-get-all",
		Email:    "test-get-all@gmail.com",
		Bio:      "test get all",
	}
	user.PasswordHash.Set("password123")
	createdUser, err := userStore.CreateUser(user)
	require.NoError(t, err)

	for i := 1; i <= 7; i++ {
		_, err := taskStore.CreateTask(&Task{
			Title:       fmt.Sprintf("Task-%d", i),
			Description: "Test desc",
			UserID:      createdUser.ID,
			RewardUSDT:  float64(i),
		})
		require.NoError(t, err)
	}
	var page, limit int64
	page, limit = 1, 5
	offset := (page - 1) * limit

	tasks, totalPage, err := taskStore.GetAllTask(limit, offset)
	require.NoError(t, err)
	assert.LessOrEqual(t, len(tasks), limit)
	assert.GreaterOrEqual(t, totalPage, 1)
}

func TestUpdateTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	taskStore := NewPostgresTaskStore(db)
	userStore := NewPostgresUserStore(db)

	user := &User{
		Username: "test-update",
		Email:    "test-update@gmail.com",
		Bio:      "test update",
	}
	user.PasswordHash.Set("password123")
	createdUser, err := userStore.CreateUser(user)
	require.NoError(t, err)

	task, err := taskStore.CreateTask(&Task{
		Title:       "Before Update",
		Description: "desc",
		UserID:      createdUser.ID,
		RewardUSDT:  10,
	})
	require.NoError(t, err)

	task.Title = "After Update"
	task.Description = "Updated desc"
	err = taskStore.EditTask(task)
	require.NoError(t, err)

	updated, err := taskStore.GetTaskByID(int64(task.ID))
	require.NoError(t, err)
	assert.Equal(t, "After Update", updated.Title)
	assert.Equal(t, "Updated desc", updated.Description)
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	taskStore := NewPostgresTaskStore(db)
	userStore := NewPostgresUserStore(db)

	user := &User{
		Username: "test-delete",
		Email:    "test-delete@gmail.com",
		Bio:      "test delete",
	}
	user.PasswordHash.Set("password123")
	createdUser, err := userStore.CreateUser(user)
	require.NoError(t, err)

	task, err := taskStore.CreateTask(&Task{
		Title:       "Delete Me",
		Description: "desc",
		UserID:      createdUser.ID,
		RewardUSDT:  5,
	})
	require.NoError(t, err)

	err = taskStore.DeleteTask(int64(task.ID))
	require.NoError(t, err)

	deleted, err := taskStore.GetTaskByID(int64(task.ID))
	require.NoError(t, err)
	assert.Nil(t, deleted)
}
