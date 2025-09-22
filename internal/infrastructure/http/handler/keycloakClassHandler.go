package handler

import (
	"Education_Dashboard/internal/models"

	"github.com/gofiber/fiber/v2"
)

type KeycloakClassHandler struct {
	keycloakClassService models.ClassService
}

func NewKeycloakClassHandler(kcc models.ClassService) KeycloakClassHandler {
	return KeycloakClassHandler{
		keycloakClassService: kcc,
	}
}

func (kcc KeycloakClassHandler) CreateClassHandler(c *fiber.Ctx) error{
	var classParams models.Class
	if err := c.BodyParser(&classParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{	
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	if classParams.ClassName == "" || classParams.TeacherID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "name, teacherid are required",
		})
	}

	_,err := kcc.keycloakClassService.CreateClass(&classParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{	
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Create the class successfully",
	})
}

func (kcc KeycloakClassHandler) DeleteClassHandler(c *fiber.Ctx) error {
	classID := c.Params("id")
	if classID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "Class ID is required",
		})
	}

	err := kcc.keycloakClassService.DeleteClass(classID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Class deleted successfully",
	})
}

func (kcc KeycloakClassHandler) UpdateClassHandler(c *fiber.Ctx) error {
	classID := c.Params("id")
	if classID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "Class ID is required",
		})
	}

	var classParams models.Class
	if err := c.BodyParser(&classParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	// Güncelleme isteğinde ID'nin body'de olması beklenmiyorsa burayı düzenleyebilirsiniz.
	// Genellikle ID URL parametresinden alınır ve body'deki diğer alanlar güncellenir.
	classParams.ID = classID // URL'den gelen ID'yi atama

	err := kcc.keycloakClassService.UpdateClass(&classParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Class updated successfully",
	})
}

func (kcc KeycloakClassHandler) GetAllClassesHandler(c *fiber.Ctx) error {
	classes, err := kcc.keycloakClassService.GetAllClasses()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.JSON(classes)
}

func (kcc KeycloakClassHandler) GetClassesByTeacherIDHandler(c *fiber.Ctx) error {
	teacherID := c.Params("teacherID") // URL parametresi olarak teacherID'yi alalım
	if teacherID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "Teacher ID is required",
		})
	}

	classes, err := kcc.keycloakClassService.GetClassesByTeacherID(teacherID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.JSON(classes)
}

func (kcc KeycloakClassHandler) GetStudentsByClassIDHandler(c *fiber.Ctx) error {
	classID := c.Params("id")
	if classID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "Class ID is required",
		})
	}

	students, err := kcc.keycloakClassService.GetStudentsByClassID(classID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.JSON(students)
}

