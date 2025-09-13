package models

type Lesson struct {
	ID          string `json:"id"`
	LessonName  string `json:"lesson_name"`
}


type LessonRepository interface {
	CreateLesson(lesson *Lesson) error
	GetLessonByID(id string) (*Lesson, error)
	UpdateLesson(lesson *Lesson) error
	DeleteLesson(id string) error
	GetAllLessons() ([]Lesson, error)
}

type LessonService interface {
	CreateLesson(lesson *Lesson) error
	GetLessonByID(id string) (*Lesson, error)
	UpdateLesson(lesson *Lesson) error
	DeleteLesson(id string) error
	GetAllLessons() ([]Lesson, error)
}
