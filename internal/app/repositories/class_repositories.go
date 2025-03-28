package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type ClassRepository struct{}

func NewClassRepository() *ClassRepository {
	return &ClassRepository{}
}

// CreateClass - INSERT ลงตาราง classes
func (r *ClassRepository) CreateClass(class *models.Class) error {
	// ถ้ายังไม่มีค่า ID ให้ gen ใหม่
	if class.ID == "" {
		class.ID = uuid.New().String()
	}

	query := `
        INSERT INTO classes (id, user_id, invite_code, title, description, accessibility)
        VALUES (?, ?, ?, ?, ?, ?)
    `
	_, err := db.DB.Exec(query,
		class.ID,
		class.UserID,
		class.InviteCode,
		class.Title,
		class.Description,
		class.Accessibility,
	)
	return err
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

func (r *ClassRepository) FindClassByTitle(title string) (*models.Class, error) {
	query := `
		SELECT id, user_id, invite_code, title, description, accessibility
		FROM classes
		WHERE title = ?
	`
	row := db.DB.QueryRow(query, title)
	var class models.Class
	err := row.Scan(
		&class.ID,
		&class.UserID,
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

// FindClassByID - SELECT class จาก id
func (r *ClassRepository) FindClassByID(id string) (*models.Class, error) {
	query := `
		SELECT id, user_id, invite_code, title, description, accessibility
		FROM classes
		WHERE id = ?
	`
	row := db.DB.QueryRow(query, id)
	var class models.Class
	err := row.Scan(
		&class.ID,
		&class.UserID,
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
func (r *ClassRepository) DeleteClass(id string) error {
	query := `DELETE FROM classes WHERE id = ?`
	_, err := db.DB.Exec(query, id)
	return err
}

// FindInviteCodeByID - ดึง invite_code จาก classes โดย id
func (r *ClassRepository) FindInviteCodeByID(id string) (string, error) {
	var inviteCode string
	query := `SELECT invite_code FROM classes WHERE id = ?`
	err := db.DB.QueryRow(query, id).Scan(&inviteCode)
	if err != nil {
		return "", err
	}
	return inviteCode, nil
}

// InsertClassEnrollment - บันทึกการ join ลงใน pivot table user_classes
func (r *ClassRepository) InsertClassEnrollment(userID, classID string) error {
	query := `INSERT INTO user_classes (id, user_id, class_id) VALUES (UUID(), ?, ?)`
	_, err := db.DB.Exec(query, userID, classID)
	return err
}

// HasJoinedClass - ตรวจสอบว่าผู้ใช้ได้ join คลาสนี้ไปแล้วหรือไม่
func (r *ClassRepository) HasJoinedClass(userID, classID string) (bool, error) {
	query := `SELECT COUNT(*) FROM user_classes WHERE user_id = ? AND class_id = ?`
	var count int
	err := db.DB.QueryRow(query, userID, classID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// DeleteClassEnrollment - ลบการ join ใน pivot table user_classes
func (r *ClassRepository) DeleteClassEnrollment(userID, classID string) error {
	query := `DELETE FROM user_classes WHERE user_id = ? AND class_id = ?`
	_, err := db.DB.Exec(query, userID, classID)
	return err
}

// DeleteClassInUserCourses - ลบการ join ในตาราง user_courses
func (r *ClassRepository) DeleteClassInUserCourses(userID, classID string) error {
	query := `DELETE FROM user_courses WHERE user_id = ? AND class_id = ?`
	_, err := db.DB.Exec(query, userID, classID)
	return err
}

// // ForceRemoveUserFromClass - บังคับให้ผู้ใช้ถูกลบออกจากคลาส
// func (r *ClassRepository) ForceRemoveUserFromClass(userID, classID string) error {
//     query := `DELETE FROM user_classes WHERE user_id = ? AND class_id = ?`
//     result, err := db.DB.Exec(query, userID, classID)
//     if err != nil {
//         return err
//     }

//     rowsAffected, err := result.RowsAffected()
//     if err != nil {
//         return err
//     }

//     if rowsAffected == 0 {
//         return errors.New("user not found in the class")
//     }

//     return nil
// }

func (r *ClassRepository) DeleteMemberFromClass(userID, classID string) error {
	log.Println(userID)
	log.Println(classID)
	query := `
		DELETE qh, uc, ucl
		FROM user_classes ucl
		JOIN user_courses uc ON uc.user_id = ucl.user_id
			LEFT JOIN courses c ON c.class_id = ucl.class_id
			LEFT JOIN quizzes q ON q.course_id = c.id
			LEFT JOIN quiz_histories qh ON qh.user_id = uc.user_id AND qh.quiz_id = q.id
		WHERE ucl.class_id = ? AND ucl.user_id = ?;`

	result, err := db.DB.Exec(query, classID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found in the class")
	}

	return nil
}

// GetClassesByUserID ดึงรายการคลาสทั้งหมดที่เจ้าของเป็น user ที่ระบุ
func (r *ClassRepository) GetClassesByUserID(userID string) ([]models.Class, error) {
	query := `
		SELECT id, user_id, invite_code, title, description, accessibility
		FROM classes
		WHERE user_id = ?
	`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching classes for user %s: %w", userID, err)
	}
	defer rows.Close()

	var classes []models.Class
	for rows.Next() {
		var c models.Class
		if err := rows.Scan(&c.ID, &c.UserID, &c.InviteCode, &c.Title, &c.Description, &c.Accessibility); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		classes = append(classes, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}
	return classes, nil
}

func (r *ClassRepository) GetUsersByClassID(classID string) ([]models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.first_name, u.last_name, u.role
		FROM user_classes uc
		JOIN users u ON u.id = uc.user_id
		WHERE uc.class_id = ?
	`
	rows, err := db.DB.Query(query, classID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching users for class: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.FirstName, &u.LastName, &u.Role); err != nil {
			return nil, fmt.Errorf("Error scanning user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over rows: %w", err)
	}

	return users, nil
}
