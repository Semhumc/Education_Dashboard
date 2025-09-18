package repo

import (
	"Education_Dashboard/internal/helper"
	"Education_Dashboard/internal/infrastructure/db/postgresql/sqlc/tutorial"
	"Education_Dashboard/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AttendanceRepository struct {
	db      *pgxpool.Pool
	queries *tutorial.Queries
}

func NewAttendanceRepository(db *pgxpool.Pool) models.AttendanceRepository {
	return &AttendanceRepository{
		db:      db,
		queries: tutorial.New(db),
	}
}

func (ar AttendanceRepository) CreateAttendance(attendance *models.Attendance) error {
	ctx := context.Background()
	studentID, err := helper.ConvertStringToUUID(attendance.StudentID)
	if err != nil {
		return fmt.Errorf("invalid student id :%w", err)
	}

	scheduleID, err := helper.ConvertStringToUUID(attendance.ScheduleID)
	if err != nil {
		return fmt.Errorf("invalid schudle id :%w", err)
	}

	params := tutorial.CreateAttendanceParams{
		StudentID:  studentID,
		ScheduleID: scheduleID,
		Here:       attendance.Here,
		Counter:    int32(attendance.Counter),
	}

	result, err := ar.queries.CreateAttendance(ctx, params)
	if err != nil {
		return fmt.Errorf("attendance create failed :%w", err)

	}

	attendance.ID = helper.ConvertUUIDToString(result.ID)
	return nil
}

func (ar AttendanceRepository) GetAttendanceByID(id string) (*models.Attendance, error) {
	ctx := context.Background()
	attendanceID, err := helper.ConvertStringToUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid attendance id:%w", err)
	}

	res, err := ar.queries.GetAttendanceByID(ctx, attendanceID)
	if err != nil {
		return nil, fmt.Errorf("get attendance by id fail:%w", err)
	}

	attendance := &models.Attendance{
		ID:        helper.ConvertUUIDToString(res.ID),
		StudentID: helper.ConvertUUIDToString(res.StudentID),
		Here:      res.Here,
		Counter:   int(res.Counter),
	}

	return attendance, nil
}

func (ar AttendanceRepository) UpdateAttendance(attendance *models.Attendance) error {
	ctx := context.Background()

	
	attendanceID, err := helper.ConvertStringToUUID(attendance.ID)
	if err != nil {
		return fmt.Errorf("invalid attendance ID: %w", err)
	}

	studentID, err := helper.ConvertStringToUUID(attendance.StudentID)
	if err != nil {
		return fmt.Errorf("invalid student ID: %w", err)
	}

	scheduleID, err := helper.ConvertStringToUUID(attendance.ScheduleID)
	if err != nil {
		return fmt.Errorf("invalid schedule ID: %w", err)
	}

	params := tutorial.UpdateAttendanceParams{
		ID:         attendanceID,
		StudentID:  studentID,
		ScheduleID: scheduleID,
		Here:       attendance.Here,
		Counter:    int32(attendance.Counter),
	}

	_,err = ar.queries.UpdateAttendance(ctx,params)
	if err != nil {
		return fmt.Errorf("failed to update attendance: %w", err)
	}
	return nil
}

func (ar AttendanceRepository) DeleteAttendance(id string) error {
	ctx := context.Background()
	attendanceID, err := helper.ConvertStringToUUID(id)
	if err != nil {
		return fmt.Errorf("invalid attendance id :%w", err)
	}

	err = ar.queries.DeleteAttendance(ctx, attendanceID)
	if err != nil {
		return fmt.Errorf("delete attendance fail:%w", err)
	}
	return nil
}

func (ar AttendanceRepository) GetAttendanceByStudentID(studentID string) ([]models.Attendance, error) {
	ctx := context.Background()
	studentUUID, err := helper.ConvertStringToUUID(studentID)
	if err != nil {
		return nil, fmt.Errorf("invalid student uuid:%w", err)
	}

	res, err := ar.queries.GetAttendanceByStudentID(ctx, studentUUID)
	if err != nil {
		return nil, fmt.Errorf("getAttendanceByStudentID fail:%w", err)
	}

	var attendances []models.Attendance
	for _, result := range res {
		attendance := models.Attendance{
			ID:         helper.ConvertUUIDToString(result.ID),
			StudentID:  helper.ConvertUUIDToString(result.StudentID),
			ScheduleID: helper.ConvertUUIDToString(result.ScheduleID),
			Here:       result.Here,
			Counter:    int(result.Counter),
		}
		attendances = append(attendances, attendance)
	}
	return attendances, nil
}
func (ar AttendanceRepository) GetAttendanceByScheduleID(scheduleID string) ([]models.Attendance, error) {
	ctx := context.Background()
	scheduleUUID, err := helper.ConvertStringToUUID(scheduleID)
	if err != nil {
		return nil, fmt.Errorf("invalid schedule ID: %w", err)
	}

	res,err := ar.queries.GetAttendanceByScheduleID(ctx,scheduleUUID)
	if err != nil{
		return nil,fmt.Errorf("getAttendanceByScheduleID failed : %w",err)
	}

	var attendances []models.Attendance
	for _, result := range res {
		attendance := models.Attendance{
			ID:         helper.ConvertUUIDToString(result.ID),
			StudentID:  helper.ConvertUUIDToString(result.StudentID),
			ScheduleID: helper.ConvertUUIDToString(result.ScheduleID),
			Here:       result.Here,
			Counter:    int(result.Counter),
		}
		attendances = append(attendances, attendance)
	}

	return attendances, nil
}
