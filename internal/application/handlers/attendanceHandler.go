package handlers

import (
	"Education_Dashboard/internal/models"

	"github.com/gofiber/fiber/v2"
)

type AttendanceHandler struct {
	attendanceService models.AttendanceService
}

func NewAttendanceHandler(as models.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceService: as,
	}
}

func (ah AttendanceHandler) CreateAttendanceHandler(c *fiber.Ctx) error {
	var attendance models.Attendance
	if err := c.BodyParser(&attendance); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	if attendance.StudentID == "" || attendance.ScheduleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "student_id and schedule_id are required",
		})
	}

	err := ah.attendanceService.CreateAttendance(&attendance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Attendance created successfully",
		"data":    attendance,
	})
}

func (ah *AttendanceHandler) UpdateAttendanceHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "attendance ID is required",
		})
	}

	var attendance models.Attendance
	if err := c.BodyParser(&attendance); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	attendance.ID = id

	err := ah.attendanceService.UpdateAttendance(&attendance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Attendance updated successfully",
		"data":    attendance,
	})
}

func (ah *AttendanceHandler) GetAttendanceByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "attendance ID is required",
		})
	}

	attendance, err := ah.attendanceService.GetAttendanceByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": attendance,
	})
}

func (ah AttendanceHandler) DeleteAttendanceHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "attendance ID is required",
		})
	}

	err := ah.attendanceService.DeleteAttendance(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Attendance deleted successfully",
	})
}

func (ah *AttendanceHandler) GetAttendanceByStudentIDHandler(c *fiber.Ctx) error {
	studentID := c.Params("studentID")
	if studentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "student ID is required",
		})
	}

	attendances, err := ah.attendanceService.GetAttendanceByStudentID(studentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": attendances,
	})
}

func (ah *AttendanceHandler) GetAttendanceByScheduleIDHandler(c *fiber.Ctx) error {
	scheduleID := c.Params("scheduleID")
	if scheduleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "schedule ID is required",
		})
	}

	attendances, err := ah.attendanceService.GetAttendanceByScheduleID(scheduleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": attendances,
	})
}

func (ah *AttendanceHandler) MarkAttendanceHandler(c *fiber.Ctx) error {
	var req struct {
		StudentID  string `json:"student_id"`
		ScheduleID string `json:"schedule_id"`
		IsPresent  bool   `json:"is_present"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	if req.StudentID == "" || req.ScheduleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "student_id and schedule_id are required",
		})
	}

	err := ah.attendanceService.MarkAttendance(req.StudentID, req.ScheduleID, req.IsPresent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Attendance marked successfully",
	})
}
