package controllers

import (
    "strconv"

    "github.com/gofiber/fiber/v2"

    "cs371-backend/internal/app/models"
    "cs371-backend/internal/app/services"
)

type UserController struct {
    service *services.UserService
}

func NewUserController() *UserController {
    return &UserController{
        service: services.NewUserService(),
    }
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	//  JWT ไม่สามารถลบ token ที่ถูกเซ็นแล้วจากฝั่ง server ได้
	// response แจ้งให้ client ทราบว่า
	return ctx.JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

// LoginHandler รับ username/password แล้วส่งผลลัพธ์ตาม Sequence
func (uc *UserController) LoginHandler(c *fiber.Ctx) error {
    // 1) รับ JSON Body
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := c.BodyParser(&req); err != nil {
        // 400 Bad Request
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    token, err := uc.service.Login(req.Username, req.Password)
    if err != nil {
        // แยกเคส error
        switch err.Error() {
        case "user not found":
            // 404 Not Found
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "error": "User not found",
            })
        case "invalid password":
            // 401 Unauthorized
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid username or password",
            })
        default:
            //อื่น ๆ
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
    }

    // 3) Login สำเร็จ -> 200 OK
    // ส่ง token กลับ
    return c.JSON(fiber.Map{
        "message": "Login successful",
        "token":   token,
    })
}

// GetAllUsers ดึงผู้ใช้ทั้งหมด
func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
    users, err := c.service.GetAllUsers()
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return ctx.JSON(users)
}

// GetUser ดึงข้อมูลผู้ใช้ด้วย ID
func (c *UserController) GetUser(ctx *fiber.Ctx) error {
    id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid ID",
        })
    }

    user, err := c.service.GetUserByID(uint(id))
    if err != nil {
        // ใน Service หรือ Repo อาจ return error ถ้าหาไม่เจอ
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
    if err := c.service.CreateUser(user); err != nil {
        // ตัวอย่าง: ถ้า username ซ้ำหรือ error อื่น จะถูก return ที่นี่
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return ctx.Status(fiber.StatusCreated).JSON(user)
}

// UpdateUser อัปเดตข้อมูลผู้ใช้
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
    id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid ID",
        })
    }

    user := new(models.User)
    if err := ctx.BodyParser(user); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    user.ID = uint(id)
    if err := c.service.UpdateUser(user); err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return ctx.JSON(user)
}

// DeleteUser ลบผู้ใช้
func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
    id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
    if err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid ID",
        })
    }

    if err := c.service.DeleteUser(uint(id)); err != nil {
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

    token, err := uc.service.GenerateResetPasswordToken(req.Email)
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

    err := uc.service.ResetPassword(req.Token, req.NewPassword)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "message": "Password reset successful",
    })
}