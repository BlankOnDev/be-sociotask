package store

import (
	"database/sql"
	"errors"
	"time"
)

type Reward struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	TaskID    int64     `json:"task_id"`
	CreatedAt time.Time `json:"created_at"`
}

type PostgresRewardsStore struct {
	db *sql.DB
}

func NewPostgresRewardsStore(db *sql.DB) *PostgresRewardsStore {
	return &PostgresRewardsStore{db: db}
}

type RewardsStore interface {
	Create(reward *Reward) (*Reward, error)
}

func (pg *PostgresRewardsStore) Create(reward *Reward) (*Reward, error) {
	if reward.UserID == 0 {
		return nil, errors.New("user id is required and cannot be zero")
	}
	if reward.TaskID == 0 {
		return nil, errors.New("task id is required and cannot be zero")
	}

	// todo next: check user and task existence first

	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	    INSERT INTO rewards (user_id, task_id)
	    VALUES ($1, $2)
	    RETURNING id, created_at
       `

	err = tx.QueryRow(query, reward.UserID, reward.TaskID).Scan(&reward.ID, &reward.CreatedAt)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return reward, nil
}
