package handlers

import (
	"Education_Dashboard/internal/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ScheduleHandler struct {
	scheduleService models.ScheduleService
}

func NewScheduleHandler(ss models.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{
		scheduleService: ss,
	}
}

func (sh *ScheduleHandler) CreateScheduleHandler(c *fiber.Ctx) error {
	var schedule models.Schedule
	if err := c.BodyParser(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	if schedule.TeacherID == "" || schedule.LessonID == "" || schedule.ClassID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "teacher_id, lesson_id, and class_id are required",
		})
	}

	err := sh.scheduleService.CreateSchedule(&schedule)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Schedule created successfully",
		"data":    schedule,
	})
}

func (sh *ScheduleHandler) GetScheduleByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "schedule ID is required",
		})
	}

	schedule, err := sh.scheduleService.GetScheduleByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": schedule,
	})
}

func (sh *ScheduleHandler) UpdateScheduleHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "schedule ID is required",
		})
	}

	var schedule models.Schedule
	if err := c.BodyParser(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	schedule.ID = id

	err := sh.scheduleService.UpdateSchedule(&schedule)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Schedule updated successfully",
		"data":    schedule,
	})
}

func (sh *ScheduleHandler) DeleteScheduleHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "schedule ID is required",
		})
	}

	err := sh.scheduleService.DeleteSchedule(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Schedule deleted successfully",
	})
}

func (sh *ScheduleHandler) GetAllSchedulesHandler(c *fiber.Ctx) error {
	schedules, err := sh.scheduleService.GetAllSchedules()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": schedules,
	})
}

func (sh *ScheduleHandler) GetSchedulesByTeacherIDHandler(c *fiber.Ctx) error {
	teacherID := c.Params("teacherID")
	if teacherID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "teacher ID is required",
		})
	}

	schedules, err := sh.scheduleService.GetSchedulesByTeacherID(teacherID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": schedules,
	})
}

func (sh *ScheduleHandler) GetSchedulesByClassIDHandler(c *fiber.Ctx) error {
	classID := c.Params("classID")
	if classID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "class ID is required",
		})
	}

	schedules, err := sh.scheduleService.GetSchedulesByClassID(classID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": schedules,
	})
}

func (sh *ScheduleHandler) GetTodaySchedulesHandler(c *fiber.Ctx) error {
	schedules, err := sh.scheduleService.GetTodaySchedules()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": schedules,
		"date": time.Now().Format("2006-01-02"),
	})
}

func (sh *ScheduleHandler) GetWeekSchedulesHandler(c *fiber.Ctx) error {
	dateParam := c.Query("start_date")
	var startDate time.Time
	var err error

	if dateParam == "" {
		// Default to current week
		startDate = time.Now()
	} else {
		startDate, err = time.Parse("2006-01-02", dateParam)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": "invalid start_date format, use YYYY-MM-DD",
			})
		}
	}

	schedules, err := sh.scheduleService.GetWeekSchedules(startDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       schedules,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   startDate.AddDate(0, 0, 7).Format("2006-01-02"),
	})
}

func (sh *ScheduleHandler) GetUpcomingSchedulesHandler(c *fiber.Ctx) error {
	teacherID := c.Params("teacherID")
	if teacherID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "teacher ID is required",
		})
	}

	daysParam := c.Query("days", "7")
	days, err := strconv.Atoi(daysParam)
	if err != nil || days <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "invalid days parameter, must be positive integer",
		})
	}

	schedules, err := sh.scheduleService.GetUpcomingSchedules(teacherID, days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       schedules,
		"teacher_id": teacherID,
		"days":       days,
	})
}

func (sh *ScheduleHandler) CheckConflictsHandler(c *fiber.Ctx) error {
	var req struct {
		TeacherID string    `json:"teacher_id"`
		ClassID   string    `json:"class_id"`
		Date      time.Time `json:"date"`
		Time      time.Time `json:"time"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	conflicts, err := sh.scheduleService.GetScheduleConflicts(req.TeacherID, req.ClassID, req.Date, req.Time)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	hasConflicts := len(conflicts) > 0

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"has_conflicts": hasConflicts,
		"conflicts":     conflicts,
		"count":         len(conflicts),
	})
}

func (sh *ScheduleHandler) RescheduleHandler(c *fiber.Ctx) error {
	scheduleID := c.Params("id")
	if scheduleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "schedule ID is required",
		})
	}

	var req struct {
		NewDate time.Time `json:"new_date"`
		NewTime time.Time `json:"new_time"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	err := sh.scheduleService.RescheduleSchedule(scheduleID, req.NewDate, req.NewTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Schedule rescheduled successfully",
		"new_date": req.NewDate.Format("2006-01-02"),
		"new_time": req.NewTime.Format("15:04"),
	})
}
