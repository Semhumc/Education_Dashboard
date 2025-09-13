package models

type Schedule struct {
	ID       string `json:"id"`
	Date     string `json:"date"`
	UserID   string `json:"user_id"`
	LessonID string `json:"lesson_id"`
	ClassID  string `json:"class_id"`
	Time     string `json:"time"`
}