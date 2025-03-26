package services

import (
	"errors"
	"math/rand"
	"time"

	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/repositories"
	"cs371-backend/internal/app/repositories/quiz"
	"cs371-backend/internal/app/repositories/quiz_history"

	"github.com/google/uuid"
)

type ClassService struct {
	repo *repositories.ClassRepository
	quizHistoryRepository *quiz_history.QuizHistoryRepository
	courseRepository *repositories.CourseRepository
	quizRepository *quiz.QuizRepository
	userCourseRepository *repositories.UserCourseRepository
}

func NewClassService() *ClassService {
	return &ClassService{
		repo: repositories.NewClassRepository(),
		quizHistoryRepository: quiz_history.NewQuizHistoryRepository(),
		courseRepository: repositories.NewCourseRepository(),
		quizRepository: quiz.NewQuizRepository(),
		userCourseRepository: repositories.NewUserCourseRepository(),
	}
}

func (s *ClassService) CreateClass(class *models.Class) error {
    // ตรวจสอบ user_id ใน test
    if class.UserID == "" {
        return errors.New("user_id cannot be empty")
    }

    // ตรวจสอบ title
    if class.Title == "" {
        return errors.New("class title cannot be empty")
    }

    // ถ้าไม่มี ID ให้ generate
    if class.ID == "" {
        class.ID = uuid.New().String()
    }

    // generate invite code
    class.InviteCode = generateInviteCode(6)

    // บันทึกผ่าน repo
    return s.repo.CreateClass(class)
}

func (s *ClassService) GetAllClasses() ([]models.Class, error) {
	return s.repo.GetAllClasses()
}

func (s *ClassService) GetClassByID(id string) (*models.Class, error) {
	return s.repo.FindClassByID(id)
}

func (s *ClassService) UpdateClass(class *models.Class) error {
	if class.Title == "" {
		return errors.New("class title cannot be empty")
	}
	// Gen invite_code ใหม่ทุกครั้งที่ update (ถ้าต้องการ)
	class.InviteCode = generateInviteCode(6)
	return s.repo.UpdateClass(class)
}

func (s *ClassService) DeleteClass(id string) error {
	return s.repo.DeleteClass(id)
}

func (s *ClassService) GetInviteCode(classID uint) (string, error) {
	inviteCode, err := s.repo.FindInviteCodeByID(classID)
	if err != nil {
		return "", err
	}
	if inviteCode == "" {
		return "", errors.New("class not found")
	}
	return inviteCode, nil
}

func (s *ClassService) JoinPublicClass(userID, classID string) error {
    class, err := s.repo.FindClassByID(classID)
    if err != nil {
        return err
    }
    if class == nil {
        return errors.New("class not found")
    }

    joined, err := s.repo.HasJoinedClass(userID, classID)
    if err != nil {
        return err
    }
    if joined {
        return errors.New("already joined this class")
    }
	
	s.repo.InsertClassEnrollment(userID, classID)

    var courseS []models.Course

	courseS,_ = s.courseRepository.FindByClassId(classID)
	for _,course := range courseS{
		user_course := new(models.UserCourse)
		user_course.CourseID = course.ID
		user_course.UserID = userID
		s.userCourseRepository.Create(user_course)
	}

	return nil
}

func (s *ClassService) JoinPrivateClass(userID string, inviteCode string) error {
	// ค้นหาคลาสที่มี invite code ตรงกับที่ผู้ใช้ส่งมา
	class, err := s.repo.FindClassByInviteCode(inviteCode)
	classID := class.ID
    if err != nil {
        return err
    }
    if class == nil {
        return errors.New("class not found")
    }

    joined, err := s.repo.HasJoinedClass(userID, classID)
    if err != nil {
        return err
    }
    if joined {
        return errors.New("already joined this class")
    }
	
	s.repo.InsertClassEnrollment(userID, classID)

    var courseS []models.Course

	courseS,_ = s.courseRepository.FindByClassId(classID)
	for _,course := range courseS{
		user_course := new(models.UserCourse)
		user_course.CourseID = course.ID
		user_course.UserID = userID
		s.userCourseRepository.Create(user_course)
	}

	return nil
}

// generateInviteCode สร้าง invite code แบบสุ่มความยาว length
func generateInviteCode(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// Exported function สำหรับ Controller ที่ต้องการ gen invite code ใหม่
func GenerateInviteCode(length int) string {
	return generateInviteCode(length)
}

// LeaveClass - กระบวนการออกจากคลาส
func (s *ClassService) LeaveClass(userID , classID string) error {
	err := s.repo.DeleteMemberFromClass(userID,classID )
	if err != nil {
		return err
	}
    return nil
}

// RemoveUserFromClass - บังคับให้ลบผู้ใช้จากคลาส
func (s *ClassService) RemoveUserFromClass(userID, classID string) error {
    return s.LeaveClass(userID, classID)
}