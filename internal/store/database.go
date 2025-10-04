package store

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/harundarat/be-socialtask/internal/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

func Open() (*sql.DB, error) {
	config := DatabaseConfig{
		Host:     utils.GetEnv("DB_HOST"),
		User:     utils.GetEnv("DB_USER"),
		Password: utils.GetEnv("DB_PASSWORD"),
		DBName:   utils.GetEnv("DB_NAME"),
		Port:     utils.GetEnv("DB_PORT"),
	}
	connectionString := config.ConnectionString()

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database")
	return db, nil
}

func MigrateFS(db *sql.DB, migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	return nil
}

func (dc *DatabaseConfig) ConnectionString() string {
	// return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dc.Host, dc.User, dc.Password, dc.DBName, dc.Port)
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require", dc.Host, dc.User, dc.Password, dc.DBName, dc.Port)
}
