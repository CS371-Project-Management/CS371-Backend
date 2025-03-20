package seeders

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	_ "database/sql"
	"log"
)

func SeedUsers() error {
	users := []models.User{
		{
			Username: "admin",
			Email:    "admin@example.com",
			Password: "hashed_password_here",
		},
		{
			Username: "user1",
			Email:    "user1@example.com",
			Password: "hashed_password_here",
		},
	}

	db := db.DB

	stmt, err := db.Prepare("INSERT IGNORE INTO users (username, email, password) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal("Error preparing SQL statement: ", err)
		return err
	}
	defer stmt.Close()

	for _, user := range users {
		_, err := stmt.Exec(user.Username, user.Email, user.Password)
		if err != nil {
			log.Printf("Error inserting user %s: %v", user.Username, err)
			return err
		}
		log.Printf("User created: %s\n", user.Username)
	}

	return nil
}
