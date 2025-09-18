package repo

import (
	"Education_Dashboard/internal/helper"
	"Education_Dashboard/internal/infrastructure/db/postgresql/sqlc/tutorial"
	"Education_Dashboard/internal/models"
	"context"
	"fmt"
	"time"

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

	
	return nil
}
func (sr *SchuedleRepository) DeleteSchedule(id string) error {
	return nil
}
func (sr *SchuedleRepository) GetAllSchedules() ([]models.Schedule, error) {
	return nil, nil
}
func (sr *SchuedleRepository) GetSchedulesByTeacherID(teacherID string) ([]models.Schedule, error) {
	return nil, nil
}
func (sr *SchuedleRepository) GetSchedulesByClassID(classID string) ([]models.Schedule, error) {
	return nil, nil
}
