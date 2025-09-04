package store

import "database/sql"

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type PostgresTaskStore struct {
	db *sql.DB
}

func NewPostgresTaskStore(db *sql.DB) *PostgresTaskStore {
	return &PostgresTaskStore{db: db}
}

type TaskStore interface {
	CreateTask(task *Task) (*Task, error)
}

func (pg *PostgresTaskStore) CreateTask(task *Task) (*Task, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO tasks (title, description)
		VALUES ($1, $2)
		RETURNING id
	`

	err = tx.QueryRow(query, task.Title, task.Description).Scan(&task.ID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return task, nil
}
