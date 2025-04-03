package services_test

import (
	"errors"
	"github.com/google/uuid"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/services"
)

// Mock Repositories
func (m *MockClassRepository) GetAllClasses() ([]models.Class, error) {
	args := m.Called()
	return args.Get(0).([]models.Class), args.Error(1)
}

func (m *MockClassRepository) FindClassByID(id string) (*models.Class, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Class), args.Error(1)
}

func (m *MockClassRepository) FindClassByTitle(title string) (*models.Class, error) {
	args := m.Called(title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Class), args.Error(1)
}

func (m *MockClassRepository) GetClassesByUserID(userID string) ([]models.Class, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Class), args.Error(1)
}

func (m *MockClassRepository) UpdateClass(class *models.Class) error {
	args := m.Called(class)
	return args.Error(0)
}

func (m *MockClassRepository) DeleteClass(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockClassRepository) FindInviteCodeByID(classID string) (string, error) {
	args := m.Called(classID)
	return args.String(0), args.Error(1)
}

func (m *MockClassRepository) FindClassByInviteCode(inviteCode string) (*models.Class, error) {
	args := m.Called(inviteCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Class), args.Error(1)
}

func (m *MockClassRepository) HasJoinedClass(userID, classID string) (bool, error) {
	args := m.Called(userID, classID)
	return args.Bool(0), args.Error(1)
}

func (m *MockClassRepository) InsertClassEnrollment(userID, classID string) error {
	args := m.Called(userID, classID)
	return args.Error(0)
}

func (m *MockClassRepository) DeleteMemberFromClass(userID, classID string) error {
	args := m.Called(userID, classID)
	return args.Error(0)
}

func (m *MockClassRepository) GetClassUserJoinByUserID(userID string) ([]models.Class, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Class), args.Error(1)
}

// Test Service Implementation
type testClassService struct {
	classRepo       *MockClassRepository
	courseRepo      *MockCourseRepository
	userCourseRepo  *MockUserCourseRepository
	quizRepo        *MockQuizRepository
	quizHistoryRepo *MockQuizHistoryRepository
}

// Implement all ClassService methods
func (s *testClassService) CreateClass(class *models.Class) error {
	if class.UserID == "" {
		return errors.New("user_id cannot be empty")
	}
	if class.Title == "" {
		return errors.New("class title cannot be empty")
	}

	existing, err := s.classRepo.FindClassByTitle(class.Title)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("title already exists")
	}

	if class.ID == "" {
		class.ID = uuid.New().String()
	}
	class.InviteCode = "MOCKINVITE"

	return s.classRepo.CreateClass(class)
}

func (s *testClassService) GetAllClasses() ([]models.Class, error) {
	return s.classRepo.GetAllClasses()
}

func (s *testClassService) GetClassByID(id string) (*models.Class, error) {
	return s.classRepo.FindClassByID(id)
}

func (s *testClassService) GetUsersByClassID(classID string) ([]models.User, error) {
	return s.classRepo.GetUsersByClassID(classID)
}

func (s *testClassService) GetClassesOwnedByUser(userID string) ([]models.Class, error) {
	return s.classRepo.GetClassesByUserID(userID)
}

func (s *testClassService) UpdateClass(class *models.Class) error {
	if class.Title == "" {
		return errors.New("class title cannot be empty")
	}
	class.InviteCode = "MOCKINVITE"
	return s.classRepo.UpdateClass(class)
}

func (s *testClassService) DeleteClass(id string) error {
	return s.classRepo.DeleteClass(id)
}

func (s *testClassService) GetInviteCode(classID string) (string, error) {
	inviteCode, err := s.classRepo.FindInviteCodeByID(classID)
	if err != nil {
		return "", err
	}
	if inviteCode == "" {
		return "", errors.New("class not found")
	}
	return inviteCode, nil
}

func (s *testClassService) JoinPublicClass(userID, classID string) error {
	class, err := s.classRepo.FindClassByID(classID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("class not found")
	}

	joined, err := s.classRepo.HasJoinedClass(userID, classID)
	if err != nil {
		return err
	}
	if joined {
		return errors.New("already joined this class")
	}

	if err := s.classRepo.InsertClassEnrollment(userID, classID); err != nil {
		return err
	}

	courses, err := s.courseRepo.FindByClassId(classID)
	if err != nil {
		return err
	}

	for _, course := range courses {
		userCourse := &models.UserCourse{
			UserID:   userID,
			CourseID: course.ID,
		}
		if err := s.userCourseRepo.Create(userCourse); err != nil {
			return err
		}
	}

	return nil
}

func (s *testClassService) JoinPrivateClass(userID, inviteCode string) error {
	class, err := s.classRepo.FindClassByInviteCode(inviteCode)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("class not found")
	}

	joined, err := s.classRepo.HasJoinedClass(userID, class.ID)
	if err != nil {
		return err
	}
	if joined {
		return errors.New("already joined this class")
	}

	if err := s.classRepo.InsertClassEnrollment(userID, class.ID); err != nil {
		return err
	}

	courses, err := s.courseRepo.FindByClassId(class.ID)
	if err != nil {
		return err
	}

	for _, course := range courses {
		userCourse := &models.UserCourse{
			UserID:   userID,
			CourseID: course.ID,
		}
		if err := s.userCourseRepo.Create(userCourse); err != nil {
			return err
		}
	}

	return nil
}

func (s *testClassService) LeaveClass(userID, classID string) error {
	return s.classRepo.DeleteMemberFromClass(userID, classID)
}

func (s *testClassService) RemoveUserFromClass(userID, classID string) error {
	return s.LeaveClass(userID, classID)
}

func (s *testClassService) GetClassUserJoinByUserID(userID string) ([]models.Class, error) {
	return s.classRepo.GetClassUserJoinByUserID(userID)
}

func NewTestClassService() *testClassService {
	return &testClassService{
		classRepo:       new(MockClassRepository),
		courseRepo:      new(MockCourseRepository),
		userCourseRepo:  new(MockUserCourseRepository),
		quizRepo:        new(MockQuizRepository),
		quizHistoryRepo: new(MockQuizHistoryRepository),
	}
}

// Test Cases
func TestCreateClass_Success(t *testing.T) {
	service := NewTestClassService()

	testClass := &models.Class{
		UserID: "user1",
		Title:  "Math Class",
	}

	service.classRepo.On("FindClassByTitle", "Math Class").Return(nil, nil)
	service.classRepo.On("CreateClass", mock.AnythingOfType("*models.Class")).Return(nil)

	err := service.CreateClass(testClass)

	assert.NoError(t, err)
	service.classRepo.AssertExpectations(t)
}

func TestCreateClass_ValidationErrors(t *testing.T) {
	tests := []struct {
		name      string
		class     *models.Class
		wantError string
	}{
		{
			name:      "Empty UserID",
			class:     &models.Class{Title: "Math"},
			wantError: "user_id cannot be empty",
		},
		{
			name:      "Empty Title",
			class:     &models.Class{UserID: "user1"},
			wantError: "class title cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTestClassService()
			err := service.CreateClass(tt.class)
			assert.Error(t, err)
			assert.Equal(t, tt.wantError, err.Error())
		})
	}
}

func TestCreateClass_TitleExists(t *testing.T) {
	service := NewTestClassService()

	testClass := &models.Class{
		UserID: "user1",
		Title:  "Math Class",
	}

	existingClass := &models.Class{ID: "existing"}
	service.classRepo.On("FindClassByTitle", "Math Class").Return(existingClass, nil)

	err := service.CreateClass(testClass)

	assert.Error(t, err)
	assert.Equal(t, "title already exists", err.Error())
	service.classRepo.AssertExpectations(t)
}

func TestGetAllClasses_Success(t *testing.T) {
	service := NewTestClassService()

	expected := []models.Class{
		{ID: "1", Title: "Math"},
		{ID: "2", Title: "Science"},
	}

	service.classRepo.On("GetAllClasses").Return(expected, nil)

	result, err := service.GetAllClasses()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	service.classRepo.AssertExpectations(t)
}

func TestGetClassByID_Success(t *testing.T) {
	service := NewTestClassService()

	expected := &models.Class{ID: "1", Title: "Math"}
	service.classRepo.On("FindClassByID", "1").Return(expected, nil)

	result, err := service.GetClassByID("1")

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	service.classRepo.AssertExpectations(t)
}

func TestGetClassByID_NotFound(t *testing.T) {
	service := NewTestClassService()

	service.classRepo.On("FindClassByID", "1").Return(nil, nil)

	result, err := service.GetClassByID("1")

	assert.NoError(t, err)
	assert.Nil(t, result)
	service.classRepo.AssertExpectations(t)
}

func TestUpdateClass_Success(t *testing.T) {
	service := NewTestClassService()

	class := &models.Class{
		ID:    "1",
		Title: "Updated Math",
	}

	service.classRepo.On("UpdateClass", class).Return(nil)

	err := service.UpdateClass(class)

	assert.NoError(t, err)
	service.classRepo.AssertExpectations(t)
}

func TestDeleteClass_Success(t *testing.T) {
	service := NewTestClassService()

	service.classRepo.On("DeleteClass", "1").Return(nil)

	err := service.DeleteClass("1")

	assert.NoError(t, err)
	service.classRepo.AssertExpectations(t)
}

func TestJoinPublicClass_Success(t *testing.T) {
	service := NewTestClassService()

	classID := "class1"
	userID := "user1"
	mockClass := &models.Class{ID: classID}

	service.classRepo.On("FindClassByID", classID).Return(mockClass, nil)
	service.classRepo.On("HasJoinedClass", userID, classID).Return(false, nil)
	service.classRepo.On("InsertClassEnrollment", userID, classID).Return(nil)
	service.courseRepo.On("FindByClassId", classID).Return([]models.Course{{ID: "course1"}}, nil)
	service.userCourseRepo.On("Create", mock.AnythingOfType("*models.UserCourse")).Return(nil)

	err := service.JoinPublicClass(userID, classID)

	assert.NoError(t, err)
	service.classRepo.AssertExpectations(t)
	service.courseRepo.AssertExpectations(t)
	service.userCourseRepo.AssertExpectations(t)
}

func TestJoinPublicClass_AlreadyJoined(t *testing.T) {
	service := NewTestClassService()

	classID := "class1"
	userID := "user1"
	mockClass := &models.Class{ID: classID}

	service.classRepo.On("FindClassByID", classID).Return(mockClass, nil)
	service.classRepo.On("HasJoinedClass", userID, classID).Return(true, nil)

	err := service.JoinPublicClass(userID, classID)

	assert.Error(t, err)
	assert.Equal(t, "already joined this class", err.Error())
	service.classRepo.AssertExpectations(t)
}

func TestJoinPrivateClass_Success(t *testing.T) {
	service := NewTestClassService()

	inviteCode := "INVITE123"
	userID := "user1"
	mockClass := &models.Class{ID: "class1"}

	service.classRepo.On("FindClassByInviteCode", inviteCode).Return(mockClass, nil)
	service.classRepo.On("HasJoinedClass", userID, "class1").Return(false, nil)
	service.classRepo.On("InsertClassEnrollment", userID, "class1").Return(nil)
	service.courseRepo.On("FindByClassId", "class1").Return([]models.Course{{ID: "course1"}}, nil)
	service.userCourseRepo.On("Create", mock.AnythingOfType("*models.UserCourse")).Return(nil)

	err := service.JoinPrivateClass(userID, inviteCode)

	assert.NoError(t, err)
	service.classRepo.AssertExpectations(t)
	service.courseRepo.AssertExpectations(t)
	service.userCourseRepo.AssertExpectations(t)
}

func TestLeaveClass_Success(t *testing.T) {
	service := NewTestClassService()

	classID := "class1"
	userID := "user1"

	service.classRepo.On("DeleteMemberFromClass", userID, classID).Return(nil)

	err := service.LeaveClass(userID, classID)

	assert.NoError(t, err)
	service.classRepo.AssertExpectations(t)
}

func TestRemoveUserFromClass_Success(t *testing.T) {
	service := NewTestClassService()

	classID := "class1"
	userID := "user1"

	service.classRepo.On("DeleteMemberFromClass", userID, classID).Return(nil)

	err := service.RemoveUserFromClass(userID, classID)

	assert.NoError(t, err)
	service.classRepo.AssertExpectations(t)
}

func TestGetClassUserJoinByUserID_Success(t *testing.T) {
	service := NewTestClassService()

	userID := "user1"
	expected := []models.Class{
		{ID: "1", Title: "Math"},
	}

	service.classRepo.On("GetClassUserJoinByUserID", userID).Return(expected, nil)

	result, err := service.GetClassUserJoinByUserID(userID)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	service.classRepo.AssertExpectations(t)
}

func TestGetInviteCode_Success(t *testing.T) {
	service := NewTestClassService()
	classID := "class1"
	expectedCode := "INVITE123"

	service.classRepo.On("FindInviteCodeByID", classID).Return(expectedCode, nil)

	result, err := service.GetInviteCode(classID)

	assert.NoError(t, err)
	assert.Equal(t, expectedCode, result)
	service.classRepo.AssertExpectations(t)
}

func TestGetInviteCode_NotFound(t *testing.T) {
	service := NewTestClassService()
	classID := "class1"

	service.classRepo.On("FindInviteCodeByID", classID).Return("", nil)

	_, err := service.GetInviteCode(classID)

	assert.Error(t, err)
	assert.Equal(t, "class not found", err.Error())
	service.classRepo.AssertExpectations(t)
}

func TestGetInviteCode_RepositoryError(t *testing.T) {
	service := NewTestClassService()
	classID := "class1"
	expectedErr := errors.New("database error")

	service.classRepo.On("FindInviteCodeByID", classID).Return("", expectedErr)

	_, err := service.GetInviteCode(classID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	service.classRepo.AssertExpectations(t)
}

func TestGenerateInviteCode(t *testing.T) {
	// Test that the generated code has the correct length
	code := services.GenerateInviteCode(6)
	assert.Equal(t, 6, len(code))

	// Test that two generated codes are different (randomness)
	code1 := services.GenerateInviteCode(6)
	time.Sleep(1 * time.Nanosecond) // Ensure different seed
	code2 := services.GenerateInviteCode(6)
	assert.NotEqual(t, code1, code2)
}
