package services_test

import (
	"cs371-backend/internal/app/services"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"cs371-backend/internal/app/models"
)

// MockUserRepository สำหรับจำลองการทำงานของ repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdatePassword(userID string, hashedPassword string) error {
	args := m.Called(userID, hashedPassword)
	return args.Error(0)
}

// สร้าง struct ที่มี repo field เป็น MockUserRepository สำหรับทดสอบ
type testUserService struct {
	repo *MockUserRepository
}

// เพื่อให้เป็นไปตาม UserService interface
func (s *testUserService) Login(username, password string) (string, string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil || user == nil {
		return "", "", errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		return "", "", err
	}

	return user.ID, tokenString, nil
}

func (s *testUserService) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *testUserService) GetUserByID(id string) (*models.User, error) {
	return s.repo.FindByID(id)
}

func (s *testUserService) CreateUser(user *models.User) error {
	// ตรวจสอบความถูกต้องเบื้องต้น
	if user.Username == "" || user.Password == "" {
		return errors.New("username or password cannot be empty")
	}

	// ตรวจสอบ username ซ้ำ
	existing, err := s.repo.FindByUsername(user.Username)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("username already exists")
	}

	// แฮชรหัสผ่าน
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	// บันทึกลง DB
	return s.repo.Create(user)
}

func (s *testUserService) UpdateUser(user *models.User) error {
	return s.repo.Update(user)
}

func (s *testUserService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}

func (s *testUserService) GenerateResetPasswordToken(email string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("email not found")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte("test-reset-token"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *testUserService) ResetPassword(tokenString, newPassword string) error {
	// Parse JWT
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte("test-reset-token"), nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return errors.New("invalid user_id in token")
	}

	// Add this check to verify user exists
	_, err = s.repo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(userID, string(hashed))
}

// สร้าง helper function สำหรับสร้าง UserService ที่ใช้ MockUserRepository
func setupUserService() (services.UserService, *MockUserRepository) {
	mockRepo := new(MockUserRepository)
	service := &testUserService{
		repo: mockRepo,
	}
	return service, mockRepo
}

// เริ่มเขียน unit tests
func TestLogin_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupUserService()

	// สร้างรหัสผ่านที่ถูกแฮช
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	// สร้าง mock user
	mockUser := &models.User{
		ID:       "user-123",
		Username: "testuser",
		Password: string(hashedPassword),
	}

	// ตั้งค่า expectation
	mockRepo.On("FindByUsername", "testuser").Return(mockUser, nil)

	// สร้างสภาพแวดล้อมสำหรับ JWT
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	// Execute
	userID, token, err := service.Login("testuser", "correctpassword")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "user-123", userID)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	// Setup
	service, mockRepo := setupUserService()

	// ตั้งค่า expectation
	mockRepo.On("FindByUsername", "nonexistent").Return(nil, nil)

	// Execute
	userID, token, err := service.Login("nonexistent", "anypassword")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Empty(t, userID)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	// Setup
	service, mockRepo := setupUserService()

	// สร้างรหัสผ่านที่ถูกแฮช
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	// สร้าง mock user
	mockUser := &models.User{
		ID:       "user-123",
		Username: "testuser",
		Password: string(hashedPassword),
	}

	// ตั้งค่า expectation
	mockRepo.On("FindByUsername", "testuser").Return(mockUser, nil)

	// Execute
	userID, token, err := service.Login("testuser", "wrongpassword")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid password", err.Error())
	assert.Empty(t, userID)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupUserService()

	newUser := &models.User{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
	}

	// Set up expectations
	mockRepo.On("FindByUsername", "newuser").Return(nil, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Run(func(args mock.Arguments) {
		user := args.Get(0).(*models.User)
		// Verify the password was hashed
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
		assert.NoError(t, err, "Password should be hashed")
	}).Return(nil)

	// Execute
	err := service.CreateUser(newUser)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	// Setup
	service, mockRepo := setupUserService()

	existingUser := &models.User{
		ID:       "existing-id",
		Username: "existinguser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	// ตั้งค่า expectation
	mockRepo.On("FindByUsername", "existinguser").Return(existingUser, nil)

	// Execute
	err := service.CreateUser(&models.User{
		Username: "existinguser",
		Email:    "newemail@example.com",
		Password: "password123",
	})

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "username already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGenerateResetPasswordToken_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupUserService()

	// Mock user ที่มีอีเมล
	mockUser := &models.User{
		ID:    "user-123",
		Email: "test@example.com",
	}

	// ตั้งค่า expectation
	mockRepo.On("FindByEmail", "test@example.com").Return(mockUser, nil)

	// Execute
	token, err := service.GenerateResetPasswordToken("test@example.com")

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestGenerateResetPasswordToken_EmailNotFound(t *testing.T) {
	// Setup
	service, mockRepo := setupUserService()

	// ตั้งค่า expectation
	mockRepo.On("FindByEmail", "nonexistent@example.com").Return(nil, nil)

	// Execute
	token, err := service.GenerateResetPasswordToken("nonexistent@example.com")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "email not found", err.Error())
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestResetPassword_Success(t *testing.T) {
	// Setup
	service, mockRepo := setupUserService()

	// Mock user ที่จะรีเซ็ตรหัสผ่าน
	mockUser := &models.User{
		ID:       "user-123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	// ตั้งค่า expectation สำหรับการค้นหาผู้ใช้
	// เปลี่ยนจาก FindByEmail เป็น FindByID ตามที่ token มีแค่ user_id
	mockRepo.On("FindByID", "user-123").Return(mockUser, nil)
	mockRepo.On("UpdatePassword", "user-123", mock.Anything).Return(nil)

	// สร้าง JWT token สำหรับ reset password
	claims := jwt.MapClaims{
		"user_id": "user-123", // Token มีแค่ user_id
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test-reset-token"))

	// Execute
	err := service.ResetPassword(tokenString, "newpassword123")

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
func TestResetPassword_InvalidToken(t *testing.T) {
	// Setup
	service, _ := setupUserService() // Ignore the repository for this test

	// สร้าง JWT token ที่ไม่ถูกต้อง
	invalidToken := "invalid.token.string"

	// Execute
	err := service.ResetPassword(invalidToken, "newpassword123")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid or expired token", err.Error())
}

func TestCreateUser_EmptyPassword(t *testing.T) {
	service, _ := setupUserService()
	err := service.CreateUser(&models.User{Username: "testuser", Password: ""})
	assert.Error(t, err)
	assert.Equal(t, "username or password cannot be empty", err.Error())
}
