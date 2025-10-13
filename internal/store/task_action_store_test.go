package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDBAction(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	require.NoError(t, err, "opening test db")

	_, err = db.Exec("TRUNCATE task_actions CASCADE")
	require.NoError(t, err, "truncating table")

	return db
}

func TestTaskActionStore(t *testing.T) {
	db := setupTestDBAction(t)
	defer db.Close()

	store := NewPostgresTaskActionStore(db)

	var createdActionID int

	t.Run("CreateAction", func(t *testing.T) {
		action := &ActionTask{
			Type:        Type1,
			Name:        "Repost",
			Description: "Repost a tweet",
		}

		id, err := store.CreateAction(action)
		require.NoError(t, err)
		require.NotNil(t, id)
		require.Greater(t, *id, 0)
		createdActionID = *id
	})

	t.Run("GetActionByID", func(t *testing.T) {
		retrieved, err := store.GetActionByID(createdActionID)
		require.NoError(t, err)
		assert.Equal(t, createdActionID, retrieved.ID)
		assert.Equal(t, Type1, retrieved.Type)
		assert.Equal(t, "Repost", retrieved.Name)
	})

	t.Run("EditAction", func(t *testing.T) {
		actionToUpdate := &ActionTask{
			ID:          createdActionID,
			Name:        "Repost on X",
			Description: "Repost a post on X (formerly Twitter)",
		}
		err := store.EditAction(actionToUpdate)
		require.NoError(t, err)

		retrieved, err := store.GetActionByID(createdActionID)
		require.NoError(t, err)
		assert.Equal(t, "Repost on X", retrieved.Name)
		assert.Equal(t, "Repost a post on X (formerly Twitter)", retrieved.Description)
		assert.Equal(t, Type1, retrieved.Type)
	})

	t.Run("GetAction", func(t *testing.T) {
		actions, err := store.GetAction()
		require.NoError(t, err)
		assert.Len(t, actions, 1)
	})

	t.Run("DeleteAction", func(t *testing.T) {
		err := store.DeleteAction(createdActionID)
		require.NoError(t, err)

		_, err = store.GetActionByID(createdActionID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})
}
