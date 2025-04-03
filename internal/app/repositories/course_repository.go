package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/utils"
	"database/sql"
	"fmt"
	"log"
)

type CourseRepository struct {
}

func NewCourseRepository() *CourseRepository { return &CourseRepository{} }

func (r *CourseRepository) Create(course *models.Course) error {
	id, err := utils.GenerateUniqueID(db.DB, "courses", "id")
	if err != nil {
		return fmt.Errorf("Error generating unique UUID: %w", err)
	}

	course.ID = id

	query := "INSERT INTO courses (id, class_id, number, title, description, difficulty_level) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = db.DB.Exec(query, course.ID, course.ClassID, course.Number, course.Title, course.Description, course.DifficultyLevel)
	if err != nil {
		return fmt.Errorf("Create: error executing query: %w", err)
	}
	return nil
}

func (r *CourseRepository) Update(course *models.UpdateCourseRequest) error {
	queryCheck := "SELECT COUNT(*) FROM courses WHERE id = ?"
	var count int
	err := db.DB.QueryRow(queryCheck, course.ID).Scan(&count)
	if err != nil {
		return fmt.Errorf("error checking if course exists: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("course not found with id %s", course.ID)
	}

	queryUpdate := "UPDATE courses SET title = ?, description = ?, difficulty_level = ? WHERE id = ?"
	_, err = db.DB.Exec(queryUpdate, course.Title, course.Description, course.DifficultyLevel, course.ID)
	if err != nil {
		return fmt.Errorf("Update: error executing query: %w", err)
	}

	return nil
}

func (r *CourseRepository) FindByClassId(classID string) ([]models.Course, error) {
	query := "SELECT id, class_id, number, title, description, difficulty_level FROM courses WHERE class_id = ?"
	rows, err := db.DB.Query(query, classID)
	if err != nil {
		return nil, fmt.Errorf("FindByClassID: error executing query: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("FindByClassID: error closing rows: %v", err)
		}
	}(rows)

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.ClassID, &course.Number, &course.Title, &course.Description, &course.DifficultyLevel); err != nil {
			return nil, fmt.Errorf("FindByClassID: error scanning row: %w", err)
		}
		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FindByClassID: error iterating rows: %w", err)
	}

	return courses, nil
}

func (r *CourseRepository) DeleteCourseByID(courseID string) error {
	log.Println("Repository: Deleting course with id: ", courseID)

	var count int
	checkQuery := "SELECT COUNT(*) FROM courses WHERE id = ?"
	err := db.DB.QueryRow(checkQuery, courseID).Scan(&count)
	if err != nil {
		return fmt.Errorf("Error checking if course exists: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("no course found with id %s", courseID)
	}

	query := "DELETE FROM courses WHERE id = ?"
	_, err = db.DB.Exec(query, courseID)
	if err != nil {
		return fmt.Errorf("DeleteCourseByID: error executing query: %w", err)
	}

	return nil
}

func (r *CourseRepository) GetCourseByID(courseID string) (*models.Course, error) {
	query := `
		SELECT id, class_id, number, title, description, difficulty_level
		FROM courses
		WHERE id = ?
	`
	row := db.DB.QueryRow(query, courseID)
	var course models.Course
	err := row.Scan(
		&course.ID,
		&course.ClassID,
		&course.Number,
		&course.Title,
		&course.Description,
		&course.DifficultyLevel,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("GetCourseByID: error scanning row: %w", err)
	}
	return &course, nil
}
