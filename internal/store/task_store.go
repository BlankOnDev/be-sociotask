package store

import (
	"database/sql"
	"errors"
)

type Task struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	UserID      int64   `json:"user_id"`
	RewardUSDT  float64 `json:"reward_usdt"`
}

type PostgresTaskStore struct {
	db *sql.DB
}

func NewPostgresTaskStore(db *sql.DB) *PostgresTaskStore {
	return &PostgresTaskStore{db: db}
}

type TaskStore interface {
	CreateTask(task *Task) (*Task, error)
	GetTaskByID(id int64) (*Task, error)
}

func (pg *PostgresTaskStore) CreateTask(task *Task) (*Task, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO tasks (title, description, user_id, reward_usdt)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err = tx.QueryRow(query, task.Title, task.Description, task.UserID, task.RewardUSDT).Scan(&task.ID)
	if task.UserID == 0 {
		return nil, errors.New("user id is required and can not be zero")
	}
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (pg *PostgresTaskStore) GetTaskByID(id int64) (*Task, error) {
	task := &Task{}

	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		SELECT id, title, description, user_id, reward_usdt
		FROM tasks
		WHERE id = $1
	`
	err = tx.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description, &task.UserID, &task.RewardUSDT)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return task, nil
}
