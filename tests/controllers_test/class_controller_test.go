package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"cs371-backend/internal/app/controllers"
	"cs371-backend/internal/app/models"
)

// MockClassService implements services.ClassService for testing
type MockClassService struct {
	mock.Mock
}

func (m *MockClassService) CreateClass(class *models.Class) error {
	args := m.Called(class)
	return args.Error(0)
}

func (m *MockClassService) GetAllClasses() ([]models.Class, error) {
	args := m.Called()
	return args.Get(0).([]models.Class), args.Error(1)
}

func (m *MockClassService) GetClassByID(id string) (*models.Class, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Class), args.Error(1)
}

func (m *MockClassService) UpdateClass(class *models.Class) error {
	args := m.Called(class)
	return args.Error(0)
}

func (m *MockClassService) GetUsersByClassID(classID string) ([]models.User, error) {
	args := m.Called(classID)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockClassService) GetClassesOwnedByUser(userID string) ([]models.Class, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Class), args.Error(1)
}

func (m *MockClassService) DeleteClass(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockClassService) GetInviteCode(id string) (string, error) {
	args := m.Called(id)
	return args.String(0), args.Error(1)
}

func (m *MockClassService) JoinPublicClass(userID, classID string) error {
	args := m.Called(userID, classID)
	return args.Error(0)
}

func (m *MockClassService) JoinPrivateClass(userID, inviteCode string) error {
	args := m.Called(userID, inviteCode)
	return args.Error(0)
}

func (m *MockClassService) LeaveClass(userID, classID string) error {
	args := m.Called(userID, classID)
	return args.Error(0)
}

func (m *MockClassService) RemoveUserFromClass(userID, classID string) error {
	args := m.Called(userID, classID)
	return args.Error(0)
}

func (m *MockClassService) GetClassUserJoinByUserID(userID string) ([]models.Class, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Class), args.Error(1)
}

// setupClassController initializes the Fiber app with mock service
func setupClassController() (*fiber.App, *MockClassService) {
	app := fiber.New()
	mockService := new(MockClassService)
	classController := controllers.NewClassController(mockService)
	// Add middleware to set user_id for all requests
	app.Use(func(c *fiber.Ctx) error {
		// For testing, we'll get user_id from header if present
		if userID := c.Get("user_id"); userID != "" {
			c.Locals("user_id", userID)
		}
		return c.Next()
	})

	// Register all routes
	app.Post("/classes", classController.CreateClassHandler)
	app.Get("/classes", classController.GetAllClassesHandler)
	app.Get("/classes/:id", classController.GetClassHandler)
	app.Put("/classes/:id", classController.UpdateClassHandler)
	app.Get("/classes/:id/users", classController.GetUsersByClassIDHandler)
	app.Get("/users/:user_id/classes", classController.GetOwnedClassesHandler)
	app.Delete("/classes/:id", classController.DeleteClassHandler)
	app.Get("/classes/:id/invite", classController.GetInviteCodeHandler)
	app.Post("/classes/:id/join", classController.JoinPublicClassHandler)
	app.Post("/classes/join", classController.JoinPrivateClassHandler)
	app.Post("/classes/:id/leave", classController.LeaveClassHandler)
	app.Delete("/classes/:class_id/users/:user_id", classController.RemoveUserFromClassHandler)
	app.Get("/classes/joined/:user_id", classController.GetClassUserJoinByUserIDHandler)

	return app, mockService
}

func testWithUserID(userID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	}
}

func executeAuthenticatedRequest(app *fiber.App, method, path string, body []byte, userID string) (*http.Response, error) {
	// Create the HTTP request
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Use app.Test() which handles the context creation internally
	_, err := app.Test(req)
	if err != nil {
		return nil, err
	}

	// For Fiber v2, we need to set the user_id in a different way
	// Since we can't directly access the context, we'll modify the handler to get user_id from header
	// Alternatively, we can set it in a middleware
	// This is a simplified approach for testing:
	req.Header.Set("user_id", userID)

	return app.Test(req)
}

// Test cases

func TestCreateClassHandler_Success(t *testing.T) {
	app, mockService := setupClassController()

	requestBody := map[string]string{
		"user_id":       "user1",
		"title":         "Math Class",
		"description":   "Basic Math",
		"accessibility": "true",
	}

	mockService.On("CreateClass", mock.AnythingOfType("*models.Class")).Run(func(args mock.Arguments) {
		class := args.Get(0).(*models.Class)
		class.ID = "class1"
		class.InviteCode = "INV123"
	}).Return(nil)

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/classes", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Class created successfully", response["message"])
	mockService.AssertExpectations(t)
}

func TestGetClassHandler_Success(t *testing.T) {
	app, mockService := setupClassController()

	expectedClass := &models.Class{
		ID:          "class1",
		Title:       "Math Class",
		Description: "Basic Math",
	}

	mockService.On("GetClassByID", "class1").Return(expectedClass, nil)

	req := httptest.NewRequest(http.MethodGet, "/classes/class1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var actualClass models.Class
	err = json.NewDecoder(resp.Body).Decode(&actualClass)
	assert.NoError(t, err)
	assert.Equal(t, *expectedClass, actualClass)
	mockService.AssertExpectations(t)
}

func TestUpdateClassHandler_Success(t *testing.T) {
	app, mockService := setupClassController()

	requestBody := map[string]string{
		"title":         "Advanced Math",
		"description":   "Updated content",
		"accessibility": "true",
	}

	mockService.On("UpdateClass", mock.AnythingOfType("*models.Class")).Return(nil)

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPut, "/classes/class1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Class updated successfully", response["message"])
	mockService.AssertExpectations(t)
}

func TestJoinPublicClassHandler_Success(t *testing.T) {
	app, mockService := setupClassController()

	mockService.On("JoinPublicClass", "user1", "class1").Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/classes/class1/join", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("user_id", "user1") // Set user_id in header

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Joined public class successfully", response["message"])
	mockService.AssertExpectations(t)
}

func TestLeaveClassHandler_Success(t *testing.T) {
	app, mockService := setupClassController()

	mockService.On("LeaveClass", "user1", "class1").Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/classes/class1/leave", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("user_id", "user1") // Set user_id in header

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Left class successfully", response["message"])
	mockService.AssertExpectations(t)
}

func TestRemoveUserFromClassHandler_Success(t *testing.T) {
	app, mockService := setupClassController()

	// Mock the actual call being made by the controller
	mockService.On("LeaveClass", "user2", "class1").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/classes/class1/users/user2", nil)
	req.Header.Set("user_id", "admin1") // Admin performing the action

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "User removed from class successfully", response["message"])

	mockService.AssertExpectations(t)
}

func TestGetClassUserJoinByUserIDHandler_Success(t *testing.T) {
	app, mockService := setupClassController()

	expectedClasses := []models.Class{
		{ID: "class1", Title: "Math"},
		{ID: "class2", Title: "Science"},
	}

	mockService.On("GetClassUserJoinByUserID", "user1").Return(expectedClasses, nil)

	req := httptest.NewRequest(http.MethodGet, "/classes/joined/user1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var actualClasses []models.Class
	err = json.NewDecoder(resp.Body).Decode(&actualClasses)
	assert.NoError(t, err)
	assert.Equal(t, expectedClasses, actualClasses)
	mockService.AssertExpectations(t)
}
