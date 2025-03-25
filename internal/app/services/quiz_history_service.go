package services

import (
	quizmodel "cs371-backend/internal/app/models/quiz"
	quizhistorymodel "cs371-backend/internal/app/models/quiz_history"
	"cs371-backend/internal/app/repositories/quiz_history"
	"fmt"
)

type CreateQuizHistoryRequest struct {
	ID             string                            `json:"id,omitempty"`
	UserID         string                            `json:"user_id"`
	QuizID         string                            `json:"quiz_id"`
	QuizType       string                            `json:"quiz_type"`
	Note           string                            `json:"note"`
	Status         string                            `json:"status"`
	ChoiceAnswer   []CreateChoiceQuizHistoryData     `json:"choice_answer,omitempty"`
	OrderingAnswer []CreateOrderQuizHistoryData      `json:"ordering_answer,omitempty"`
	MissingAnswer  *CreateMissingWordQuizHistoryData `json:"missing_answer,omitempty"`
}

type CreateChoiceQuizHistoryData struct {
	Answer string `json:"answer"`
	Result bool   `json:"result"`
}

type CreateOrderQuizHistoryData struct {
	Answer string `json:"answer"`
	Order  int    `json:"order"`
	Result bool   `json:"result"`
}

type CreateMissingWordQuizHistoryData struct {
	Answer string `json:"answer"`
	Result bool   `json:"result"`
}

type QuizHistoryService struct {
	quizHistoryRepository            *quiz_history.QuizHistoryRepository
	choiceHistoryRepository          *quiz_history.ChoiceQuizHistoryRepository
	orderHistoryRepository           *quiz_history.OrderingQuizHistoryRepository
	missingWordQuizHistoryRepository *quiz_history.MissingWordQuizHistoryRepository
}

func NewQuizHistoryService() *QuizHistoryService {
	return &QuizHistoryService{
		quizHistoryRepository:            quiz_history.NewQuizHistoryRepository(),
		choiceHistoryRepository:          quiz_history.NewChoiceQuizHistoryRepository(),
		orderHistoryRepository:           quiz_history.NewOrderingQuizHistoryRepository(),
		missingWordQuizHistoryRepository: quiz_history.NewMissingWordQuizHistoryRepository(),
	}
}

func (s *QuizHistoryService) TakeQuiz(request *CreateQuizHistoryRequest) (*quizhistorymodel.QuizHistory, error) {
	// ลบประวัติเดิมก่อน (ถ้ามี)
	err := s.quizHistoryRepository.DeleteQuizHistoryByQuizID(request.QuizID, request.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete previous quiz history: %w", err)
	}

	quizType, err := quizmodel.ValidateQuizType(request.QuizType)
	if err != nil {
		return nil, err
	}

	quizStatus, err := quizhistorymodel.ValidateQuizStatus(request.Status)
	if err != nil {
		return nil, err
	}

	quizHistory := new(quizhistorymodel.QuizHistory)
	quizHistory.UserID = request.UserID
	quizHistory.QuizID = request.QuizID
	quizHistory.QuizType = quizType
	quizHistory.Note = request.Note
	quizHistory.Status = quizStatus

	err = s.quizHistoryRepository.CreateQuizHistory(quizHistory)
	if err != nil {
		return nil, err
	}

	switch quizType {
	case quizmodel.QuizTypeChoice:
		for _, answer := range request.ChoiceAnswer {
			choiceAnswer := new(quizhistorymodel.ChoiceHistory)
			choiceAnswer.QuizHistoryID = quizHistory.ID // เปลี่ยนจาก request.QuizID เป็น quizHistory.ID
			choiceAnswer.Answer = answer.Answer
			choiceAnswer.Result = answer.Result

			err = s.choiceHistoryRepository.CreateChoiceHistory(choiceAnswer)
			if err != nil {
				return nil, err
			}
		}

	case quizmodel.QuizTypeOrdering:
		for _, answer := range request.OrderingAnswer {
			orderAnswer := new(quizhistorymodel.OrderingHistory)
			orderAnswer.QuizHistoryID = quizHistory.ID // เปลี่ยนจาก request.QuizID เป็น quizHistory.ID
			orderAnswer.Order = answer.Order
			orderAnswer.Answer = answer.Answer
			orderAnswer.Result = answer.Result

			err = s.orderHistoryRepository.CreateOrderingHistory(orderAnswer)
			if err != nil {
				return nil, err
			}
		}

	case quizmodel.QuizTypeMissingWords:
		missingWordAnswer := new(quizhistorymodel.MissingWordHistory)
		missingWordAnswer.QuizHistoryID = quizHistory.ID // เปลี่ยนจาก request.QuizID เป็น quizHistory.ID
		missingWordAnswer.Answer = request.MissingAnswer.Answer
		missingWordAnswer.Result = request.MissingAnswer.Result

		err = s.missingWordQuizHistoryRepository.CreateMissingWordHistory(missingWordAnswer)
		if err != nil {
			return nil, err
		}
	}

	return quizHistory, nil
}

func (s *QuizHistoryService) GetQuizHistoryByID(historyID string) (*CreateQuizHistoryRequest, error) {
	// ดึงข้อมูลประวัติหลัก
	history, err := s.quizHistoryRepository.GetQuizHistoryByID(historyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quiz history: %w", err)
	}

	// สร้าง response object
	response := &CreateQuizHistoryRequest{
		ID:       history.ID,
		UserID:   history.UserID,
		QuizID:   history.QuizID,
		QuizType: quizmodel.QuizTypeToString(history.QuizType),
		Note:     history.Note,
		Status:   quizhistorymodel.QuizStatusToString(history.Status),
	}

	// ดึงข้อมูลคำตอบตามประเภท
	switch history.QuizType {
	case quizmodel.QuizTypeChoice:
		answers, err := s.choiceHistoryRepository.GetChoiceHistoryByQuizHistoryID(history.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get choice answers: %w", err)
		}

		choiceAnswers := make([]CreateChoiceQuizHistoryData, 0)
		for _, a := range answers {
			choiceAnswers = append(choiceAnswers, CreateChoiceQuizHistoryData{
				Answer: a.Answer,
				Result: a.Result,
			})
		}
		response.ChoiceAnswer = choiceAnswers

	case quizmodel.QuizTypeOrdering:
		answers, err := s.orderHistoryRepository.GetOrderingHistoryByQuizHistoryID(history.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get ordering answers: %w", err)
		}

		orderingAnswers := make([]CreateOrderQuizHistoryData, 0)
		for _, a := range answers {
			orderingAnswers = append(orderingAnswers, CreateOrderQuizHistoryData{
				Answer: a.Answer,
				Order:  a.Order,
				Result: a.Result,
			})
		}
		response.OrderingAnswer = orderingAnswers

	case quizmodel.QuizTypeMissingWords:
		answer, err := s.missingWordQuizHistoryRepository.GetMissingWordHistoryByQuizHistoryID(history.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get missing word answer: %w", err)
		}

		response.MissingAnswer = &CreateMissingWordQuizHistoryData{
			Answer: answer.Answer,
			Result: answer.Result,
		}
	}

	return response, nil
}

func (s *QuizHistoryService) GetAllQuizHistoryByQuizID(quizID string) ([]CreateQuizHistoryRequest, error) {
	quizHistories, err := s.quizHistoryRepository.GetAllQuizHistoryByQuizID(quizID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quiz histories: %w", err)
	}

	var historyRequests []CreateQuizHistoryRequest

	for _, history := range quizHistories {
		request := CreateQuizHistoryRequest{
			ID:       history.ID,
			UserID:   history.UserID,
			QuizID:   history.QuizID,
			QuizType: quizmodel.QuizTypeToString(history.QuizType),
			Note:     history.Note,
			Status:   quizhistorymodel.QuizStatusToString(history.Status),
		}

		switch history.QuizType {
		case quizmodel.QuizTypeChoice:
			choiceAnswers, err := s.choiceHistoryRepository.GetChoiceHistoryByQuizHistoryID(history.ID)
			if err != nil {
				return nil, err
			}

			var choiceData []CreateChoiceQuizHistoryData
			for _, answer := range choiceAnswers {
				choiceData = append(choiceData, CreateChoiceQuizHistoryData{
					Answer: answer.Answer,
					Result: answer.Result,
				})
			}
			request.ChoiceAnswer = choiceData

		case quizmodel.QuizTypeOrdering:
			orderingAnswers, err := s.orderHistoryRepository.GetOrderingHistoryByQuizHistoryID(history.ID)
			if err != nil {
				return nil, err
			}

			var orderingData []CreateOrderQuizHistoryData
			for _, answer := range orderingAnswers {
				orderingData = append(orderingData, CreateOrderQuizHistoryData{
					Answer: answer.Answer,
					Order:  answer.Order,
					Result: answer.Result,
				})
			}
			request.OrderingAnswer = orderingData

		case quizmodel.QuizTypeMissingWords:
			missingWordAnswer, err := s.missingWordQuizHistoryRepository.GetMissingWordHistoryByQuizHistoryID(history.ID)
			if err != nil {
				return nil, err
			}

			request.MissingAnswer = &CreateMissingWordQuizHistoryData{
				Answer: missingWordAnswer.Answer,
				Result: missingWordAnswer.Result,
			}
		}

		historyRequests = append(historyRequests, request)
	}

	if len(historyRequests) == 0 {
		return nil, fmt.Errorf("no quiz histories found for quiz ID: %s", quizID)
	}

	return historyRequests, nil
}
