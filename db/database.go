package db

import (
	"cs371-backend/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetRequiredEnv("DB_USER"),
		config.GetRequiredEnv("DB_PASSWORD"),
		config.GetRequiredEnv("DB_HOST"),
		config.GetEnv("DB_PORT", "3306"),
		config.GetRequiredEnv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("Error connecting to database: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("Error pinging database: %w", err)
	}

	return nil
}
