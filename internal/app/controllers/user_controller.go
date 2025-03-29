package controllers

import (
	"github.com/gofiber/fiber/v2"

	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/services"
)

type UserController struct {
	Service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{
		Service: service,
	}
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	//  JWT ไม่สามารถลบ token ที่ถูกเซ็นแล้วจากฝั่ง server ได้
	// response แจ้งให้ client ทราบว่า
	return ctx.JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

// controllers/UserController.go
func (uc *UserController) LoginHandler(c *fiber.Ctx) error {
	// รับ JSON Body
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// เรียก Service Login ที่คืนค่า userID, token, err
	userID, token, err := uc.Service.Login(req.Username, req.Password)
	if err != nil {
		switch err.Error() {
		case "user not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		case "invalid password":
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// ส่งกลับ response พร้อม user_id
	return c.JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"user_id": userID,
	})
}

// GetAllUsers ดึงผู้ใช้ทั้งหมด
func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	users, err := c.Service.GetAllUsers()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(users)
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	user, err := c.Service.GetUserByID(id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return ctx.JSON(user)
}

// CreateUser สร้างผู้ใช้ใหม่
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// เรียก Service เพื่อสร้างผู้ใช้
	if err := c.Service.CreateUser(user); err != nil {
		// ตัวอย่าง: ถ้า username ซ้ำหรือ error อื่น จะถูก return ที่นี่
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

// UpdateUser อัปเดตข้อมูลผู้ใช้
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	user := new(models.User)
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user.ID = id
	if err := c.Service.UpdateUser(user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(user)
}

// DeleteUser ลบผู้ใช้
func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := c.Service.DeleteUser(id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 204 No Content
	return ctx.SendStatus(fiber.StatusNoContent)
}

// RequestResetPassword รับ email จากผู้ใช้ เพื่อขอรีเซ็ตรหัสผ่าน
func (uc *UserController) RequestResetPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	token, err := uc.Service.GenerateResetPasswordToken(req.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// ส่งอีเมลจริง หรือ mock แสดงใน console
	resetLink := "https://example.com/reset-password?token=" + token
	println("Reset link:", resetLink)

	return c.JSON(fiber.Map{
		"message": "Reset password link sent",
	})
}

// ResetPassword รับ reset token และ newPassword แล้วเปลี่ยนรหัสผ่าน
func (uc *UserController) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := uc.Service.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password reset successful",
	})
}
