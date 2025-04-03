package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"cs371-backend/internal/app/controllers"
	"cs371-backend/internal/app/models"
)

// Mock สำหรับ UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Login(username, password string) (string, string, error) {
	args := m.Called(username, password)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockUserService) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) UpdateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) GenerateResetPasswordToken(email string) (string, error) {
	args := m.Called(email)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) ResetPassword(tokenString, newPassword string) error {
	args := m.Called(tokenString, newPassword)
	return args.Error(0)
}

// setupApp สร้าง Fiber App และ MockUserService
func setupApp() (*fiber.App, *MockUserService) {
	app := fiber.New()
	mockService := new(MockUserService)
	userController := controllers.NewUserController(mockService)

	// ลงทะเบียน routes
	app.Post("/login", userController.LoginHandler)
	app.Post("/logout", userController.Logout)
	app.Get("/users", userController.GetAllUsers)
	app.Get("/users/:id", userController.GetUser)
	app.Post("/users", userController.CreateUser)
	app.Put("/users/:id", userController.UpdateUser)
	app.Delete("/users/:id", userController.DeleteUser)
	app.Post("/reset-password/request", userController.RequestResetPassword)
	app.Post("/reset-password", userController.ResetPassword)

	return app, mockService
}

// ทดสอบ Login สำเร็จ
func TestLoginHandler_Success(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	expectedUserID := "user-123"
	expectedToken := "jwt-token-example"

	mockService.On("Login", "testuser", "password").Return(expectedUserID, expectedToken, nil)

	// สร้าง request
	loginData := map[string]string{
		"username": "testuser",
		"password": "password",
	}
	jsonBody, _ := json.Marshal(loginData)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)

	assert.Equal(t, "Login successful", response["message"])
	assert.Equal(t, expectedToken, response["token"])
	assert.Equal(t, expectedUserID, response["user_id"])
	mockService.AssertExpectations(t)
}

// ทดสอบ Login ล้มเหลว - ผู้ใช้ไม่พบ
func TestLoginHandler_UserNotFound(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	mockService.On("Login", "nonexistent", "password").Return("", "", errors.New("user not found"))

	// สร้าง request
	loginData := map[string]string{
		"username": "nonexistent",
		"password": "password",
	}
	jsonBody, _ := json.Marshal(loginData)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "User not found")
	mockService.AssertExpectations(t)
}

// ทดสอบ Login ล้มเหลว - รหัสผ่านไม่ถูกต้อง
func TestLoginHandler_InvalidPassword(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	mockService.On("Login", "testuser", "wrongpass").Return("", "", errors.New("invalid password"))

	// สร้าง request
	loginData := map[string]string{
		"username": "testuser",
		"password": "wrongpass",
	}
	jsonBody, _ := json.Marshal(loginData)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid username or password")
	mockService.AssertExpectations(t)
}

// ทดสอบ Logout
func TestLogout(t *testing.T) {
	// Setup
	app, _ := setupApp()

	// สร้าง request
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Successfully logged out", response["message"])
}

// ทดสอบ GetAllUsers สำเร็จ
func TestGetAllUsers_Success(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	expectedUsers := []models.User{
		{ID: "user-1", Username: "user1", Email: "user1@example.com"},
		{ID: "user-2", Username: "user2", Email: "user2@example.com"},
	}

	mockService.On("GetAllUsers").Return(expectedUsers, nil)

	// สร้าง request
	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var actualUsers []models.User
	json.NewDecoder(resp.Body).Decode(&actualUsers)

	assert.Equal(t, expectedUsers, actualUsers)
	mockService.AssertExpectations(t)
}

// ทดสอบ GetAllUsers ล้มเหลว
func TestGetAllUsers_Error(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	mockService.On("GetAllUsers").Return([]models.User{}, errors.New("database error"))

	// สร้าง request
	req := httptest.NewRequest(http.MethodGet, "/users", nil)

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "database error")
	mockService.AssertExpectations(t)
}

// ทดสอบ GetUser สำเร็จ
func TestGetUser_Success(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	expectedUser := &models.User{ID: "user-1", Username: "user1", Email: "user1@example.com"}

	mockService.On("GetUserByID", "user-1").Return(expectedUser, nil)

	// สร้าง request
	req := httptest.NewRequest(http.MethodGet, "/users/user-1", nil)

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var actualUser models.User
	json.NewDecoder(resp.Body).Decode(&actualUser)

	assert.Equal(t, *expectedUser, actualUser)
	mockService.AssertExpectations(t)
}

// ทดสอบ GetUser ไม่พบผู้ใช้
func TestGetUser_NotFound(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	mockService.On("GetUserByID", "nonexistent").Return(nil, errors.New("user not found"))

	// สร้าง request
	req := httptest.NewRequest(http.MethodGet, "/users/nonexistent", nil)

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "User not found", response["error"])
	mockService.AssertExpectations(t)
}

// ทดสอบ CreateUser สำเร็จ
func TestCreateUser_Success(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	newUser := &models.User{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
	}

	mockService.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)

	// สร้าง request
	jsonBody, _ := json.Marshal(newUser)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockService.AssertExpectations(t)
}

// ทดสอบ CreateUser ล้มเหลว - username ซ้ำ
func TestCreateUser_DuplicateUsername(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	newUser := &models.User{
		Username: "existinguser",
		Email:    "user@example.com",
		Password: "password123",
	}

	mockService.On("CreateUser", mock.AnythingOfType("*models.User")).Return(errors.New("username already exists"))

	// สร้าง request
	jsonBody, _ := json.Marshal(newUser)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "username already exists", response["error"])
	mockService.AssertExpectations(t)
}

// ทดสอบ UpdateUser สำเร็จ
func TestUpdateUser_Success(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	updatedUser := &models.User{
		ID:       "user-1",
		Username: "updateduser",
		Email:    "updated@example.com",
	}

	mockService.On("UpdateUser", mock.AnythingOfType("*models.User")).Return(nil)

	// สร้าง request
	jsonBody, _ := json.Marshal(updatedUser)
	req := httptest.NewRequest(http.MethodPut, "/users/user-1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

// ทดสอบ DeleteUser สำเร็จ
func TestDeleteUser_Success(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	mockService.On("DeleteUser", "user-1").Return(nil)

	// สร้าง request
	req := httptest.NewRequest(http.MethodDelete, "/users/user-1", nil)

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	mockService.AssertExpectations(t)
}

// ทดสอบ RequestResetPassword สำเร็จ
func TestRequestResetPassword_Success(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	email := "user@example.com"
	expectedToken := "reset-token-123"

	mockService.On("GenerateResetPasswordToken", email).Return(expectedToken, nil)

	// สร้าง request
	reqData := map[string]string{"email": email}
	jsonBody, _ := json.Marshal(reqData)
	req := httptest.NewRequest(http.MethodPost, "/reset-password/request", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Reset password link sent", response["message"])
	mockService.AssertExpectations(t)
}

// ทดสอบ RequestResetPassword ล้มเหลว - ไม่พบอีเมล
func TestRequestResetPassword_EmailNotFound(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	email := "nonexistent@example.com"

	mockService.On("GenerateResetPasswordToken", email).Return("", errors.New("email not found"))

	// สร้าง request
	reqData := map[string]string{"email": email}
	jsonBody, _ := json.Marshal(reqData)
	req := httptest.NewRequest(http.MethodPost, "/reset-password/request", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "email not found", response["error"])
	mockService.AssertExpectations(t)
}

// ทดสอบ ResetPassword สำเร็จ
func TestResetPassword_Success(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	token := "valid-reset-token"
	newPassword := "newpassword123"

	mockService.On("ResetPassword", token, newPassword).Return(nil)

	// สร้าง request
	reqData := map[string]string{
		"token":       token,
		"newPassword": newPassword,
	}
	jsonBody, _ := json.Marshal(reqData)
	req := httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Password reset successful", response["message"])
	mockService.AssertExpectations(t)
}

// ทดสอบ ResetPassword ล้มเหลว - โทเค็นไม่ถูกต้อง
func TestResetPassword_InvalidToken(t *testing.T) {
	// Setup
	app, mockService := setupApp()
	token := "invalid-token"
	newPassword := "newpassword123"

	mockService.On("ResetPassword", token, newPassword).Return(errors.New("invalid or expired token"))

	// สร้าง request
	reqData := map[string]string{
		"token":       token,
		"newPassword": newPassword,
	}
	jsonBody, _ := json.Marshal(reqData)
	req := httptest.NewRequest(http.MethodPost, "/reset-password", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Execute
	resp, _ := app.Test(req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "invalid or expired token", response["error"])
	mockService.AssertExpectations(t)
}
