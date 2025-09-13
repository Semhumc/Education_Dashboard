package models

type Lesson struct {
	ID          string `json:"id"`
	LessonName  string `json:"lesson_name"`
	UserID      string `json:"user_id"`
}
