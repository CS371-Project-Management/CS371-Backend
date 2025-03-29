package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"cs371-backend/internal/app/controllers"
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/services"
)

// MockCourseService for testing
type MockCourseService struct {
	mock.Mock
}

func (m *MockCourseService) GetCoursesByClassID(classID string) ([]models.Course, error) {
	args := m.Called(classID)
	return args.Get(0).([]models.Course), args.Error(1)
}

func (m *MockCourseService) CreateCourse(request *services.CreateCourseRequest) (*models.Course, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Course), args.Error(1)
}

func (m *MockCourseService) UpdateCourse(course *models.UpdateCourseRequest) error {
	args := m.Called(course)
	return args.Error(0)
}

func (m *MockCourseService) DeleteCourseByID(courseID string) error {
	args := m.Called(courseID)
	return args.Error(0)
}

// Add the missing EnrollCourse method
func (m *MockCourseService) EnrollCourse(userCourse *models.UserCourse) error {
	args := m.Called(userCourse)
	return args.Error(0)
}

// Helper function to setup Fiber app with mock service
func setupCourseController() (*fiber.App, *MockCourseService) {
	app := fiber.New()
	mockService := new(MockCourseService)
	courseController := controllers.NewCourseController(mockService)

	app.Get("/classes/:classID/courses", courseController.GetCoursesByClassID)
	app.Post("/courses", courseController.CreateCourse)
	app.Put("/courses/:course_id", courseController.UpdateCourse)
	app.Delete("/courses/:course_id", courseController.DeleteCourseByID)

	return app, mockService
}

// Test Cases

func TestGetCoursesByClassID_Success(t *testing.T) {
	app, mockService := setupCourseController()

	expectedCourses := []models.Course{
		{ID: "1", ClassID: "class1", Title: "Course 1"},
		{ID: "2", ClassID: "class1", Title: "Course 2"},
	}

	mockService.On("GetCoursesByClassID", "class1").Return(expectedCourses, nil)

	req := httptest.NewRequest(http.MethodGet, "/classes/class1/courses", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var actualCourses []models.Course
	err = json.NewDecoder(resp.Body).Decode(&actualCourses)
	assert.NoError(t, err)
	assert.Equal(t, expectedCourses, actualCourses)
	mockService.AssertExpectations(t)
}

func TestGetCoursesByClassID_NotFound(t *testing.T) {
	app, mockService := setupCourseController()

	mockService.On("GetCoursesByClassID", "nonexistent").Return([]models.Course{}, errors.New("no courses found"))

	req := httptest.NewRequest(http.MethodGet, "/classes/nonexistent/courses", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "no courses found")
	mockService.AssertExpectations(t)
}

func TestGetCoursesByClassID_BadRequest(t *testing.T) {
	app, mockService := setupCourseController()

	mockService.On("GetCoursesByClassID", "invalid").Return([]models.Course{}, errors.New("error executing query"))

	req := httptest.NewRequest(http.MethodGet, "/classes/invalid/courses", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestCreateCourse_Success(t *testing.T) {
	app, mockService := setupCourseController()

	courseRequest := &services.CreateCourseRequest{
		ClassID:         "class1",
		Title:           "New Course",
		Description:     "Description",
		DifficultyLevel: "Beginner",
		Number:          1,
	}

	expectedCourse := &models.Course{
		ID:              "new-course",
		ClassID:         "class1",
		Title:           "New Course",
		Description:     "Description",
		DifficultyLevel: "Beginner",
		Number:          1,
	}

	mockService.On("CreateCourse", courseRequest).Return(expectedCourse, nil)

	jsonBody, _ := json.Marshal(courseRequest)
	req := httptest.NewRequest(http.MethodPost, "/courses", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var actualCourse models.Course
	err = json.NewDecoder(resp.Body).Decode(&actualCourse)
	assert.NoError(t, err)
	assert.Equal(t, *expectedCourse, actualCourse)
	mockService.AssertExpectations(t)
}

func TestCreateCourse_InvalidBody(t *testing.T) {
	app, _ := setupCourseController()

	// Invalid JSON body
	req := httptest.NewRequest(http.MethodPost, "/courses", strings.NewReader("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestUpdateCourse_Success(t *testing.T) {
	app, mockService := setupCourseController()

	courseID := "course1"
	updateRequest := models.UpdateCourseRequest{
		Title:           "Updated Title",
		Description:     "Updated Description",
		DifficultyLevel: "Intermediate",
	}

	mockService.On("UpdateCourse", &models.UpdateCourseRequest{
		ID:              courseID,
		Title:           "Updated Title",
		Description:     "Updated Description",
		DifficultyLevel: "Intermediate",
	}).Return(nil)

	jsonBody, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest(http.MethodPut, "/courses/"+courseID, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var actualResponse models.UpdateCourseRequest
	err = json.NewDecoder(resp.Body).Decode(&actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, updateRequest.Title, actualResponse.Title)
	mockService.AssertExpectations(t)
}

func TestUpdateCourse_NotFound(t *testing.T) {
	app, mockService := setupCourseController()

	courseID := "nonexistent"
	updateRequest := models.UpdateCourseRequest{
		Title: "Updated Title",
	}

	mockService.On("UpdateCourse", mock.Anything).Return(errors.New("course not found"))

	jsonBody, _ := json.Marshal(updateRequest)
	req := httptest.NewRequest(http.MethodPut, "/courses/"+courseID, bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "course not found")
	mockService.AssertExpectations(t)
}

func TestDeleteCourseByID_Success(t *testing.T) {
	app, mockService := setupCourseController()

	courseID := "course1"
	mockService.On("DeleteCourseByID", courseID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/courses/"+courseID, nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Course deleted successfully", response["message"])
	mockService.AssertExpectations(t)
}

func TestDeleteCourseByID_NotFound(t *testing.T) {
	app, mockService := setupCourseController()

	courseID := "nonexistent"
	mockService.On("DeleteCourseByID", courseID).Return(errors.New("Course not found"))

	req := httptest.NewRequest(http.MethodDelete, "/courses/"+courseID, nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "Course not found")
	mockService.AssertExpectations(t)
}
