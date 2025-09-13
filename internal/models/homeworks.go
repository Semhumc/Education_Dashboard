package models

import "time"

type Homework struct {
	ID        string    `json:"id"`
	TeacherID string    `json:"teacher_id"`
	LessonID  string    `json:"lesson_id"`
	ClassID   string    `json:"class_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	DueDate   time.Time `json:"due_date"`
}

type HomeworkRepository interface {
	CreateHomework(homework *Homework) error
	GetHomeworkByID(id string) (*Homework, error)
	UpdateHomework(homework *Homework) error
	DeleteHomework(id string) error
	GetAllHomeworks() ([]Homework, error)
	GetHomeworksByTeacherID(teacherID string) ([]Homework, error)
	GetHomeworksByLessonID(lessonID string) ([]Homework, error)
	GetHomeworksByClassID(classID string) ([]Homework, error)
}

type HomeworkService interface {
	CreateHomework(homework *Homework) error
	GetHomeworkByID(id string) (*Homework, error)
	UpdateHomework(homework *Homework) error
	DeleteHomework(id string) error
	GetAllHomeworks() ([]Homework, error)
}
