package models


type Class struct {
	ID        string `json:"id"`
	ClassName string `json:"class_name"`
	TeacherID string `json:"teacher_id"`
}

type ClassRepository interface {
	CreateClass(class *Class) error
	UpdateClass(class *Class) error
	DeleteClass(classID string) error
	GetAllClasses() ([]Class, error)
	GetClassesByTeacherID(teacherID string) ([]Class, error)
}

type ClassService interface {
	CreateClass(class *Class) (string,error)
	UpdateClass(class *Class) error
	DeleteClass(classID string) error
	GetAllClasses() ([]Class, error)
	GetClassesByTeacherID(teacherID string) ([]Class, error)
	GetStudentsByClassID(classID string) ([]User, error) // New method
}
