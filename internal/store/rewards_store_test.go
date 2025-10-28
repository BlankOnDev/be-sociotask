package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDBRewards(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	require.NoError(t, err, "opening test db")

	err = Migrate(db, "../../migrations/")
	require.NoError(t, err, "migration test db error")

	_, err = db.Exec("TRUNCATE rewards, tasks, users CASCADE")
	require.NoError(t, err, "truncating table")

	return db
}

func TestCreateReward(t *testing.T) {
	db := setupTestDBRewards(t)
	defer db.Close()

	rewardsStore := NewPostgresRewardsStore(db)
	userStore := NewPostgresUserStore(db)
	taskStore := NewPostgresTaskStore(db)

	// Create initial user
	initialUser := &User{
		Username: "test-user",
		Email:    "test@gmail.com",
		Fullname: sql.NullString{
			String: "Test User",
			Valid:  true,
		},
	}
	initialUser.PasswordHash.Set("password123")
	createdUser, err := userStore.CreateUser(initialUser)
	require.NoError(t, err, "failed to create initial user")

	// Create initial task
	initialTask := &Task{
		Title:       "Test Task",
		Description: "Test Description",
		UserID:      createdUser.ID,
		RewardUSDT:  10,
	}
	createdTask, err := taskStore.CreateTask(initialTask)
	require.NoError(t, err, "failed to create initial task")

	tests := []struct {
		name    string
		reward  *Reward
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid reward",
			reward: &Reward{
				UserID: createdUser.ID,
				TaskID: int64(createdTask.ID),
			},
			wantErr: false,
		},
		{
			name: "missing user id",
			reward: &Reward{
				UserID: 0,
				TaskID: int64(createdTask.ID),
			},
			wantErr: true,
			errMsg:  "user id is required and cannot be zero",
		},
		{
			name: "missing task id",
			reward: &Reward{
				UserID: createdUser.ID,
				TaskID: 0,
			},
			wantErr: true,
			errMsg:  "task id is required and cannot be zero",
		},
		{
			name: "missing both ids",
			reward: &Reward{
				UserID: 0,
				TaskID: 0,
			},
			wantErr: true,
			errMsg:  "user id is required and cannot be zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := rewardsStore.Create(tt.reward)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Greater(t, result.ID, int64(0))
			assert.Equal(t, tt.reward.UserID, result.UserID)
			assert.Equal(t, tt.reward.TaskID, result.TaskID)
			assert.False(t, result.CreatedAt.IsZero())
		})
	}
}
