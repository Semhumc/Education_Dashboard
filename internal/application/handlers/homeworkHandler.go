package handlers

import (
	"Education_Dashboard/internal/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HomeworkHandler struct {
	homeworkService models.HomeworkService
}

func NewHomeworkHandler(hs models.HomeworkService) *HomeworkHandler {
	return &HomeworkHandler{
		homeworkService: hs,
	}
}

func (hh *HomeworkHandler) CreateHomeworkHandler(c *fiber.Ctx) error {
	var homework models.Homework
	if err := c.BodyParser(&homework); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	if homework.TeacherID == "" || homework.LessonID == "" || homework.ClassID == "" || homework.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "teacher_id, lesson_id, class_id, and title are required",
		})
	}

	err := hh.homeworkService.CreateHomework(&homework)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Homework created successfully",
		"data":    homework,
	})
}

func (hh *HomeworkHandler) GetHomeworkByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "homework ID is required",
		})
	}

	homework, err := hh.homeworkService.GetHomeworkByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": homework,
	})
}

func (hh *HomeworkHandler) UpdateHomeworkHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "homework ID is required",
		})
	}

	var homework models.Homework
	if err := c.BodyParser(&homework); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	homework.ID = id

	err := hh.homeworkService.UpdateHomework(&homework)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Homework updated successfully",
		"data":    homework,
	})
}

func (hh *HomeworkHandler) DeleteHomeworkHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "homework ID is required",
		})
	}

	err := hh.homeworkService.DeleteHomework(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Homework deleted successfully",
	})
}

func (hh *HomeworkHandler) GetAllHomeworksHandler(c *fiber.Ctx) error {
	homeworks, err := hh.homeworkService.GetAllHomeworks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": homeworks,
	})
}

func (hh *HomeworkHandler) GetHomeworksByTeacherIDHandler(c *fiber.Ctx) error {
	teacherID := c.Params("teacherID")
	if teacherID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "teacher ID is required",
		})
	}

	homeworks, err := hh.homeworkService.GetHomeworksByTeacherID(teacherID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": homeworks,
	})
}

func (hh *HomeworkHandler) GetHomeworksByLessonIDHandler(c *fiber.Ctx) error {
	lessonID := c.Params("lessonID")
	if lessonID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "lesson ID is required",
		})
	}

	homeworks, err := hh.homeworkService.GetHomeworksByLessonID(lessonID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": homeworks,
	})
}

func (hh *HomeworkHandler) GetHomeworksByClassIDHandler(c *fiber.Ctx) error {
	classID := c.Params("classID")
	if classID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "class ID is required",
		})
	}

	homeworks, err := hh.homeworkService.GetHomeworksByClassID(classID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": homeworks,
	})
}

func (hh *HomeworkHandler) GetActiveHomeworksHandler(c *fiber.Ctx) error {
	homeworks, err := hh.homeworkService.GetActiveHomeworks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": homeworks,
	})
}

func (hh *HomeworkHandler) GetOverdueHomeworksHandler(c *fiber.Ctx) error {
	homeworks, err := hh.homeworkService.GetOverdueHomeworks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": homeworks,
	})
}

func (hh *HomeworkHandler) GetHomeworksDueSoonHandler(c *fiber.Ctx) error {
	hoursParam := c.Query("hours", "24")
	hours, err := strconv.Atoi(hoursParam)
	if err != nil || hours <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "invalid hours parameter, must be positive integer",
		})
	}

	homeworks, err := hh.homeworkService.GetHomeworksDueSoon(hours)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": homeworks,
		"hours": hours,
	})
}

func (hh *HomeworkHandler) ExtendDueDateHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "homework ID is required",
		})
	}

	var req struct {
		NewDueDate time.Time `json:"new_due_date"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	err := hh.homeworkService.ExtendDueDate(id, req.NewDueDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Due date extended successfully",
		"new_due_date": req.NewDueDate,
	})
}