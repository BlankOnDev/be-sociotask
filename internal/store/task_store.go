package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	UserID         int64     `json:"user_id"`
	RewardID       int       `json:"reward_task"`
	RewardUSDT     float64   `json:"reward_usdt"` // total reward kah?
	DueDate        time.Time `json:"due_date"`
	MaxParticipant string    `json:"max_participant"`
	CreatedAt      time.Time `json:"created_at"`
	TaskImage      string    `json:"task_image"`
	ActionID       int       `json:"action_id"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PostgresTaskStore struct {
	db *sql.DB
}

func NewPostgresTaskStore(db *sql.DB) *PostgresTaskStore {
	return &PostgresTaskStore{db: db}
}

type TaskStore interface {
	CreateTask(task *Task) (*Task, error)
	GetAllTask(limit, offset int64) ([]Task, int64, error)
	GetTaskByID(id int64) (*Task, error)
	EditTask(t *Task) error
	DeleteTask(id int64) error
}

func (pg *PostgresTaskStore) CreateTask(task *Task) (*Task, error) {
	if task.UserID == 0 {
		return nil, errors.New("user id is required and can not be zero")
	}

	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `
	INSERT INTO tasks (
		title, 
		description, 
		user_id, 
		reward_id, 
		reward_usdt, 
		due_date, 
		max_participant, 
		task_image, 
		action_id
	) 
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9
	)
	RETURNING id
`

	err = tx.QueryRow(query, task.Title, task.Description, task.UserID, task.RewardID, task.RewardUSDT, task.DueDate, task.MaxParticipant, task.TaskImage, task.ActionID).Scan(&task.ID)
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
		SELECT 
			t.id, 
			t.title, 
			t.description, 
			t.user_id, 
			t.reward_id,
			r.reward_type,
			r.reward_name,
			t.reward_usdt, 
			t.due_date, 
			t.max_participant, 
			t.task_image, 
			t.action_id, 
			a.type AS action_type,
			a.name AS action_name,
			a.description AS action_description
		FROM tasks t
		LEFT JOIN task_rewards r ON t.reward_id = r.id
		LEFT JOIN task_actions a ON t.action_id = a.id
		WHERE t.id = $1
	`

	err = tx.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.UserID,
		&task.RewardID,
		&task.RewardUSDT,
		&task.DueDate,
		&task.MaxParticipant,
		&task.TaskImage,
		&task.ActionID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (pg *PostgresTaskStore) GetAllTask(limit, offset int64) ([]Task, int64, error) {
	cQuery := `
		SELECT COUNT(*)
		FROM tasks 
	`
	var pageNum int64
	err := pg.db.QueryRow(cQuery).Scan(&pageNum)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT 
			id, 
			title, 
			description, 
			user_id, 
			reward_id,
			reward_usdt, 
			due_date, 
			max_participant, 
			task_image, 
			action_id 
		FROM tasks 
		LIMIT $1 OFFSET $2
	`

	rows, err := pg.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID,
			&t.Title,
			&t.Description,
			&t.UserID,
			&t.RewardID,
			&t.RewardUSDT,
			&t.DueDate,
			&t.MaxParticipant,
			&t.TaskImage,
			&t.ActionID); err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, t)
	}
	defer rows.Close()

	return tasks, pageNum, nil
}

func (pg *PostgresTaskStore) EditTask(t *Task) error {
	var args []interface{}
	var setClause []string
	argCount := 1

	if t.Title != "" {
		setClause = append(setClause, fmt.Sprintf("title = $%d", argCount))
		args = append(args, t.Title)
		argCount++
	}

	if t.Description != "" {
		setClause = append(setClause, fmt.Sprintf("description = $%d", argCount))
		args = append(args, t.Description)
		argCount++
	}

	if t.RewardID != 0 {
		setClause = append(setClause, fmt.Sprintf("reward_id = $%d", argCount))
		args = append(args, t.RewardID)
		argCount++
	}

	if t.RewardUSDT != 0 {
		setClause = append(setClause, fmt.Sprintf("reward_usdt = $%d", argCount))
		args = append(args, t.RewardUSDT)
		argCount++
	}

	if !t.DueDate.IsZero() {
		setClause = append(setClause, fmt.Sprintf("due_date = $%d", argCount))
		args = append(args, t.DueDate)
		argCount++
	}

	if t.MaxParticipant != "" {
		setClause = append(setClause, fmt.Sprintf("max_participant = $%d", argCount))
		args = append(args, t.MaxParticipant)
		argCount++
	}

	if t.TaskImage != "" {
		setClause = append(setClause, fmt.Sprintf("task_image = $%d", argCount))
		args = append(args, t.TaskImage)
		argCount++
	}

	if t.ActionID != 0 {
		setClause = append(setClause, fmt.Sprintf("action_id = $%d", argCount))
		args = append(args, t.ActionID)
		argCount++
	}

	if len(setClause) == 0 {
		return fmt.Errorf("no fields to update for task id %d", t.ID)
	}

	setClause = append(setClause, "updated_at = NOW()")
	query := fmt.Sprintf(`
		UPDATE tasks
		SET %s
		WHERE id = $%d
	`, strings.Join(setClause, ", "), argCount)

	args = append(args, t.ID)

	result, err := pg.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", t.ID)
	}

	return nil

}

func (pg *PostgresTaskStore) DeleteTask(id int64) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}
