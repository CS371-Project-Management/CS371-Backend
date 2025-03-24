package services

import (
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/repositories"
	"errors"
)

type CreateQuizRequest struct {
	QuizType        string                     `json:"quiz_type"`
	CourseID        string                     `json:"course_id"`
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
	quizRepository         *repositories.QuizRepository
	choiceRepository       *repositories.ChoiceQuizRepository
	choiceAnswerRepository *repositories.ChoiceAnswerRepository
	orderQuizRepository    *repositories.OrderingQuizRepository
	orderAnswerRepository  *repositories.OrderingAnswerRepository
	missingWordRepository  *repositories.MissingWordQuizRepository
}

func NewQuizService() *QuizService {
	return &QuizService{
		quizRepository:         repositories.NewQuizRepository(),
		choiceRepository:       repositories.NewChoiceQuizRepository(),
		choiceAnswerRepository: repositories.NewChoiceAnswerRepository(),
		orderQuizRepository:    repositories.NewOrderingQuizRepository(),
		orderAnswerRepository:  repositories.NewOrderingAnswerRepository(),
		missingWordRepository:  repositories.NewMissingWordQuizRepository(),
	}
}

func (s *QuizService) CreateQuiz(request *CreateQuizRequest) (*models.Quiz, error) {
	quizType, err := models.ValidateQuizType(request.QuizType)
	if err != nil {
		return nil, err
	}

	quiz := new(models.Quiz)
	quiz.CourseID = request.CourseID
	quiz.Number = request.Number
	quiz.QuizType = quizType
	quiz.Title = request.Title
	quiz.Lesson = request.Lesson

	err = s.quizRepository.Create(quiz)
	if err != nil {
		return nil, err
	}

	switch quizType {
	case models.QuizTypeChoice:
		choiceType, err := models.ValidateChoiceType(request.ChoiceData.Type)
		if err != nil {
			return nil, err
		}

		choiceQuiz := new(models.ChoiceQuiz)
		choiceQuiz.QuizID = quiz.ID
		choiceQuiz.Question = request.ChoiceData.Question
		choiceQuiz.Type = choiceType

		err = s.choiceRepository.CreateChoiceQuiz(choiceQuiz)
		if err != nil {
			return nil, err
		}

		for _, answer := range request.ChoiceData.Answers {
			choiceAnswer := new(models.ChoiceAnswer)
			choiceAnswer.QuizID = quiz.ID
			choiceAnswer.Answer = answer.Answer
			choiceAnswer.Result = answer.Result

			err = s.choiceAnswerRepository.CreateChoiceAnswer(choiceAnswer)
			if err != nil {
				return nil, err
			}
		}
	case models.QuizTypeOrdering:
		orderingQuiz := new(models.OrderingQuiz)
		orderingQuiz.QuizID = quiz.ID
		orderingQuiz.Question = request.OrderingData.Question

		err = s.orderQuizRepository.CreateOrderingQuiz(orderingQuiz)
		if err != nil {
			return nil, err
		}

		for _, answer := range request.OrderingData.Answers {
			orderingAnswer := new(models.OrderingAnswer)
			orderingAnswer.QuizID = quiz.ID
			orderingAnswer.Answer = answer.Answer
			orderingAnswer.Order = answer.Order

			err = s.orderAnswerRepository.CreateOrderingAnswer(orderingAnswer)
			if err != nil {
				return nil, err
			}
		}

	case models.QuizTypeMissingWords:
		missingWordQuiz := new(models.MissingWordQuiz)
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
