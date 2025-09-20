package models

import "time"

type Schedule struct {
	ID        string    `json:"id"`
	Date      time.Time `json:"date"`
	TeacherID string    `json:"teacher_id"`
	LessonID  string    `json:"lesson_id"`
	ClassID   string    `json:"class_id"`
	Time      time.Time `json:"time"`
}

type ScheduleRepository interface {
	CreateSchedule(schedule *Schedule) error
	GetScheduleByID(id string) (*Schedule, error)
	UpdateSchedule(schedule *Schedule) error
	DeleteSchedule(id string) error
	GetAllSchedules() ([]Schedule, error)
	GetSchedulesByTeacherID(teacherID string) ([]Schedule, error)
	GetSchedulesByClassID(classID string) ([]Schedule, error)
}

type ScheduleService interface {
	CreateSchedule(schedule *Schedule) error
	GetScheduleByID(id string) (*Schedule, error)
	UpdateSchedule(schedule *Schedule) error
	DeleteSchedule(id string) error
	GetAllSchedules() ([]Schedule, error)
	GetSchedulesByTeacherID(teacherID string) ([]Schedule, error)
	GetSchedulesByClassID(classID string) ([]Schedule, error)
	RescheduleSchedule(scheduleID string, newDate time.Time, newTime time.Time) error
	GetScheduleConflicts(teacherID, classID string, date time.Time, startTime time.Time) ([]Schedule, error)
	GetUpcomingSchedules(teacherID string, days int) ([]Schedule, error)
	GetWeekSchedules(startDate time.Time) ([]Schedule, error)
	GetTodaySchedules() ([]Schedule, error)
}
