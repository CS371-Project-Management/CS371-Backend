package repositories

import (
    "database/sql"
    "errors"
    "fmt"
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
    return &UserRepository{}
}

// ข้อมูลผู้ใช้ทั้งหมด
func (r *UserRepository) FindAll() ([]models.User, error) {
    var users []models.User

    query := "SELECT id, username, email, password FROM users"
    rows, err := db.DB.Query(query)
    if err != nil {
        return nil, fmt.Errorf("FindAll: error executing query: %w", err)
    }
    defer rows.Close()

    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
            return nil, fmt.Errorf("FindAll: error scanning row: %w", err)
        }
        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("FindAll: error iterating rows: %w", err)
    }

    return users, nil
}

//ดึงข้อมูลผู้ใช้ด้วย ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
    var user models.User

    query := "SELECT id, username, email, password FROM users WHERE id = ?"
    err := db.DB.QueryRow(query, id).
        Scan(&user.ID, &user.Username, &user.Email, &user.Password)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            // ไม่พบข้อมูล
            return nil, nil
        }
        return nil, fmt.Errorf("FindByID: error executing query: %w", err)
    }

    return &user, nil
}

//เพิ่มผู้ใช้ใหม่
func (r *UserRepository) Create(user *models.User) error {
    query := "INSERT INTO users (username, email, password) VALUES (?, ?, ?)"
    _, err := db.DB.Exec(query, user.Username, user.Email, user.Password)
    if err != nil {
        return fmt.Errorf("Create: error inserting user: %w", err)
    }
    return nil
}

// แก้ไขข้อมูล
func (r *UserRepository) Update(user *models.User) error {
    query := "UPDATE users SET username = ?, email = ?, password = ? WHERE id = ?"
    _, err := db.DB.Exec(query, user.Username, user.Email, user.Password, user.ID)
    if err != nil {
        return fmt.Errorf("Update: error updating user: %w", err)
    }
    return nil
}

// ลบผู้ใช้ตาม ID
func (r *UserRepository) Delete(id uint) error {
    query := "DELETE FROM users WHERE id = ?"
    _, err := db.DB.Exec(query, id)
    if err != nil {
        return fmt.Errorf("Delete: error deleting user: %w", err)
    }
    return nil
}

// FindByUsername ดึงข้อมูลผู้ใช้จากตาราง users ด้วย username
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
    var user models.User
    query := "SELECT id, username, email, password FROM users WHERE username = ?"
    err := db.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
    if err != nil {
        if err == sql.ErrNoRows {
            // ไม่พบข้อมูล
            return nil, nil
        }
        return nil, fmt.Errorf("FindByUsername error: %w", err)
    }
    return &user, nil
}
