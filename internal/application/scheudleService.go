package application

import (
	"Education_Dashboard/internal/models"
	"fmt"
	"time"
)

type ScheduleService struct {
	scheduleRepo    models.ScheduleRepository
	lessonRepo      models.LessonRepository
	attendanceRepo  models.AttendanceRepository
}

func NewScheduleService(scheduleRepo models.ScheduleRepository, lessonRepo models.LessonRepository, attendanceRepo models.AttendanceRepository) models.ScheduleService {
	return &ScheduleService{
		scheduleRepo:   scheduleRepo,
		lessonRepo:     lessonRepo,
		attendanceRepo: attendanceRepo,
	}
}

func (ss *ScheduleService) CreateSchedule(schedule *models.Schedule) error {
	// Validate required fields
	if schedule.TeacherID == "" {
		return fmt.Errorf("teacher ID is required")
	}

	if schedule.LessonID == "" {
		return fmt.Errorf("lesson ID is required")
	}

	if schedule.ClassID == "" {
		return fmt.Errorf("class ID is required")
	}

	// Validate lesson exists
	_, err := ss.lessonRepo.GetLessonByID(schedule.LessonID)
	if err != nil {
		return fmt.Errorf("lesson not found: %w", err)
	}

	// Validate date is not in the past
	today := time.Now().Truncate(24 * time.Hour)
	scheduleDate := schedule.Date.Truncate(24 * time.Hour)
	
	if scheduleDate.Before(today) {
		return fmt.Errorf("cannot create schedule for past dates")
	}

	// Check for conflicts - same teacher, same time, same day
	existingSchedules, err := ss.scheduleRepo.GetSchedulesByTeacherID(schedule.TeacherID)
	if err != nil {
		return fmt.Errorf("failed to check teacher schedules: %w", err)
	}

	for _, existing := range existingSchedules {
		if isSameDay(existing.Date, schedule.Date) && isSameTime(existing.Time, schedule.Time) {
			return fmt.Errorf("teacher already has a schedule at this time on this date")
		}
	}

	// Check for class conflicts - same class, same time, same day
	classSchedules, err := ss.scheduleRepo.GetSchedulesByClassID(schedule.ClassID)
	if err != nil {
		return fmt.Errorf("failed to check class schedules: %w", err)
	}

	for _, existing := range classSchedules {
		if isSameDay(existing.Date, schedule.Date) && isSameTime(existing.Time, schedule.Time) {
			return fmt.Errorf("class already has a schedule at this time on this date")
		}
	}

	return ss.scheduleRepo.CreateSchedule(schedule)
}

func (ss *ScheduleService) GetScheduleByID(id string) (*models.Schedule, error) {
	if id == "" {
		return nil, fmt.Errorf("schedule ID is required")
	}

	return ss.scheduleRepo.GetScheduleByID(id)
}

func (ss *ScheduleService) UpdateSchedule(schedule *models.Schedule) error {
	// Validate schedule exists
	existing, err := ss.scheduleRepo.GetScheduleByID(schedule.ID)
	if err != nil {
		return fmt.Errorf("schedule not found: %w", err)
	}

	// Validate required fields
	if schedule.TeacherID == "" {
		return fmt.Errorf("teacher ID is required")
	}

	if schedule.LessonID == "" {
		return fmt.Errorf("lesson ID is required")
	}

	if schedule.ClassID == "" {
		return fmt.Errorf("class ID is required")
	}

	// Validate lesson exists if changed
	if existing.LessonID != schedule.LessonID {
		_, err := ss.lessonRepo.GetLessonByID(schedule.LessonID)
		if err != nil {
			return fmt.Errorf("lesson not found: %w", err)
		}
	}

	// Business rule: Cannot change date/time if schedule is in the past
	if existing.Date.Before(time.Now().Truncate(24 * time.Hour)) {
		if !isSameDay(existing.Date, schedule.Date) || !isSameTime(existing.Time, schedule.Time) {
			return fmt.Errorf("cannot modify date/time for past schedules")
		}
	}

	// Validate new date is not in the past
	today := time.Now().Truncate(24 * time.Hour)
	scheduleDate := schedule.Date.Truncate(24 * time.Hour)
	
	if scheduleDate.Before(today) {
		return fmt.Errorf("cannot schedule for past dates")
	}

	// Check for conflicts only if date/time/teacher/class changed
	if !isSameDay(existing.Date, schedule.Date) || 
	   !isSameTime(existing.Time, schedule.Time) || 
	   existing.TeacherID != schedule.TeacherID ||
	   existing.ClassID != schedule.ClassID {
		
		// Check teacher conflicts
		if existing.TeacherID != schedule.TeacherID || 
		   !isSameDay(existing.Date, schedule.Date) || 
		   !isSameTime(existing.Time, schedule.Time) {
			
			teacherSchedules, err := ss.scheduleRepo.GetSchedulesByTeacherID(schedule.TeacherID)
			if err != nil {
				return fmt.Errorf("failed to check teacher schedules: %w", err)
			}

			for _, other := range teacherSchedules {
				if other.ID != schedule.ID && isSameDay(other.Date, schedule.Date) && isSameTime(other.Time, schedule.Time) {
					return fmt.Errorf("teacher already has a schedule at this time on this date")
				}
			}
		}

		// Check class conflicts
		if existing.ClassID != schedule.ClassID || 
		   !isSameDay(existing.Date, schedule.Date) || 
		   !isSameTime(existing.Time, schedule.Time) {
			
			classSchedules, err := ss.scheduleRepo.GetSchedulesByClassID(schedule.ClassID)
			if err != nil {
				return fmt.Errorf("failed to check class schedules: %w", err)
			}

			for _, other := range classSchedules {
				if other.ID != schedule.ID && isSameDay(other.Date, schedule.Date) && isSameTime(other.Time, schedule.Time) {
					return fmt.Errorf("class already has a schedule at this time on this date")
				}
			}
		}
	}

	return ss.scheduleRepo.UpdateSchedule(schedule)
}

func (ss *ScheduleService) DeleteSchedule(id string) error {
	if id == "" {
		return fmt.Errorf("schedule ID is required")
	}

	// Validate schedule exists
	schedule, err := ss.scheduleRepo.GetScheduleByID(id)
	if err != nil {
		return fmt.Errorf("schedule not found: %w", err)
	}

	// Business rule: Cannot delete schedules from the past
	if schedule.Date.Before(time.Now().Truncate(24 * time.Hour)) {
		return fmt.Errorf("cannot delete past schedules")
	}

	// Check for associated attendances
	attendances, err := ss.attendanceRepo.GetAttendanceByScheduleID(id)
	if err != nil {
		return fmt.Errorf("failed to check associated attendances: %w", err)
	}

	if len(attendances) > 0 {
		return fmt.Errorf("cannot delete schedule with existing attendance records (%d found)", len(attendances))
	}

	return ss.scheduleRepo.DeleteSchedule(id)
}

func (ss *ScheduleService) GetAllSchedules() ([]models.Schedule, error) {
	return ss.scheduleRepo.GetAllSchedules()
}

func (ss *ScheduleService) GetSchedulesByTeacherID(teacherID string) ([]models.Schedule, error) {
	if teacherID == "" {
		return nil, fmt.Errorf("teacher ID is required")
	}

	return ss.scheduleRepo.GetSchedulesByTeacherID(teacherID)
}

func (ss *ScheduleService) GetSchedulesByClassID(classID string) ([]models.Schedule, error) {
	if classID == "" {
		return nil, fmt.Errorf("class ID is required")
	}

	return ss.scheduleRepo.GetSchedulesByClassID(classID)
}

// Additional business methods

func (ss *ScheduleService) GetTodaySchedules() ([]models.Schedule, error) {
	allSchedules, err := ss.scheduleRepo.GetAllSchedules()
	if err != nil {
		return nil, fmt.Errorf("failed to get all schedules: %w", err)
	}

	today := time.Now().Truncate(24 * time.Hour)
	var todaySchedules []models.Schedule

	for _, schedule := range allSchedules {
		if isSameDay(schedule.Date, today) {
			todaySchedules = append(todaySchedules, schedule)
		}
	}

	return todaySchedules, nil
}

func (ss *ScheduleService) GetWeekSchedules(startDate time.Time) ([]models.Schedule, error) {
	allSchedules, err := ss.scheduleRepo.GetAllSchedules()
	if err != nil {
		return nil, fmt.Errorf("failed to get all schedules: %w", err)
	}

	startOfWeek := startDate.Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	var weekSchedules []models.Schedule
	for _, schedule := range allSchedules {
		scheduleDate := schedule.Date.Truncate(24 * time.Hour)
		if (scheduleDate.Equal(startOfWeek) || scheduleDate.After(startOfWeek)) && scheduleDate.Before(endOfWeek) {
			weekSchedules = append(weekSchedules, schedule)
		}
	}

	return weekSchedules, nil
}

func (ss *ScheduleService) GetUpcomingSchedules(teacherID string, days int) ([]models.Schedule, error) {
	if teacherID == "" {
		return nil, fmt.Errorf("teacher ID is required")
	}

	if days <= 0 {
		return nil, fmt.Errorf("days must be positive")
	}

	teacherSchedules, err := ss.scheduleRepo.GetSchedulesByTeacherID(teacherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher schedules: %w", err)
	}

	now := time.Now()
	endDate := now.AddDate(0, 0, days)

	var upcomingSchedules []models.Schedule
	for _, schedule := range teacherSchedules {
		if schedule.Date.After(now) && schedule.Date.Before(endDate) {
			upcomingSchedules = append(upcomingSchedules, schedule)
		}
	}

	return upcomingSchedules, nil
}

func (ss *ScheduleService) GetScheduleConflicts(teacherID, classID string, date time.Time, startTime time.Time) ([]models.Schedule, error) {
	var conflicts []models.Schedule

	// Check teacher conflicts
	if teacherID != "" {
		teacherSchedules, err := ss.scheduleRepo.GetSchedulesByTeacherID(teacherID)
		if err != nil {
			return nil, fmt.Errorf("failed to get teacher schedules: %w", err)
		}

		for _, schedule := range teacherSchedules {
			if isSameDay(schedule.Date, date) && isSameTime(schedule.Time, startTime) {
				conflicts = append(conflicts, schedule)
			}
		}
	}

	// Check class conflicts
	if classID != "" {
		classSchedules, err := ss.scheduleRepo.GetSchedulesByClassID(classID)
		if err != nil {
			return nil, fmt.Errorf("failed to get class schedules: %w", err)
		}

		for _, schedule := range classSchedules {
			if isSameDay(schedule.Date, date) && isSameTime(schedule.Time, startTime) {
				// Avoid duplicates
				found := false
				for _, conflict := range conflicts {
					if conflict.ID == schedule.ID {
						found = true
						break
					}
				}
				if !found {
					conflicts = append(conflicts, schedule)
				}
			}
		}
	}

	return conflicts, nil
}

func (ss *ScheduleService) RescheduleSchedule(scheduleID string, newDate time.Time, newTime time.Time) error {
	if scheduleID == "" {
		return fmt.Errorf("schedule ID is required")
	}

	schedule, err := ss.scheduleRepo.GetScheduleByID(scheduleID)
	if err != nil {
		return fmt.Errorf("schedule not found: %w", err)
	}

	// Check if the schedule is in the past
	if schedule.Date.Before(time.Now().Truncate(24 * time.Hour)) {
		return fmt.Errorf("cannot reschedule past schedules")
	}

	// Check for conflicts
	conflicts, err := ss.GetScheduleConflicts(schedule.TeacherID, schedule.ClassID, newDate, newTime)
	if err != nil {
		return fmt.Errorf("failed to check conflicts: %w", err)
	}

	for _, conflict := range conflicts {
		if conflict.ID != scheduleID {
			return fmt.Errorf("schedule conflict found at new time slot")
		}
	}

	// Update the schedule
	schedule.Date = newDate
	schedule.Time = newTime

	return ss.scheduleRepo.UpdateSchedule(schedule)
}

// Helper functions
func isSameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func isSameTime(time1, time2 time.Time) bool {
	h1, m1, _ := time1.Clock()
	h2, m2, _ := time2.Clock()
	return h1 == h2 && m1 == m2
}