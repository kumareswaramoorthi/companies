package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/kumareswaramoorthi/companies/api/utils"
	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	connVerifyErr := db.Ping()
	if connVerifyErr != nil {
		return nil, connVerifyErr
	}

	return db, nil
}

func GetDBConfig() Config {
	return Config{
		Host:     utils.GetEnvVars("DB_HOST", "127.0.0.1"),
		Port:     utils.GetEnvVars("DB_PORT", "5432"),
		Username: utils.GetEnvVars("DB_USER", "postgres"),
		Password: utils.GetEnvVars("DB_PWD", "password"),
		DBName:   utils.GetEnvVars("DB_NAME", "postgres"),
		SSLMode:  utils.GetEnvVars("DB_SSL_MODE", "disable"),
	}
}
