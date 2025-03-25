package controllers

import (
    "github.com/gofiber/fiber/v2"
    "cs371-backend/internal/app/models"
    "cs371-backend/internal/app/services"
    "strconv"
)

type ClassController struct {
    service *services.ClassService
}

func NewClassController() *ClassController {
    return &ClassController{
        service: services.NewClassService(), // เรียก service constructor
    }
}

func (cc *ClassController) CreateClassHandler(c *fiber.Ctx) error {
    // รับ JSON Body จาก Frontend
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

    // สร้างโมเดล Class (ไม่ต้องใส่ invite_code เพราะ Service จะ gen)
    class := &models.Class{
        Title:         req.Title,
        Description:   req.Description,
        Accessibility: req.Accessibility,
    }

    // เรียก Service เพื่อบันทึกคลาส
    if err := cc.service.CreateClass(class); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    // ตอบกลับ Frontend ด้วยข้อมูลสำคัญ เช่น class_id, invite_code
    return c.JSON(fiber.Map{
        "message":     "Class created successfully",
        "class_id":    class.ID,
        "invite_code": class.InviteCode,
    })
}

// GetAllClassesHandler - GET /classes
func (cc *ClassController) GetAllClassesHandler(c *fiber.Ctx) error {
	classes, err := cc.service.GetAllClasses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(classes)
}

// GetClassHandler - GET /classes/:id
func (cc *ClassController) GetClassHandler(c *fiber.Ctx) error {
    idParam := c.Params("id")
    idUint, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid class ID",
        })
    }

    class, err := cc.service.GetClassByID(uint(idUint))
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

// UpdateClassHandler - PUT /classes/:id
func (cc *ClassController) UpdateClassHandler(c *fiber.Ctx) error {
    // แปลง parameter id
    idParam := c.Params("id")
    idUint, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid class ID",
        })
    }

    // รับ JSON body (ไม่รวม invite_code)
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

    // สร้าง invite_code ใหม่โดยระบบ (ใช้ฟังก์ชัน GenerateInviteCode ใน service)
    newInviteCode := services.GenerateInviteCode(6)

    // สร้างโมเดล Class ที่จะอัปเดต
    class := &models.Class{
        ID:            uint(idUint),
        InviteCode:    newInviteCode, // ใช้ invite_code ที่ gen ใหม่
        Title:         req.Title,
        Description:   req.Description,
        Accessibility: req.Accessibility,
    }

    err = cc.service.UpdateClass(class)
    if err != nil {
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

// DeleteClassHandler - DELETE /classes/:id
func (cc *ClassController) DeleteClassHandler(c *fiber.Ctx) error {
    idParam := c.Params("id")
    idUint, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid class ID",
        })
    }

    err = cc.service.DeleteClass(uint(idUint))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "message": "Class deleted successfully",
    })
}

// GetInviteCodeHandler - GET /classes/:id/invite_code
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

func (cc *ClassController) JoinPublicClassHandler(c *fiber.Ctx) error {
    // ดึง userID จาก session หรือ token (ตัวอย่างนี้สมมติว่า middleware ได้เก็บไว้ใน c.Locals("user_id"))
    userIDAny := c.Locals("user_id")
    if userIDAny == nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
    }
    userID, ok := userIDAny.(uint)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user id"})
    }
    
    idParam := c.Params("id")
    classID64, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid class ID"})
    }
    classID := uint(classID64)

    err = cc.service.JoinPublicClass(userID, classID)
    if err != nil {
        // handle error...
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Joined public class successfully"})
}

func (cc *ClassController) JoinPrivateClassHandler(c *fiber.Ctx) error {
    // ดึง userID จาก session หรือ token
    userIDAny := c.Locals("user_id")
    if userIDAny == nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
    }
    userID, ok := userIDAny.(uint)
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

    err := cc.service.JoinPrivateClass(userID, req.InviteCode)
    if err != nil {
        // handle errors accordingly...
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"message": "Joined private class successfully"})
}