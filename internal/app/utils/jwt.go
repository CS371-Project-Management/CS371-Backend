package utils

import (
	"cs371-backend/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateToken สร้าง JWT token สำหรับผู้ใช้
func GenerateToken(userID uint) (string, error) {
	// สร้าง token ใหม่
	token := jwt.New(jwt.SigningMethodHS256)

	// ตั้งค่า claims
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // หมดอายุใน 24 ชั่วโมง

	// ลงชื่อ token ด้วย secret key
	secretKey := config.GetRequiredEnv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken ตรวจสอบความถูกต้องของ token
func ValidateToken(tokenString string) (uint, error) {
	// แปลง token string เป็น token object
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบว่าเป็นวิธีการลงชื่อที่ถูกต้อง
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// ส่งคืน secret key สำหรับการตรวจสอบลายเซ็น
		secretKey := config.GetRequiredEnv("JWT_SECRET")
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, err
	}

	// ตรวจสอบความถูกต้องของ token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// แปลง user_id เป็น uint
		userID := uint(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, fmt.Errorf("invalid token")
}
