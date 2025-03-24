package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"database/sql"
	"errors"
)

type ClassRepository struct {}

func NewClassRepository() *ClassRepository {
    return &ClassRepository{}
}

// func NewUserRepository() *UserRepository {
//     return &UserRepository{}
// }

// CreateClass - INSERT ลงตาราง classes
func (r *ClassRepository) CreateClass(class *models.Class) error {
    query := `
        INSERT INTO classes (invite_code, title, description, accessibility)
        VALUES (?, ?, ?, ?)
    `
    result, err := db.DB.Exec(query,
        class.InviteCode,
        class.Title,
        class.Description,
        class.Accessibility,
    )
    if err != nil {
        return err
    }

    lastID, err := result.LastInsertId()
    if err != nil {
        return err
    }
    class.ID = uint(lastID)
    return nil
}

// GetAllClasses ดึงรายการคลาสทั้งหมดจากตาราง classes
func (r *ClassRepository) GetAllClasses() ([]models.Class, error) {
	query := `
		SELECT id, invite_code, title, description, accessibility
		FROM classes
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		var c models.Class
		if err := rows.Scan(&c.ID, &c.InviteCode, &c.Title, &c.Description, &c.Accessibility); err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return classes, nil
}

func (r *ClassRepository) FindClassByID(id uint) (*models.Class, error) {
	query := `
		SELECT id, invite_code, title, description, accessibility
		FROM classes
		WHERE id = ?
	`
	row := db.DB.QueryRow(query, id)
	var class models.Class
	err := row.Scan(
		&class.ID,
		&class.InviteCode,
		&class.Title,
		&class.Description,
		&class.Accessibility,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &class, nil
}

func (r *ClassRepository) UpdateClass(class *models.Class) error {
	query := `
    UPDATE classes
    SET invite_code = ?, title = ?, description = ?, accessibility = ?
    WHERE id = ?
	`
	_, err := db.DB.Exec(query,
		class.InviteCode,
		class.Title,
		class.Description,
		class.Accessibility,
		class.ID,
	)
	return err
}

func (r *ClassRepository) DeleteClass(id uint) error {
	query := `
		DELETE FROM classes
		WHERE id = ?
	`
	_, err := db.DB.Exec(query, id)
	return err
}

func (r *ClassRepository) FindInviteCodeByID(id uint) (string, error) {
	var inviteCode string
	query := `SELECT invite_code FROM classes WHERE id = ?`
	err := db.DB.QueryRow(query, id).Scan(&inviteCode)
	if err != nil {
		return "", err
	}
	return inviteCode, nil
}
