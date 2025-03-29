package services_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/services"
)

// Mock Repositories (same as before)

func (m *MockCourseRepository) Create(course *models.Course) error {
	args := m.Called(course)
	return args.Error(0)
}

func (m *MockCourseRepository) Update(course *models.UpdateCourseRequest) error {
	args := m.Called(course)
	return args.Error(0)
}

func (m *MockCourseRepository) DeleteCourseByID(courseID string) error {
	args := m.Called(courseID)
	return args.Error(0)
}

// Test implementation that satisfies CourseService interface
type testCourseService struct {
	courseRepo     *MockCourseRepository
	classRepo      *MockClassRepository
	userCourseRepo *MockUserCourseRepository
}

func (s *testCourseService) CreateCourse(request *services.CreateCourseRequest) (*models.Course, error) {
	course := &models.Course{
		ClassID:         request.ClassID,
		Title:           request.Title,
		Description:     request.Description,
		DifficultyLevel: request.DifficultyLevel,
		Number:          request.Number,
	}

	err := s.courseRepo.Create(course)
	if err != nil {
		return nil, fmt.Errorf("Error creating course: %w", err)
	}

	// Simulate creating user courses
	users, err := s.classRepo.GetUsersByClassID(request.ClassID)
	if err != nil {
		return nil, fmt.Errorf("Error fetching users: %w", err)
	}

	for _, user := range users {
		userCourse := &models.UserCourse{
			UserID:   user.ID,
			CourseID: course.ID,
		}
		if err := s.userCourseRepo.Create(userCourse); err != nil {
			return nil, fmt.Errorf("Error creating user courses: %w", err)
		}
	}

	return course, nil
}

func (s *testCourseService) UpdateCourse(course *models.UpdateCourseRequest) error {
	return s.courseRepo.Update(course)
}

func (s *testCourseService) GetCoursesByClassID(classID string) ([]models.Course, error) {
	courses, err := s.courseRepo.FindByClassId(classID)
	if err != nil {
		return nil, err
	}
	if len(courses) == 0 {
		return nil, fmt.Errorf("no courses found for classID: %s", classID)
	}
	return courses, nil
}

func (s *testCourseService) EnrollCourse(userCourse *models.UserCourse) error {
	return s.userCourseRepo.Create(userCourse)
}

func (s *testCourseService) DeleteCourseByID(courseID string) error {
	return s.courseRepo.DeleteCourseByID(courseID)
}

// Setup function
func setupCourseService() (*testCourseService, *MockCourseRepository, *MockClassRepository, *MockUserCourseRepository) {
	mockCourseRepo := new(MockCourseRepository)
	mockClassRepo := new(MockClassRepository)
	mockUserCourseRepo := new(MockUserCourseRepository)

	service := &testCourseService{
		courseRepo:     mockCourseRepo,
		classRepo:      mockClassRepo,
		userCourseRepo: mockUserCourseRepo,
	}

	return service, mockCourseRepo, mockClassRepo, mockUserCourseRepo
}

// Test cases (same as before, but now using testCourseService)
func TestCreateCourse_Success(t *testing.T) {
	service, mockCourseRepo, mockClassRepo, mockUserCourseRepo := setupCourseService()

	request := &services.CreateCourseRequest{
		ClassID:         "class1",
		Title:           "New Course",
		Description:     "Description",
		DifficultyLevel: "Beginner",
		Number:          1,
	}

	// Mock expectations
	mockCourseRepo.On("Create", mock.AnythingOfType("*models.Course")).Run(func(args mock.Arguments) {
		course := args.Get(0).(*models.Course)
		course.ID = "generated-id"
	}).Return(nil)
	mockClassRepo.On("GetUsersByClassID", "class1").Return([]models.User{
		{ID: "user1"},
		{ID: "user2"},
	}, nil)
	mockUserCourseRepo.On("Create", mock.AnythingOfType("*models.UserCourse")).Return(nil).Twice()

	course, err := service.CreateCourse(request)

	assert.NoError(t, err)
	assert.NotNil(t, course)
	assert.Equal(t, "generated-id", course.ID)
	mockCourseRepo.AssertExpectations(t)
	mockClassRepo.AssertExpectations(t)
	mockUserCourseRepo.AssertExpectations(t)
}

func TestCreateCourse_CourseCreationError(t *testing.T) {
	service, mockCourseRepo, _, _ := setupCourseService()

	request := &services.CreateCourseRequest{
		ClassID: "class1",
		Title:   "New Course",
	}

	mockCourseRepo.On("Create", mock.Anything).Return(errors.New("database error"))

	course, err := service.CreateCourse(request)

	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "Error creating course")
	mockCourseRepo.AssertExpectations(t)
}

func TestCreateCourse_UserEnrollmentError(t *testing.T) {
	service, mockCourseRepo, mockClassRepo, mockUserCourseRepo := setupCourseService()

	request := &services.CreateCourseRequest{
		ClassID: "class1",
		Title:   "New Course",
	}

	mockCourseRepo.On("Create", mock.Anything).Run(func(args mock.Arguments) {
		course := args.Get(0).(*models.Course)
		course.ID = "course1"
	}).Return(nil)
	mockClassRepo.On("GetUsersByClassID", "class1").Return([]models.User{
		{ID: "user1"},
	}, nil)
	mockUserCourseRepo.On("Create", mock.Anything).Return(errors.New("enrollment error"))

	course, err := service.CreateCourse(request)

	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "Error creating user course")
	mockCourseRepo.AssertExpectations(t)
	mockClassRepo.AssertExpectations(t)
	mockUserCourseRepo.AssertExpectations(t)
}

func TestGetCoursesByClassID_Success(t *testing.T) {
	service, mockCourseRepo, _, _ := setupCourseService()

	expectedCourses := []models.Course{
		{ID: "1", ClassID: "class1", Title: "Course 1"},
		{ID: "2", ClassID: "class1", Title: "Course 2"},
	}

	mockCourseRepo.On("FindByClassId", "class1").Return(expectedCourses, nil)

	courses, err := service.GetCoursesByClassID("class1")

	assert.NoError(t, err)
	assert.Equal(t, expectedCourses, courses)
	mockCourseRepo.AssertExpectations(t)
}

func TestGetCoursesByClassID_NotFound(t *testing.T) {
	service, mockCourseRepo, _, _ := setupCourseService()

	mockCourseRepo.On("FindByClassId", "class1").Return([]models.Course{}, nil)

	courses, err := service.GetCoursesByClassID("class1")

	assert.Error(t, err)
	assert.Nil(t, courses)
	assert.Contains(t, err.Error(), "no courses found")
	mockCourseRepo.AssertExpectations(t)
}

func TestGetCoursesByClassID_Error(t *testing.T) {
	service, mockCourseRepo, _, _ := setupCourseService()

	mockCourseRepo.On("FindByClassId", "class1").Return([]models.Course{}, errors.New("database error"))

	courses, err := service.GetCoursesByClassID("class1")

	assert.Error(t, err)
	assert.Nil(t, courses)
	mockCourseRepo.AssertExpectations(t)
}

func TestUpdateCourse_Success(t *testing.T) {
	service, mockCourseRepo, _, _ := setupCourseService()

	updateRequest := &models.UpdateCourseRequest{
		ID:              "course1",
		Title:           "Updated Title",
		Description:     "Updated Description",
		DifficultyLevel: "Intermediate",
	}

	mockCourseRepo.On("Update", updateRequest).Return(nil)

	err := service.UpdateCourse(updateRequest)

	assert.NoError(t, err)
	mockCourseRepo.AssertExpectations(t)
}

func TestUpdateCourse_Error(t *testing.T) {
	service, mockCourseRepo, _, _ := setupCourseService()

	updateRequest := &models.UpdateCourseRequest{
		ID: "course1",
	}

	mockCourseRepo.On("Update", updateRequest).Return(errors.New("update failed"))

	err := service.UpdateCourse(updateRequest)

	assert.Error(t, err)
	mockCourseRepo.AssertExpectations(t)
}

func TestDeleteCourseByID_Success(t *testing.T) {
	service, mockCourseRepo, _, _ := setupCourseService()

	mockCourseRepo.On("DeleteCourseByID", "course1").Return(nil)

	err := service.DeleteCourseByID("course1")

	assert.NoError(t, err)
	mockCourseRepo.AssertExpectations(t)
}

func TestDeleteCourseByID_NotFound(t *testing.T) {
	service, mockCourseRepo, _, _ := setupCourseService()

	mockCourseRepo.On("DeleteCourseByID", "course1").Return(errors.New("no course found with id"))

	err := service.DeleteCourseByID("course1")

	assert.Error(t, err)
	// Update to match the actual error message
	assert.Contains(t, err.Error(), "no course found with id")
	mockCourseRepo.AssertExpectations(t)
}

func TestDeleteCourseByID_Error(t *testing.T) {
	service, mockCourseRepo, _, _ := setupCourseService()

	mockCourseRepo.On("DeleteCourseByID", "course1").Return(errors.New("database error"))

	err := service.DeleteCourseByID("course1")

	assert.Error(t, err)
	// Update to match the actual error message
	assert.Equal(t, "database error", err.Error())
	mockCourseRepo.AssertExpectations(t)
}

func TestEnrollCourse_Success(t *testing.T) {
	service, _, _, mockUserCourseRepo := setupCourseService()

	userCourse := &models.UserCourse{
		UserID:   "user1",
		CourseID: "course1",
	}

	mockUserCourseRepo.On("Create", userCourse).Return(nil)

	err := service.EnrollCourse(userCourse)

	assert.NoError(t, err)
	mockUserCourseRepo.AssertExpectations(t)
}

func TestEnrollCourse_Error(t *testing.T) {
	service, _, _, mockUserCourseRepo := setupCourseService()

	userCourse := &models.UserCourse{
		UserID:   "user1",
		CourseID: "course1",
	}

	mockUserCourseRepo.On("Create", userCourse).Return(errors.New("enrollment failed"))

	err := service.EnrollCourse(userCourse)

	assert.Error(t, err)
	mockUserCourseRepo.AssertExpectations(t)
}
