package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"cs371-backend/internal/app/models"
	"cs371-backend/internal/app/services"
)

type ClassController struct {
	service *services.ClassService
}

func NewClassController() *ClassController {
	return &ClassController{
		service: services.NewClassService(),
	}
}

func (cc *ClassController) CreateClassHandler(c *fiber.Ctx) error {
    // รับ JSON Body
    var req struct {
        UserID        string `json:"user_id"`
        Title         string `json:"title"`
        Description   string `json:"description"`
        Accessibility string `json:"accessibility"`
    }
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }

    // แปลง string เป็น bool สำหรับ accessibility
    accessibility, err := strconv.ParseBool(req.Accessibility)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid value for accessibility",
        })
    }

    // สร้าง model Class
    class := &models.Class{
        UserID:        req.UserID,
        Title:         req.Title,
        Description:   req.Description,
        Accessibility: accessibility,
    }

    // เรียก Service เพื่อสร้างคลาส
    if err := cc.service.CreateClass(class); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    // ตอบกลับข้อมูลที่สำคัญ
    return c.JSON(fiber.Map{
        "message":     "Class created successfully",
        "class_id":    class.ID,
        "invite_code": class.InviteCode,
    })
}

// GetAllClassesHandler ดึงรายการคลาสทั้งหมด
func (cc *ClassController) GetAllClassesHandler(c *fiber.Ctx) error {
	classes, err := cc.service.GetAllClasses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(classes)
}

// GetClassHandler ดึงข้อมูลคลาสจาก id (ใช้เป็น string)
func (cc *ClassController) GetClassHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid class ID",
		})
	}

	class, err := cc.service.GetClassByID(idParam)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if class == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Class not found",
		})
	}
	return c.JSON(class)
}

// UpdateClassHandler อัปเดตข้อมูลคลาส
func (cc *ClassController) UpdateClassHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid class ID",
		})
	}

	var req struct {
		Title         string `json:"title"`
		Description   string `json:"description"`
		Accessibility string `json:"accessibility"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	accessibility, err := strconv.ParseBool(req.Accessibility)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid value for accessibility",
		})
	}

	// สร้าง invite_code ใหม่โดยระบบ
	newInviteCode := services.GenerateInviteCode(6)

	// สร้าง model Class ที่จะอัปเดต โดยใช้ id เป็น string
	class := &models.Class{
		ID:            idParam,
		InviteCode:    newInviteCode,
		Title:         req.Title,
		Description:   req.Description,
		Accessibility: accessibility,
	}

	if err := cc.service.UpdateClass(class); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Class updated successfully",
		"class_id":    class.ID,
		"invite_code": class.InviteCode,
	})
}

func (cc *ClassController) DeleteClassHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid class ID",
		})
	}

	err = cc.service.DeleteClass(uint(idUint64))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Class deleted successfully",
	})
}

func (cc *ClassController) GetInviteCodeHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	idUint64, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid class ID",
		})
	}
	inviteCode, err := cc.service.GetInviteCode(uint(idUint64))
	if err != nil {
		switch err.Error() {
		case "class not found":
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Class not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}
	return c.JSON(fiber.Map{
		"invite_code": inviteCode,
	})
}

// JoinPublicClassHandler สำหรับผู้ใช้เข้าร่วมคลาสสาธารณะ
func (cc *ClassController) JoinPublicClassHandler(c *fiber.Ctx) error {
	// ดึง userID จาก session หรือ token (สมมติว่า middleware เก็บเป็น string)
	userIDAny := c.Locals("user_id")
	if userIDAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}
	userID, ok := userIDAny.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user id"})
	}

	classID := c.Params("id")
	if classID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid class ID"})
	}

	if err := cc.service.JoinPublicClass(userID, classID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Joined public class successfully"})
}

// JoinPrivateClassHandler สำหรับผู้ใช้เข้าร่วมคลาสแบบส่วนตัวผ่าน invite code
func (cc *ClassController) JoinPrivateClassHandler(c *fiber.Ctx) error {
	userIDAny := c.Locals("user_id")
	if userIDAny == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}
	userID, ok := userIDAny.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user id"})
	}

	var req struct {
		InviteCode string `json:"invite_code"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.InviteCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invite code is required"})
	}

	if err := cc.service.JoinPrivateClass(userID, req.InviteCode); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Joined private class successfully"})
}

// LeaveClassHandler - Handler สำหรับออกจากคลาส
func (cc *ClassController) LeaveClassHandler(c *fiber.Ctx) error {
    // ดึง userID จาก middleware (เช่น c.Locals("user_id"))
    userIDAny := c.Locals("user_id")
    if userIDAny == nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
    }
    userID, ok := userIDAny.(string)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user id"})
    }

    // ดึง classID จาก URL parameter เช่น /api/classes/leaveClass/:id
    classID := c.Params("id")
    if classID == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid class ID"})
    }

    // เรียกใช้งาน Service เพื่อทำกระบวนการ leave class
    if err := cc.service.LeaveClass(userID, classID); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    // สำเร็จ
    return c.JSON(fiber.Map{"message": "Left class successfully"})
}

// RemoveUserFromClassHandler - API สำหรับลบผู้ใช้จากคลาส
func (cc *ClassController) RemoveUserFromClassHandler(c *fiber.Ctx) error {
	classID := c.Params("class_id") // รับ class ID จาก URL
	userID := c.Params("user_id")   // รับ user ID จาก URL

	if classID == "" || userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID or class ID"})
	}

	// เรียกใช้ Service เพื่อลบผู้ใช้
	if err := cc.service.RemoveUserFromClass(userID, classID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User removed from class successfully"})
}