package db

import (
	"fmt"
	"log"

	"cs371-backend/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() error {
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		config.GetRequiredEnv("DB_USER"),
		config.GetRequiredEnv("DB_PASSWORD"),
		config.GetRequiredEnv("DB_HOST"),
		config.GetEnv("DB_PORT", "3306"),
		config.GetRequiredEnv("DB_NAME"),
	)

	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}

func RunMigrationsDown() error {
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?multiStatements=true",
		config.GetRequiredEnv("DB_USER"),
		config.GetRequiredEnv("DB_PASSWORD"),
		config.GetRequiredEnv("DB_HOST"),
		config.GetEnv("DB_PORT", "3306"),
		config.GetRequiredEnv("DB_NAME"),
	)

	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		return err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations down completed successfully")
	return nil
}
