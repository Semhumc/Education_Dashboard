package application

import "Education_Dashboard/internal/models"

type AttendanceService struct {
	attendanceService models.AttendanceService
}

func NewAttendanceService(as models.AttendanceService) models.AttendanceService{
	return &AttendanceService{
		attendanceService: as,
	}
}
func(as AttendanceService)CreateAttendance(attendance *models.Attendance) error{
	return nil
}

func(as AttendanceService) GetAttendanceByID(id string) (*models.Attendance, error){
	return nil,nil
}
func(as AttendanceService) UpdateAttendance(attendance *models.Attendance) error{
	return nil
}

func(as AttendanceService) DeleteAttendance(id string) error{

	return nil
}

func(as AttendanceService) GetAttendanceByStudentID(studentID string) ([]models.Attendance, error){
	return nil,nil
}

func(as AttendanceService) GetAttendanceByScheduleID(scheduleID string) ([]models.Attendance, error){
	return nil,nil
}