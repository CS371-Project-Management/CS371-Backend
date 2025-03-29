// tests/mocks_test.go
package services_test

import (
	"cs371-backend/internal/app/models"
	"github.com/stretchr/testify/mock"
)

// Shared mock repositories

type MockClassRepository struct {
	mock.Mock
}

func (m *MockClassRepository) CreateClass(class *models.Class) error {
	args := m.Called(class)
	return args.Error(0)
}

// Implement all ClassRepository methods...
// [Copy all MockClassRepository methods from your existing test file]

type MockCourseRepository struct {
	mock.Mock
}

func (m *MockCourseRepository) FindByClassId(classID string) ([]models.Course, error) {
	args := m.Called(classID)
	return args.Get(0).([]models.Course), args.Error(1)
}

type MockUserCourseRepository struct {
	mock.Mock
}

func (m *MockUserCourseRepository) Create(userCourse *models.UserCourse) error {
	args := m.Called(userCourse)
	return args.Error(0)
}

type MockQuizRepository struct {
	mock.Mock
}

type MockQuizHistoryRepository struct {
	mock.Mock
}

func (m *MockClassRepository) GetUsersByClassID(classID string) ([]models.User, error) {
	args := m.Called(classID)
	return args.Get(0).([]models.User), args.Error(1)
}
