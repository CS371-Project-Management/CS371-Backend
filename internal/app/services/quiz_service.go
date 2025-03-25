package services

import (
	quiz2 "cs371-backend/internal/app/models/quiz"
	"cs371-backend/internal/app/repositories/quiz"
	"errors"
	"fmt"
	"log"
	"strings"
)

type CreateQuizRequest struct {
	ID              string                     `json:"id,omitempty"`
	QuizType        string                     `json:"quiz_type"`
	CourseID        string                     `json:"course_id"`
	Point           int                        `json:"point"`
	Number          int                        `json:"number"`
	Title           string                     `json:"title"`
	Lesson          string                     `json:"lesson"`
	ChoiceData      *CreateChoiceQuizData      `json:"choice_data,omitempty"`
	OrderingData    *CreateOrderingQuizData    `json:"ordering_data,omitempty"`
	MissingWordData *CreateMissingWordQuizData `json:"missing_word_data,omitempty"`
}

type CreateChoiceQuizData struct {
	Question string         `json:"question"`
	Type     string         `json:"type"`
	Answers  []ChoiceAnswer `json:"answers"`
}

type ChoiceAnswer struct {
	Answer string `json:"answer"`
	Result bool   `json:"result"`
}

type CreateOrderingQuizData struct {
	Question string           `json:"question"`
	Answers  []OrderingAnswer `json:"answers"`
}

type OrderingAnswer struct {
	Answer string `json:"answer"`
	Order  int    `json:"order"`
}

type CreateMissingWordQuizData struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type QuizService struct {
	quizRepository         *quiz.QuizRepository
	choiceRepository       *quiz.ChoiceQuizRepository
	choiceAnswerRepository *quiz.ChoiceAnswerRepository
	orderQuizRepository    *quiz.OrderingQuizRepository
	orderAnswerRepository  *quiz.OrderingAnswerRepository
	missingWordRepository  *quiz.MissingWordQuizRepository
}

func NewQuizService() *QuizService {
	return &QuizService{
		quizRepository:         quiz.NewQuizRepository(),
		choiceRepository:       quiz.NewChoiceQuizRepository(),
		choiceAnswerRepository: quiz.NewChoiceAnswerRepository(),
		orderQuizRepository:    quiz.NewOrderingQuizRepository(),
		orderAnswerRepository:  quiz.NewOrderingAnswerRepository(),
		missingWordRepository:  quiz.NewMissingWordQuizRepository(),
	}
}

func (s *QuizService) CreateQuiz(request *CreateQuizRequest) (*quiz2.Quiz, error) {
	quizType, err := quiz2.ValidateQuizType(request.QuizType)
	if err != nil {
		return nil, err
	}

	quiz := new(quiz2.Quiz)
	quiz.CourseID = request.CourseID
	quiz.Point = request.Point
	quiz.Number = request.Number
	quiz.QuizType = quizType
	quiz.Title = request.Title
	quiz.Lesson = request.Lesson

	err = s.quizRepository.Create(quiz)
	if err != nil {
		return nil, err
	}

	switch quizType {
	case quiz2.QuizTypeChoice:
		choiceType, err := quiz2.ValidateChoiceType(request.ChoiceData.Type)
		if err != nil {
			return nil, err
		}

		choiceQuiz := new(quiz2.ChoiceQuiz)
		choiceQuiz.QuizID = quiz.ID
		choiceQuiz.Question = request.ChoiceData.Question
		choiceQuiz.Type = choiceType

		err = s.choiceRepository.CreateChoiceQuiz(choiceQuiz)
		if err != nil {
			return nil, err
		}

		for _, answer := range request.ChoiceData.Answers {
			choiceAnswer := new(quiz2.ChoiceAnswer)
			choiceAnswer.QuizID = quiz.ID
			choiceAnswer.Answer = answer.Answer
			choiceAnswer.Result = answer.Result

			err = s.choiceAnswerRepository.CreateChoiceAnswer(choiceAnswer)
			if err != nil {
				return nil, err
			}
		}
	case quiz2.QuizTypeOrdering:
		orderingQuiz := new(quiz2.OrderingQuiz)
		orderingQuiz.QuizID = quiz.ID
		orderingQuiz.Question = request.OrderingData.Question

		err = s.orderQuizRepository.CreateOrderingQuiz(orderingQuiz)
		if err != nil {
			return nil, err
		}

		for _, answer := range request.OrderingData.Answers {
			orderingAnswer := new(quiz2.OrderingAnswer)
			orderingAnswer.QuizID = quiz.ID
			orderingAnswer.Answer = answer.Answer
			orderingAnswer.Order = answer.Order

			err = s.orderAnswerRepository.CreateOrderingAnswer(orderingAnswer)
			if err != nil {
				return nil, err
			}
		}

	case quiz2.QuizTypeMissingWords:
		missingWordQuiz := new(quiz2.MissingWordQuiz)
		missingWordQuiz.QuizID = quiz.ID
		missingWordQuiz.Question = request.MissingWordData.Question
		missingWordQuiz.Answer = request.MissingWordData.Answer

		err = s.missingWordRepository.CreateMissingWordQuiz(missingWordQuiz)
		if err != nil {
			return nil, err
		}

	default:
		return nil, errors.New("unsupported quiz type")

	}
	return quiz, nil
}

func (s *QuizService) GetAllQuizByCourseID(courseID string) ([]CreateQuizRequest, error) {
	quizzes, err := s.quizRepository.GetAllQuizByCourseID(courseID)
	if err != nil {
		return nil, err
	}

	var quizRequests []CreateQuizRequest

	for _, quiz := range quizzes {
		switch quiz.QuizType {
		case quiz2.QuizTypeChoice:
			choiceQuiz, err := s.choiceRepository.GetChoiceByQuizID(quiz.ID)
			if err != nil {
				return nil, err
			}

			choiceAnswers, err := s.choiceAnswerRepository.GetChoiceAnswerByQuizID(quiz.ID)
			if err != nil {
				return nil, err
			}

			answers := make([]ChoiceAnswer, 0)
			for _, answer := range choiceAnswers {
				answers = append(answers, ChoiceAnswer{
					Answer: answer.Answer,
					Result: answer.Result,
				})
			}

			quizRequest := CreateQuizRequest{
				ID:       quiz.ID,
				QuizType: quiz2.QuizTypeToString(quiz.QuizType),
				Point:    quiz.Point,
				CourseID: quiz.CourseID,
				Number:   quiz.Number,
				Title:    quiz.Title,
				Lesson:   quiz.Lesson,
				ChoiceData: &CreateChoiceQuizData{
					Question: choiceQuiz.Question,
					Type:     quiz2.ChoiceTypeToString(choiceQuiz.Type),
					Answers:  answers,
				},
			}
			quizRequests = append(quizRequests, quizRequest)

		case quiz2.QuizTypeOrdering:
			orderingQuiz, err := s.orderQuizRepository.GetOrderingByQuizID(quiz.ID)
			if err != nil {
				return nil, err
			}

			orderingAnswers, err := s.orderAnswerRepository.GetOrderingAnswerByQuizID(quiz.ID)
			if err != nil {
				return nil, err
			}

			answers := make([]OrderingAnswer, 0)
			for _, answer := range orderingAnswers {
				answers = append(answers, OrderingAnswer{
					Answer: answer.Answer,
					Order:  answer.Order,
				})
			}

			quizRequest := CreateQuizRequest{
				ID:       quiz.ID,
				QuizType: quiz2.QuizTypeToString(quiz.QuizType),
				Point:    quiz.Point,
				CourseID: quiz.CourseID,
				Number:   quiz.Number,
				Title:    quiz.Title,
				Lesson:   quiz.Lesson,
				OrderingData: &CreateOrderingQuizData{
					Question: orderingQuiz.Question,
					Answers:  answers,
				},
			}
			quizRequests = append(quizRequests, quizRequest)

		case quiz2.QuizTypeMissingWords:
			missingWordQuiz, err := s.missingWordRepository.GetMissingWordByQuizID(quiz.ID)
			if err != nil {
				return nil, err
			}

			quizRequest := CreateQuizRequest{
				ID:       quiz.ID,
				QuizType: quiz2.QuizTypeToString(quiz.QuizType),
				Point:    quiz.Point,
				CourseID: quiz.CourseID,
				Number:   quiz.Number,
				Title:    quiz.Title,
				Lesson:   quiz.Lesson,
				MissingWordData: &CreateMissingWordQuizData{
					Question: missingWordQuiz.Question,
					Answer:   missingWordQuiz.Answer,
				},
			}
			quizRequests = append(quizRequests, quizRequest)

		default:
			return nil, fmt.Errorf("unsupported quiz type: %s", quiz.QuizType)
		}
		log.Println(len(quizRequests))
	}

	if len(quizRequests) > 0 {
		return quizRequests, nil
	}

	return nil, fmt.Errorf("no quizzes found for course ID: %s", courseID)
}

func (s *QuizService) DeleteQuizByID(quizID string) error {
	log.Println("Service: Deleting quiz with id: ", quizID)
	err := s.quizRepository.DeleteQuizByID(quizID)
	if err != nil {
		if strings.Contains(err.Error(), "no quiz found with id") {
			return fmt.Errorf("Quiz not found: %w", err)
		}
		return fmt.Errorf("Error deleting quiz: %w", err)
	}

	return nil
}
