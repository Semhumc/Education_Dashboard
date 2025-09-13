package models

type Class struct {
	ID        string `json:"id"`
	ClassName string `json:"class_name"`
	UserID    string `json:"user_id"`
}

type ClassRepository interface {
	CreateClass(class *Class) error
	GetClassByID(id string) (*Class, error)
	UpdateClass(class *Class) error
	DeleteClass(id string) error
	GetAllClasses() ([]Class, error)
	GetClassesByUserID(userID string) ([]Class, error)
}

type ClassService interface {
	CreateClass(class *Class) error
	GetClassByID(id string) (*Class, error)
	UpdateClass(class *Class) error
	DeleteClass(id string) error
	GetAllClasses() ([]Class, error)
	GetClassesByUserID(userID string) ([]Class, error)
}
