package services

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/repositories"

	"time"

	"github.com/golang-jwt/jwt/v4"
)

// กำหนด secret key
var jwtSecret = []byte("secret-reset-token")

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: repositories.NewUserRepository(),
	}
}

func (s *UserService) Login(username, password string) (string, error) {
	// user ตาม username
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", err
	}
	if user == nil {
		// ไม่มี user
		return "", errors.New("user not found")
	}

	// ตรวจสอบ password ด้วย bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// password ไม่ตรง
		return "", errors.New("invalid password")
	}

	// สร้าง JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}

// CreateUser ตรวจสอบซ้ำ username + แฮช password ก่อนบันทึก
func (s *UserService) CreateUser(user *models.User) error {
	// 1) ตรวจสอบความถูกต้องเบื้องต้น
	if user.Username == "" || user.Password == "" {
		return errors.New("username or password cannot be empty")
	}

	// 2) ตรวจสอบ username ซ้ำ
	existing, err := s.repo.FindByUsername(user.Username)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("username already exists")
	}

	// 3) แฮชรหัสผ่าน
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashed)

	// 4) บันทึกลง DB
	return s.repo.Create(user)
}

func (s *UserService) UpdateUser(user *models.User) error {
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}

// GenerateResetToken สร้าง reset token สำหรับรีเซ็ตรหัสผ่าน แล้วบันทึกลง DB
// แล้วส่ง (หรือแสดง) ลิงก์สำหรับรีเซ็ตรหัสผ่าน
func (s *UserService) GenerateResetPasswordToken(email string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("email not found")
	}

	// สร้าง claims สำหรับ token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(), // token หมดอายุใน 15 นาที
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ResetPassword ตรวจสอบ reset token แล้วเปลี่ยนรหัสผ่านใหม่
func (s *UserService) ResetPassword(tokenString, newPassword string) error {
	// แกะ JWT
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// ตรวจสอบ Signing Method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid or expired token")
	}

	// ดึง claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid token claims")
	}

	// ดึง user_id จาก claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return errors.New("invalid user_id in token")
	}
	userID := uint(userIDFloat)

	// แฮชรหัสผ่านใหม่
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// อัปเดตใน DB
	return s.repo.UpdatePassword(userID, string(hashed))
}
