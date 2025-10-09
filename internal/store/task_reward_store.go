package store

import (
	"database/sql"
	"fmt"
	"strings"
)

type JenisCategory string

const (
	CryptoUsdt1 JenisCategory = "crypto_usdt_1"
	CryptoUsdt2 JenisCategory = "crypto_usdt_2"
	CryptoUsdt3 JenisCategory = "crypto_usdt_3"
)

type RewardTask struct {
	ID         int           `json:"id"`
	RewardType JenisCategory `json:"reward_type"`
	RewardName string        `json:"reward_name"`
}

type PostgresTaskRewardStore struct {
	db *sql.DB
}

func NewPostgresTaskRewardStore(db *sql.DB) *PostgresTaskRewardStore {
	return &PostgresTaskRewardStore{db: db}
}

type TaskRewardStore interface {
	CreateReward(req *RewardTask) (*int, error)
	EditReward(req *RewardTask) error
	DeleteReward(id int) error
	GetReward() ([]RewardTask, error)
	GetRewardByID(id int) (*RewardTask, error)
}

func (pg *PostgresTaskRewardStore) CreateReward(req *RewardTask) (*int, error) {
	query := `INSERT INTO task_rewards (reward_type, reward_name) VALUES ($1, $2) RETURNING id`
	err := pg.db.QueryRow(query, req.RewardType, req.RewardName).Scan(&req.ID)
	if err != nil {
		return nil, err
	}

	return &req.ID, err
}

func (pg *PostgresTaskRewardStore) EditReward(req *RewardTask) error {
	var args []interface{}
	var setClause []string
	argCount := 1

	if req.RewardType != "" {
		setClause = append(setClause, fmt.Sprintf("reward_type = $%d", argCount))
		args = append(args, req.RewardType)
		argCount++
	}

	if req.RewardName != "" {
		setClause = append(setClause, fmt.Sprintf("reward_name = $%d", argCount))
		args = append(args, req.RewardName)
		argCount++
	}

	if len(setClause) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`UPDATE task_rewards SET %s WHERE id = $%d`, strings.Join(setClause, ", "), argCount)
	args = append(args, req.ID)

	_, err := pg.db.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresTaskRewardStore) DeleteReward(id int) error {
	query := `DELETE FROM task_rewards WHERE id = $1`
	_, err := pg.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (pg *PostgresTaskRewardStore) GetReward() ([]RewardTask, error) {
	query := `SELECT id, reward_type, reward_name FROM task_rewards`
	rows, err := pg.db.Query(query)
	if err != nil {
		return nil, err
	}
	var resp []RewardTask
	for rows.Next() {
		var r RewardTask
		if err := rows.Scan(&r.ID, &r.RewardType, &r.RewardName); err != nil {
			return nil, err
		}
		resp = append(resp, r)
	}
	defer rows.Close()

	return resp, nil
}

func (pg *PostgresTaskRewardStore) GetRewardByID(id int) (*RewardTask, error) {
	query := `SELECT id, reward_type, reward_name FROM task_rewards WHERE id = $1`

	var r RewardTask
	err := pg.db.QueryRow(query, id).Scan(&r.ID, &r.RewardName, &r.RewardType)
	if err != nil {
		return nil, err
	}

	return &r, nil
}
