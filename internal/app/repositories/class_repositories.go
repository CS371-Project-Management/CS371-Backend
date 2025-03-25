package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"database/sql"
	"errors"
)

type ClassRepository struct{}

func NewClassRepository() *ClassRepository {
	return &ClassRepository{}
}

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

// GetAllClasses - SELECT รายการคลาสทั้งหมด
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

// FindClassByID - SELECT class จาก id
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

// FindClassByInviteCode - ค้นหาคลาสตาม invite code
func (r *ClassRepository) FindClassByInviteCode(inviteCode string) (*models.Class, error) {
	query := `
		SELECT id, invite_code, title, description, accessibility
		FROM classes
		WHERE invite_code = ?
	`
	row := db.DB.QueryRow(query, inviteCode)
	var class models.Class
	err := row.Scan(&class.ID, &class.InviteCode, &class.Title, &class.Description, &class.Accessibility)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &class, nil
}

// UpdateClass - UPDATE ตาราง classes
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

// DeleteClass - DELETE จากตาราง classes
func (r *ClassRepository) DeleteClass(id uint) error {
	query := `DELETE FROM classes WHERE id = ?`
	_, err := db.DB.Exec(query, id)
	return err
}

// FindInviteCodeByID - ดึง invite_code จาก classes โดย id
func (r *ClassRepository) FindInviteCodeByID(id uint) (string, error) {
	var inviteCode string
	query := `SELECT invite_code FROM classes WHERE id = ?`
	err := db.DB.QueryRow(query, id).Scan(&inviteCode)
	if err != nil {
		return "", err
	}
	return inviteCode, nil
}

// InsertClassEnrollment - บันทึกการ join ลงใน pivot table user_classes
func (r *ClassRepository) InsertClassEnrollment(userID, classID uint) error {
	query := `INSERT INTO user_classes (id, user_id, class_id) VALUES (UUID(), ?, ?)`
	_, err := db.DB.Exec(query, userID, classID)
	return err
}

// HasJoinedClass - ตรวจสอบว่าผู้ใช้ได้ join คลาสนี้ไปแล้วหรือไม่
func (r *ClassRepository) HasJoinedClass(userID, classID uint) (bool, error) {
	query := `SELECT COUNT(*) FROM user_classes WHERE user_id = ? AND class_id = ?`
	var count int
	err := db.DB.QueryRow(query, userID, classID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}