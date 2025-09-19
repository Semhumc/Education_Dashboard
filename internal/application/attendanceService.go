package application

import (
	"Education_Dashboard/internal/models"
	"fmt"
)

type AttendanceService struct {
	attendanceRepo models.AttendanceRepository
	scheduleRepo   models.ScheduleRepository
}

func NewAttendanceService(attendanceRepo models.AttendanceRepository, scheduleRepo models.ScheduleRepository) models.AttendanceService {
	return &AttendanceService{
		attendanceRepo: attendanceRepo,
		scheduleRepo:   scheduleRepo,
	}
}

func (as *AttendanceService) CreateAttendance(attendance *models.Attendance) error {
	// Validate schedule exists
	_, err := as.scheduleRepo.GetScheduleByID(attendance.ScheduleID)
	if err != nil {
		return fmt.Errorf("schedule not found: %w", err)
	}

	// Validate attendance data
	if attendance.StudentID == "" {
		return fmt.Errorf("student ID is required")
	}

	if attendance.ScheduleID == "" {
		return fmt.Errorf("schedule ID is required")
	}

	// Check if attendance already exists for this student and schedule
	existingAttendances, err := as.attendanceRepo.GetAttendanceByStudentID(attendance.StudentID)
	if err != nil {
		return fmt.Errorf("failed to check existing attendance: %w", err)
	}

	for _, existing := range existingAttendances {
		if existing.ScheduleID == attendance.ScheduleID {
			return fmt.Errorf("attendance already exists for this student and schedule")
		}
	}

	return as.attendanceRepo.CreateAttendance(attendance)
}

func (as *AttendanceService) GetAttendanceByID(id string) (*models.Attendance, error) {
	if id == "" {
		return nil, fmt.Errorf("attendance ID is required")
	}

	return as.attendanceRepo.GetAttendanceByID(id)
}

func (as *AttendanceService) UpdateAttendance(attendance *models.Attendance) error {
	// Validate attendance exists
	existing, err := as.attendanceRepo.GetAttendanceByID(attendance.ID)
	if err != nil {
		return fmt.Errorf("attendance not found: %w", err)
	}

	// Validate schedule exists if changed
	if existing.ScheduleID != attendance.ScheduleID {
		_, err := as.scheduleRepo.GetScheduleByID(attendance.ScheduleID)
		if err != nil {
			return fmt.Errorf("schedule not found: %w", err)
		}
	}

	// Validate required fields
	if attendance.StudentID == "" {
		return fmt.Errorf("student ID is required")
	}

	if attendance.ScheduleID == "" {
		return fmt.Errorf("schedule ID is required")
	}

	// Update counter logic - increment if marking as present
	if !existing.Here && attendance.Here {
		attendance.Counter = existing.Counter + 1
	} else if existing.Here && !attendance.Here {
		attendance.Counter = existing.Counter - 1
		if attendance.Counter < 0 {
			attendance.Counter = 0
		}
	} else {
		attendance.Counter = existing.Counter
	}

	return as.attendanceRepo.UpdateAttendance(attendance)
}

func (as *AttendanceService) DeleteAttendance(id string) error {
	if id == "" {
		return fmt.Errorf("attendance ID is required")
	}

	// Validate attendance exists
	_, err := as.attendanceRepo.GetAttendanceByID(id)
	if err != nil {
		return fmt.Errorf("attendance not found: %w", err)
	}

	return as.attendanceRepo.DeleteAttendance(id)
}

func (as *AttendanceService) GetAttendanceByStudentID(studentID string) ([]models.Attendance, error) {
	if studentID == "" {
		return nil, fmt.Errorf("student ID is required")
	}

	return as.attendanceRepo.GetAttendanceByStudentID(studentID)
}

func (as *AttendanceService) GetAttendanceByScheduleID(scheduleID string) ([]models.Attendance, error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("schedule ID is required")
	}

	// Validate schedule exists
	_, err := as.scheduleRepo.GetScheduleByID(scheduleID)
	if err != nil {
		return nil, fmt.Errorf("schedule not found: %w", err)
	}

	return as.attendanceRepo.GetAttendanceByScheduleID(scheduleID)
}

// Additional business methods

func (as *AttendanceService) GetAttendanceRateByStudent(studentID string) (float64, error) {
	if studentID == "" {
		return 0, fmt.Errorf("student ID is required")
	}

	attendances, err := as.attendanceRepo.GetAttendanceByStudentID(studentID)
	if err != nil {
		return 0, fmt.Errorf("failed to get student attendances: %w", err)
	}

	if len(attendances) == 0 {
		return 0, nil
	}

	totalPresent := 0
	for _, attendance := range attendances {
		if attendance.Here {
			totalPresent++
		}
	}

	return float64(totalPresent) / float64(len(attendances)) * 100, nil
}

func (as *AttendanceService) MarkAttendance(studentID, scheduleID string, isPresent bool) error {
	if studentID == "" || scheduleID == "" {
		return fmt.Errorf("student ID and schedule ID are required")
	}

	// Check if attendance already exists
	existingAttendances, err := as.attendanceRepo.GetAttendanceByStudentID(studentID)
	if err != nil {
		return fmt.Errorf("failed to check existing attendance: %w", err)
	}

	for _, existing := range existingAttendances {
		if existing.ScheduleID == scheduleID {
			// Update existing attendance
			existing.Here = isPresent
			return as.UpdateAttendance(&existing)
		}
	}

	// Create new attendance record
	newAttendance := &models.Attendance{
		StudentID:  studentID,
		ScheduleID: scheduleID,
		Here:       isPresent,
		Counter:    0,
	}

	if isPresent {
		newAttendance.Counter = 1
	}

	return as.CreateAttendance(newAttendance)
}