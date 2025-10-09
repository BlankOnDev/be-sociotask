package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDBReward(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	require.NoError(t, err, "opening test db")

	_, err = db.Exec("TRUNCATE task_rewards CASCADE")
	require.NoError(t, err, "truncating table")

	return db
}

func TestTaskRewardStore(t *testing.T) {
	db := setupTestDBReward(t)
	defer db.Close()

	store := NewPostgresTaskRewardStore(db)

	var createdRewardID int

	t.Run("CreateReward", func(t *testing.T) {
		reward := &RewardTask{
			RewardType: CryptoUsdt1,
			RewardName: "1 USDT",
		}
		id, err := store.CreateReward(reward)
		require.NoError(t, err)
		require.NotNil(t, id)
		require.Greater(t, *id, 0)
		createdRewardID = *id
	})

	t.Run("GetRewardByID", func(t *testing.T) {
		retrieved, err := store.GetRewardByID(createdRewardID)
		require.NoError(t, err)
		assert.Equal(t, createdRewardID, retrieved.ID)
		assert.Equal(t, CryptoUsdt1, retrieved.RewardType)
		assert.Equal(t, "1 USDT", retrieved.RewardName)
	})

	t.Run("EditReward", func(t *testing.T) {
		rewardToUpdate := &RewardTask{
			ID:         createdRewardID,
			RewardName: "2 USDT Bonus",
		}
		err := store.EditReward(rewardToUpdate)
		require.NoError(t, err)

		retrieved, err := store.GetRewardByID(createdRewardID)
		require.NoError(t, err)
		assert.Equal(t, "2 USDT Bonus", retrieved.RewardName)
		assert.Equal(t, CryptoUsdt1, retrieved.RewardType)
	})

	t.Run("GetReward", func(t *testing.T) {
		rewards, err := store.GetReward()
		require.NoError(t, err)
		assert.Len(t, rewards, 1)
	})

	t.Run("DeleteReward", func(t *testing.T) {
		err := store.DeleteReward(createdRewardID)
		require.NoError(t, err)

		_, err = store.GetRewardByID(createdRewardID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}
