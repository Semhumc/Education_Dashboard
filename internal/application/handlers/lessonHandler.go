package handlers

import (
	"Education_Dashboard/internal/models"

	"github.com/gofiber/fiber/v2"
)

type LessonHandler struct {
	lessonService models.LessonService
}

func NewLessonHandler(ls models.LessonService) *LessonHandler {
	return &LessonHandler{
		lessonService: ls,
	}
}

func (lh *LessonHandler) CreateLessonHandler(c *fiber.Ctx) error {
	var lesson models.Lesson
	if err := c.BodyParser(&lesson); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	if lesson.LessonName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "lesson_name is required",
		})
	}

	err := lh.lessonService.CreateLesson(&lesson)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Lesson created successfully",
		"data":    lesson,
	})
}

func (lh *LessonHandler) GetLessonByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "lesson ID is required",
		})
	}

	lesson, err := lh.lessonService.GetLessonByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": lesson,
	})
}

func (lh *LessonHandler) UpdateLessonHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "lesson ID is required",
		})
	}

	var lesson models.Lesson
	if err := c.BodyParser(&lesson); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	lesson.ID = id

	if lesson.LessonName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "lesson_name is required",
		})
	}

	err := lh.lessonService.UpdateLesson(&lesson)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Lesson updated successfully",
		"data":    lesson,
	})
}

func (lh *LessonHandler) DeleteLessonHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "lesson ID is required",
		})
	}

	err := lh.lessonService.DeleteLesson(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Lesson deleted successfully",
	})
}

func (lh *LessonHandler) GetAllLessonsHandler(c *fiber.Ctx) error {
	lessons, err := lh.lessonService.GetAllLessons()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": lessons,
	})
}


