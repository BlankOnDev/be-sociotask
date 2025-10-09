package store

import (
	"database/sql"
	"fmt"
	"strings"
)

type TypeAction string

const (
	Type1 TypeAction = "type_1"
	Type2 TypeAction = "type_2"
	Type3 TypeAction = "type_3"
)

type ActionTask struct {
	ID          int        `json:"id"`
	Type        TypeAction `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
}

type PostgresTaskActionStore struct {
	db *sql.DB
}

func NewPostgresTaskActionStore(db *sql.DB) *PostgresTaskActionStore {
	return &PostgresTaskActionStore{db: db}
}

type TaskActionStore interface {
	CreateAction(req *ActionTask) (*int, error)
	EditAction(req *ActionTask) error
	DeleteAction(id int) error
	GetActionByID(id int) (*ActionTask, error)
	GetAction() ([]ActionTask, error)
}

func (pg *PostgresTaskActionStore) CreateAction(req *ActionTask) (*int, error) {
	query := `INSERT INTO task_actions (type, name, description) VALUES ($1, $2, $3) RETURNING id`
	err := pg.db.QueryRow(query, req.Type, req.Name, req.Description).Scan(&req.ID)
	if err != nil {
		return nil, err
	}

	return &req.ID, err
}

func (pg *PostgresTaskActionStore) EditAction(req *ActionTask) error {
	var args []interface{}
	var setClause []string
	argCount := 1

	if req.Type != "" {
		setClause = append(setClause, fmt.Sprintf("type = $%d", argCount))
		args = append(args, req.Type)
		argCount++
	}

	if req.Name != "" {
		setClause = append(setClause, fmt.Sprintf("name = $%d", argCount))
		args = append(args, req.Name)
		argCount++
	}

	if req.Description != "" {
		setClause = append(setClause, fmt.Sprintf("description = $%d", argCount))
		args = append(args, req.Description)
		argCount++
	}

	if len(setClause) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`UPDATE task_actions SET %s WHERE id = $%d`, strings.Join(setClause, ", "), argCount)
	args = append(args, req.ID)

	_, err := pg.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresTaskActionStore) DeleteAction(id int) error {
	query := `DELETE FROM task_actions WHERE id = $1`
	_, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresTaskActionStore) GetAction() ([]ActionTask, error) {
	query := `SELECT id, name, type, description FROM task_actions`

	rows, err := pg.db.Query(query)
	if err != nil {
		return nil, err
	}
	var resp []ActionTask
	for rows.Next() {
		var r ActionTask
		if err := rows.Scan(&r.ID, &r.Name, &r.Type, &r.Description); err != nil {
			return nil, err
		}
		resp = append(resp, r)
	}
	defer rows.Close()

	return resp, nil
}

func (pg *PostgresTaskActionStore) GetActionByID(id int) (*ActionTask, error) {
	query := `SELECT id, name, type, description FROM task_actions WHERE id = $1`

	var r ActionTask
	err := pg.db.QueryRow(query, id).Scan(&r.ID, &r.Name, &r.Type, &r.Description)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
