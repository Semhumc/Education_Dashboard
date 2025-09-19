package repo

import (
	"Education_Dashboard/internal/helper"
	"Education_Dashboard/internal/infrastructure/db/postgresql/sqlc/tutorial"
	"Education_Dashboard/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LessonRepository struct {
	db      *pgxpool.Pool
	queries *tutorial.Queries
}

func NewLessonRepository(db *pgxpool.Pool) models.LessonRepository {
	return &LessonRepository{
		db:      db,
		queries: tutorial.New(db),
	}
}

func (lr *LessonRepository) CreateLesson(lesson *models.Lesson) error {
	ctx := context.Background()

	res, err := lr.queries.CreateLesson(ctx, lesson.LessonName)
	if err != nil {
		return fmt.Errorf("create lesson fail:%w", err)
	}

	lesson.ID = helper.ConvertUUIDToString(res.ID)
	return nil
}

func (lr *LessonRepository) GetLessonByID(id string) (*models.Lesson, error) {
	ctx := context.Background()

	lessonID, err := helper.ConvertStringToUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid lesson ID: %w", err)
	}

	result, err := lr.queries.GetLessonByID(ctx, lessonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lesson: %w", err)
	}

	lesson := &models.Lesson{
		ID:         helper.ConvertUUIDToString(result.ID),
		LessonName: result.LessonName,
	}

	return lesson, nil
}

func (lr *LessonRepository) UpdateLesson(lesson *models.Lesson) error {
	ctx := context.Background()

	lessonID, err := helper.ConvertStringToUUID(lesson.ID)
	if err != nil{
		return fmt.Errorf("invalid lesson id:%w",err)
	}
	params := tutorial.UpdateLessonParams{
		ID: lessonID,
		LessonName: lesson.LessonName,
	}

	_ ,err = lr.queries.UpdateLesson(ctx,params)
	if err != nil{
		return fmt.Errorf("update lesson fail:%w",err)
	}

	return nil
}

func (lr *LessonRepository) DeleteLesson(id string) error {
	ctx := context.Background()
	lessonID, err := helper.ConvertStringToUUID(id)
	if err != nil{
		return fmt.Errorf("invalid lesson id:%w",err)
	}
	err = lr.queries.DeleteLesson(ctx,lessonID)
	if err != nil{
		return fmt.Errorf("delete lesson fail:%w",err)
	}
	return nil
}

func (lr *LessonRepository) GetAllLessons() ([]models.Lesson, error) {
	ctx := context.Background()

	results, err := lr.queries.GetAllLessons(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all lessons: %w", err)
	}

	var lessons []models.Lesson
	for _, result := range results {
		lesson := models.Lesson{
			ID:         helper.ConvertUUIDToString(result.ID),
			LessonName: result.LessonName,
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}