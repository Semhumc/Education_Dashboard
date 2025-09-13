package models

type Class struct {
	ID        string `json:"id"`
	ClassName string `json:"class_name"`
	TeacherID string `json:"teacher_id"`
}

type ClassRepository interface {
	CreateClass(class *Class) error
	GetClassByID(id string) (*Class, error)
	UpdateClass(class *Class) error
	DeleteClass(id string) error
	GetAllClasses() ([]Class, error)
	GetClassesByTeacherID(teacherID string) ([]Class, error)
}

type ClassService interface {
	CreateClass(class *Class) error
	GetClassByID(id string) (*Class, error)
	UpdateClass(class *Class) error
	DeleteClass(id string) error
	GetAllClasses() ([]Class, error)
	GetClassesByTeacherID(teacherID string) ([]Class, error)
}
