package services

import (
	"errors"
	"math/rand"
	"time"

	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/repositories"
)

type ClassService struct {
	repo *repositories.ClassRepository
}

func NewClassService() *ClassService {
	return &ClassService{
		repo: repositories.NewClassRepository(),
	}
}

func (s *ClassService) CreateClass(class *models.Class) error {
	if class.Title == "" {
		return errors.New("class title cannot be empty")
	}
	if class.Accessibility == "" {
		class.Accessibility = "private"
	}
	// Gen invite_code อัตโนมัติ
	class.InviteCode = generateInviteCode(6)
	return s.repo.CreateClass(class)
}

func (s *ClassService) GetAllClasses() ([]models.Class, error) {
	return s.repo.GetAllClasses()
}

func (s *ClassService) GetClassByID(id uint) (*models.Class, error) {
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

func (s *ClassService) DeleteClass(id uint) error {
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

// JoinPublicClass สำหรับ join public class (ไม่ต้องส่ง invite code)
func (s *ClassService) JoinPublicClass(userID, classID uint) error {
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

	// สำหรับ public class ไม่ตรวจสอบ invite code
	return s.repo.InsertClassEnrollment(userID, classID)
}

// JoinPrivateClass สำหรับ join private class โดยใช้ invite code
func (s *ClassService) JoinPrivateClass(userID uint, inviteCode string) error {
	// ค้นหาคลาสที่มี invite code ตรงกับที่ผู้ใช้ส่งมา
	class, err := s.repo.FindClassByInviteCode(inviteCode)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("invalid invite code")
	}

	joined, err := s.repo.HasJoinedClass(userID, class.ID)
	if err != nil {
		return err
	}
	if joined {
		return errors.New("already joined this class")
	}

	return s.repo.InsertClassEnrollment(userID, class.ID)
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
