package models

type Attendance struct {
	ID         string `json:"id"`
	StudentID  string `json:"student_id"`
	ScheduleID string `json:"schedule_id"`
	Here       bool   `json:"here"`
	Counter    int    `json:"counter"`
}

type AttendanceRepository interface {
	CreateAttendance(attendance *Attendance) error
	GetAttendanceByID(id string) (*Attendance, error)
	UpdateAttendance(attendance *Attendance) error
	DeleteAttendance(id string) error
	GetAttendanceByStudentID(studentID string) ([]Attendance, error)
	GetAttendanceByScheduleID(scheduleID string) ([]Attendance, error)
}

type AttendanceService interface {
	CreateAttendance(attendance *Attendance) error
	GetAttendanceByID(id string) (*Attendance, error)
	UpdateAttendance(attendance *Attendance) error
	DeleteAttendance(id string) error
	GetAttendanceByStudentID(studentID string) ([]Attendance, error)
	GetAttendanceByScheduleID(scheduleID string) ([]Attendance, error)
}
