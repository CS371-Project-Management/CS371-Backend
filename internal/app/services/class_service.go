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

// func NewUserService() *UserService {
// 	return &UserService{
// 		repo: repositories.NewUserRepository(),
// 	}
// }

func (s *ClassService) CreateClass(class *models.Class) error {
    // 1) ตรวจสอบความถูกต้องเบื้องต้น
    if class.Title == "" {
        return errors.New("class title cannot be empty")
    }
    // ถ้าไม่ระบุ accessibility ให้เป็น "private" เป็นค่า default
    if class.Accessibility == "" {
        class.Accessibility = "private"
    }

    // 2) สร้าง invite code โดยอัตโนมัติ
    class.InviteCode = generateInviteCode(6)

    // 3) เรียก repository เพื่อบันทึกใน DB
    return s.repo.CreateClass(class)
}

// GetAllClasses ติดต่อ repository เพื่อดึงรายการคลาสทั้งหมด
func (s *ClassService) GetAllClasses() ([]models.Class, error) {
	return s.repo.GetAllClasses()
}

// ฟังก์ชัน generate code
// ขโมยมา
func generateInviteCode(length int) string {
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, length)
    for i := 0; i < length; i++ {
        b[i] = chars[rand.Intn(len(chars))]
    }
    return string(b)
}


func (s *ClassService) GetClassByID(id uint) (*models.Class, error) {
    return s.repo.FindClassByID(id)
}

func (s *ClassService) UpdateClass(class *models.Class) error {
    if class.Title == "" {
        return errors.New("class title cannot be empty")
    }
    return s.repo.UpdateClass(class)
}

func (s *ClassService) DeleteClass(id uint) error {
    return s.repo.DeleteClass(id)
}

// GetInviteCode คืน invite_code ของคลาสตาม id
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

func GenerateInviteCode(length int) string {
    const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    rand.Seed(time.Now().UnixNano())
    b := make([]byte, length)
    for i := 0; i < length; i++ {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}