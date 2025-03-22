package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"fmt"
	"github.com/google/uuid"
)

type CourseRepository struct {
}

func NewCourseRepository() *CourseRepository { return &CourseRepository{} }

func (r *CourseRepository) CreateCourse(course *models.Course) error {
	course.ID = uuid.New().String()
	query := "INSERT INTO courses (id ,class_id, number, title, description, difficulty_level) VALUES (?,?, ?, ?, ?, ?)"
	_, err := db.DB.Exec(query, course.ID, course.ClassID, course.Number, course.Title, course.Description, course.DifficultyLevel)
	if err != nil {
		return fmt.Errorf("CreateCourse: error executing query: %w", err)
	}
	return nil
}
