package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"database/sql"
	"fmt"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	query := "SELECT id, username, email, password FROM users"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Error executing query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
			return nil, fmt.Errorf("Error scanning row: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over rows: %w", err)
	}

	return users, nil
}

func (r *UserRepository) FindByID(id uint) (models.User, error) {
	var user models.User
	query := "SELECT id, username, email, password FROM users WHERE id = ?"
	err := db.DB.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil // หากไม่พบข้อมูล
		}
		return user, fmt.Errorf("Error executing query: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
	_, err := db.DB.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("Error inserting user: %w", err)
	}
	return nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := "UPDATE users SET username = ?, email = ?, password = ? WHERE id = ?"
	_, err := db.DB.Exec(query, user.Username, user.Email, user.Password, user.ID)
	if err != nil {
		return fmt.Errorf("Error updating user: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(id uint) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := db.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("Error deleting user: %w", err)
	}
	return nil
}
