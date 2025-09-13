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
}
