package quiz

import (
	"cs371-backend/db"
	"cs371-backend/internal/app/models/quiz"
	"cs371-backend/internal/app/utils"
	"fmt"
	"log"
)

type QuizRepository struct{}

func NewQuizRepository() *QuizRepository { return &QuizRepository{} }

func (r *QuizRepository) Create(quiz *quiz.Quiz) error {

	id, err := utils.GenerateUniqueID(db.DB, "quizzes", "id")
	if err != nil {
		return fmt.Errorf("Error generating unique UUID: %w", err)
	}

	quiz.ID = id

	query := "INSERT INTO quizzes (id,course_id,point,number,quiz_type,title,lesson) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err = db.DB.Exec(query, quiz.ID, quiz.CourseID, quiz.Point, quiz.Number, quiz.QuizType, quiz.Title, quiz.Lesson)
	if err != nil {
		return fmt.Errorf("Create: error executing query: %w", err)
	}
	return nil
}

func (r *QuizRepository) GetAllQuizByCourseID(courseID string) ([]quiz.Quiz, error) {
	query := "SELECT id, course_id,point, number, quiz_type, title, lesson FROM quizzes WHERE course_id = ?"
	rows, err := db.DB.Query(query, courseID)
	if err != nil {
		return nil, fmt.Errorf("GetAllQuizByCourseID: error executing query: %w", err)
	}

	var quizzes []quiz.Quiz
	for rows.Next() {
		var quiz quiz.Quiz
		err = rows.Scan(&quiz.ID, &quiz.CourseID, &quiz.Point, &quiz.Number, &quiz.QuizType, &quiz.Title, &quiz.Lesson)
		if err != nil {
			return nil, fmt.Errorf("GetAllQuizByCourseID: error scanning row: %w", err)
		}
		quizzes = append(quizzes, quiz)
	}

	return quizzes, nil
}

func (r *QuizRepository) DeleteQuizByID(quizID string) error {
	log.Println("Repository: Deleting quiz with id: ", quizID)

	var count int
	checkQuery := "SELECT COUNT(*) FROM quizzes WHERE id = ?"
	err := db.DB.QueryRow(checkQuery, quizID).Scan(&count)
	if err != nil {
		return fmt.Errorf("Error checking if quiz exists: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("no quiz found with id %s", quizID)
	}

	query := "DELETE FROM quizzes WHERE id = ?"
	_, err = db.DB.Exec(query, quizID)
	if err != nil {
		return fmt.Errorf("DeleteQuizByID: error executing query: %w", err)
	}

	return nil
}
