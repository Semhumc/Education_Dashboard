package repo

import (
	"Education_Dashboard/internal/helper"
	"Education_Dashboard/internal/infrastructure/db/postgresql/sqlc/tutorial"
	"Education_Dashboard/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SchuedleRepository struct {
	db      *pgxpool.Pool
	queries *tutorial.Queries
}

func NewSchuedleRepository(db *pgxpool.Pool) models.ScheduleRepository {
	return &SchuedleRepository{
		db:      db,
		queries: tutorial.New(db),
	}
}

func (sr *SchuedleRepository) CreateSchedule(schedule *models.Schedule) error {
	ctx := context.Background()
	lessonID, err := helper.ConvertStringToUUID(schedule.LessonID)
	if err != nil {
		return fmt.Errorf("invalid lesson id:%w", err)
	}

	teacherID, err := helper.ConvertStringToUUID(schedule.TeacherID)
	if err != nil {
		return fmt.Errorf("invalid teacher id:%w", err)
	}

	classID, err := helper.ConvertStringToUUID(schedule.ClassID)
	if err != nil {
		return fmt.Errorf("invalid class id:%w", err)
	}
	params := tutorial.CreateScheduleParams{
		Date:      pgtype.Date{Time: schedule.Date, Valid: true},
		Time:      pgtype.Time{Microseconds: int64(schedule.Time.Hour()*3600000000 + schedule.Time.Minute()*60000000), Valid: true},
		TeacherID: teacherID,
		LessonID:  lessonID,
		ClassID:   classID,
	}

	res, err := sr.queries.CreateSchedule(ctx, params)
	if err != nil {
		return fmt.Errorf("create schuedle fail:%w", err)
	}
	schedule.ID = helper.ConvertUUIDToString(res.ID)
	return nil
}

func (sr *SchuedleRepository) GetScheduleByID(id string) (*models.Schedule, error) {
	ctx := context.Background()
	schuedleID, err := helper.ConvertStringToUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid schuedle id: %w", err)
	}

	res, err := sr.queries.GetScheduleByID(ctx, schuedleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	schedule := &models.Schedule{
		ID:        helper.ConvertUUIDToString(res.ID),
		Date:      res.Date.Time,
		Time:      time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(res.Time.Microseconds) * time.Microsecond),
		TeacherID: helper.ConvertUUIDToString(res.TeacherID),
		LessonID:  helper.ConvertUUIDToString(res.LessonID),
		ClassID:   helper.ConvertUUIDToString(res.ClassID),
	}

	return schedule, nil
}
func (sr *SchuedleRepository) UpdateSchedule(schedule *models.Schedule) error {

	ctx := context.Background()
	schuedleID, err := helper.ConvertStringToUUID(schedule.ID)
	if err != nil {
		return fmt.Errorf("invalid schuedle id:%w", err)
	}

	lessonID, err := helper.ConvertStringToUUID(schedule.LessonID)
	if err != nil {
		return fmt.Errorf("invalid lesson id:%w", err)
	}

	teacherID, err := helper.ConvertStringToUUID(schedule.TeacherID)
	if err != nil {
		return fmt.Errorf("invalid teacher id:%w", err)
	}

	classID, err := helper.ConvertStringToUUID(schedule.ClassID)
	if err != nil {
		return fmt.Errorf("invalid class id:%w", err)
	}

	params := tutorial.UpdateScheduleParams{
		ID:        schuedleID,
		Date:      pgtype.Date{Time: schedule.Date, Valid: true},
		Time:      pgtype.Time{Microseconds: int64(schedule.Time.Hour()*3600000000 + schedule.Time.Minute()*60000000), Valid: true},
		TeacherID: teacherID,
		LessonID:  lessonID,
		ClassID:   classID,
	}

	_, err = sr.queries.UpdateSchedule(ctx, params)
	if err != nil {
		return fmt.Errorf("update schuedle fail:%w", err)
	}
	return nil

}
func (sr *SchuedleRepository) DeleteSchedule(id string) error {
	ctx := context.Background()
	schuedleID, err := helper.ConvertStringToUUID(id)
	if err != nil {
		return fmt.Errorf("schuedle id invalid :%w", err)
	}

	err = sr.queries.DeleteSchedule(ctx, schuedleID)
	if err != nil {
		return fmt.Errorf("delete schuedle fail:%w", err)
	}
	return nil
}
func (sr *SchuedleRepository) GetAllSchedules() ([]models.Schedule, error) {
	ctx := context.Background()

	res, err := sr.queries.GetAllSchedules(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get all schedules: %w", err)
	}

	var schedules []models.Schedule
	for _, result := range res {
		schedule := models.Schedule{
			ID:        helper.ConvertUUIDToString(result.ID),
			Date:      result.Date.Time,
			Time:      time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(result.Time.Microseconds) * time.Microsecond),
			TeacherID: helper.ConvertUUIDToString(result.TeacherID),
			LessonID:  helper.ConvertUUIDToString(result.LessonID),
			ClassID:   helper.ConvertUUIDToString(result.ClassID),
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}
func (sr *SchuedleRepository) GetSchedulesByTeacherID(teacherID string) ([]models.Schedule, error) {
	ctx := context.Background()

	teacherUUID, err := helper.ConvertStringToUUID(teacherID)
	if err != nil {
		return nil, fmt.Errorf("invalid teacher ID: %w", err)
	}

	results, err := sr.queries.GetSchedulesByTeacherID(ctx, teacherUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedules by teacher ID: %w", err)
	}

	var schedules []models.Schedule
	for _, result := range results {
		schedule := models.Schedule{
			ID:        helper.ConvertUUIDToString(result.ID),
			Date:      result.Date.Time,
			Time:      time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(result.Time.Microseconds) * time.Microsecond),
			TeacherID: helper.ConvertUUIDToString(result.TeacherID),
			LessonID:  helper.ConvertUUIDToString(result.LessonID),
			ClassID:   helper.ConvertUUIDToString(result.ClassID),
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (sr *SchuedleRepository) GetSchedulesByClassID(classID string) ([]models.Schedule, error) {
	ctx := context.Background()

	classUUID, err := helper.ConvertStringToUUID(classID)
	if err != nil {
		return nil, fmt.Errorf("invalid class ID: %w", err)
	}

	results, err := sr.queries.GetSchedulesByClassID(ctx, classUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedules by class ID: %w", err)
	}

	var schedules []models.Schedule
	for _, result := range results {
		schedule := models.Schedule{
			ID:        helper.ConvertUUIDToString(result.ID),
			Date:      result.Date.Time,
			Time:      time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC).Add(time.Duration(result.Time.Microseconds) * time.Microsecond),
			TeacherID: helper.ConvertUUIDToString(result.TeacherID),
			LessonID:  helper.ConvertUUIDToString(result.LessonID),
			ClassID:   helper.ConvertUUIDToString(result.ClassID),
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}
