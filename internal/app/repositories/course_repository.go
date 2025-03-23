package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log"
)

type CourseRepository struct {
}

func NewCourseRepository() *CourseRepository { return &CourseRepository{} }

func (r *CourseRepository) Create(course *models.CreateCourseRequest) error {
	id := uuid.New().String()
	query := "INSERT INTO courses (id ,class_id, number, title, description, difficulty_level) VALUES (?,?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(query, id, course.ClassID, course.Number, course.Title, course.Description, course.DifficultyLevel)
	if err != nil {
		return fmt.Errorf("Create: error executing query: %w", err)
	}
	return nil
}

func (r *CourseRepository) Update(course *models.UpdateCourseRequest) error {
	query := "UPDATE courses SET  title = ?, description = ?, difficulty_level = ? WHERE id = ?"
	_, err := db.DB.Exec(query, course.Title, course.Description, course.DifficultyLevel, course.ID)
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
