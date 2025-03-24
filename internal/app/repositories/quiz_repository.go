package repositories

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/utils"
	"fmt"
)

type QuizRepository struct{}

func NewQuizRepository() *QuizRepository { return &QuizRepository{} }

func (r *QuizRepository) Create(quiz *models.Quiz) error {

	id, err := utils.GenerateUniqueID(db.DB, "quizzes", "id")
	if err != nil {
		return fmt.Errorf("Error generating unique UUID: %w", err)
	}

	quiz.ID = id

	query := "INSERT INTO quizzes (id,course_id,number,quiz_type,title,lesson) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = db.DB.Exec(query, quiz.ID, quiz.CourseID, quiz.Number, quiz.QuizType, quiz.Title, quiz.Lesson)
	if err != nil {
		return fmt.Errorf("Create: error executing query: %w", err)
	}
	return nil
}

func (r *QuizRepository) GetAllQuizByCourseID(courseID string) ([]models.Quiz, error) {
	query := "SELECT id, course_id, number, quiz_type, title, lesson FROM quizzes WHERE course_id = ?"
	rows, err := db.DB.Query(query, courseID)
	if err != nil {
		return nil, fmt.Errorf("GetAllQuizByCourseID: error executing query: %w", err)
	}

	var quizzes []models.Quiz
	for rows.Next() {
		var quiz models.Quiz
		err = rows.Scan(&quiz.ID, &quiz.CourseID, &quiz.Number, &quiz.QuizType, &quiz.Title, &quiz.Lesson)
		if err != nil {
			return nil, fmt.Errorf("GetAllQuizByCourseID: error scanning row: %w", err)
		}
		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}
