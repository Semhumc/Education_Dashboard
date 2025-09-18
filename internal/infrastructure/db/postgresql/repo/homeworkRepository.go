package repo

import (
	"Education_Dashboard/internal/helper"
	"Education_Dashboard/internal/infrastructure/db/postgresql/sqlc/tutorial"
	"Education_Dashboard/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HomeworkRepository struct {
	db      *pgxpool.Pool
	queries *tutorial.Queries
}

func NewHomeworkRepository(db *pgxpool.Pool) models.HomeworkRepository {
	return &HomeworkRepository{
		db:      db,
		queries: tutorial.New(db),
	}
}

func (hr HomeworkRepository) CreateHomework(homework *models.Homework) error {
	ctx := context.Background()

	teacherID, err := helper.ConvertStringToUUID(homework.TeacherID)
	if err != nil {
		return fmt.Errorf("invalid teacher ıd:%w", err)
	}

	lessonID, err := helper.ConvertStringToUUID(homework.LessonID)
	if err != nil {
		return fmt.Errorf("invalid lesson ıd:%w", err)
	}

	classID, err := helper.ConvertStringToUUID(homework.ClassID)
	if err != nil {
		return fmt.Errorf("invalid class ıd:%w", err)
	}
	hwparams := tutorial.CreateHomeworkParams{
		TeacherID: teacherID,
		LessonID:  lessonID,
		ClassID:   classID,
		Title:     homework.Title,
		Content:   pgtype.Text{String: homework.Content, Valid: homework.Content != ""},
		DueDate:   pgtype.Timestamp{Time: homework.DueDate, Valid: true},
	}

	result, err := hr.queries.CreateHomework(ctx, hwparams)
	if err != nil {
		return fmt.Errorf("failed to create homework: %w", err)
	}

	homework.ID = helper.ConvertUUIDToString(result.ID)

	return nil
}

func (hr HomeworkRepository) GetHomeworkByID(id string) (*models.Homework, error) {
	ctx := context.Background()

	homeworkID, err := helper.ConvertStringToUUID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid homework ID: %w", err)
	}

	result, err := hr.queries.GetHomeworkByID(ctx, homeworkID)
	if err != nil {
		return nil, fmt.Errorf("failed to get homework: %w", err)
	}

	homework := &models.Homework{
		ID:        helper.ConvertUUIDToString(result.ID),
		TeacherID: helper.ConvertUUIDToString(result.TeacherID),
		LessonID:  helper.ConvertUUIDToString(result.LessonID),
		ClassID:   helper.ConvertUUIDToString(result.ClassID),
		Title:     result.Title,
		Content:   result.Content.String,
		DueDate:   result.DueDate.Time,
	}

	return homework, nil
}

func (hr HomeworkRepository) UpdateHomework(homework *models.Homework) error {
	ctx := context.Background()

	homeworkID, err := helper.ConvertStringToUUID(homework.ID)
	if err != nil {
		return fmt.Errorf("invalid homework ID: %w", err)
	}

	teacherID, err := helper.ConvertStringToUUID(homework.TeacherID)
	if err != nil {
		return fmt.Errorf("invalid teacher ID: %w", err)
	}

	lessonID, err := helper.ConvertStringToUUID(homework.LessonID)
	if err != nil {
		return fmt.Errorf("invalid lesson ID: %w", err)
	}

	classID, err := helper.ConvertStringToUUID(homework.ClassID)
	if err != nil {
		return fmt.Errorf("invalid class ID: %w", err)
	}

	params := tutorial.UpdateHomeworkParams{
		ID:        homeworkID,
		TeacherID: teacherID,
		LessonID:  lessonID,
		ClassID:   classID,
		Title:     homework.Title,
		Content:   pgtype.Text{String: homework.Content, Valid: homework.Content != ""},
		DueDate:   pgtype.Timestamp{Time: homework.DueDate, Valid: true},
	}

	_, err = hr.queries.UpdateHomework(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update homework: %w", err)
	}

	return nil
}

func (hr HomeworkRepository) DeleteHomework(id string) error {
	ctx := context.Background()
	hwid,err := helper.ConvertStringToUUID(id)
	if err != nil{
		return fmt.Errorf("invalid hw id:%w",err)
	}
	err = hr.queries.DeleteHomework(ctx,hwid)
	if err != nil{
		return fmt.Errorf("delete hw fail:%w",err)
	}
	return nil
}

func (hr HomeworkRepository) GetAllHomeworks() ([]models.Homework, error) {
	ctx := context.Background()

	results, err := hr.queries.GetAllHomeworks(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all homeworks: %w", err)
	}

	var homeworks []models.Homework
	for _, result := range results {
		homework := models.Homework{
			ID:        helper.ConvertUUIDToString(result.ID),
			TeacherID: helper.ConvertUUIDToString(result.TeacherID),
			LessonID:  helper.ConvertUUIDToString(result.LessonID),
			ClassID:   helper.ConvertUUIDToString(result.ClassID),
			Title:     result.Title,
			Content:   result.Content.String,
			DueDate:   result.DueDate.Time,
		}
		homeworks = append(homeworks, homework)
	}

	return homeworks, nil
}

func (hr HomeworkRepository) GetHomeworksByTeacherID(teacherID string) ([]models.Homework, error) {
	ctx := context.Background()
	teacherUUID, err := helper.ConvertStringToUUID(teacherID)
	if err != nil {
		return nil, fmt.Errorf("invalid teacher ID: %w", err)
	}

	results, err := hr.queries.GetHomeworksByTeacherID(ctx, teacherUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get homeworks by teacher ID: %w", err)
	}

	var homeworks []models.Homework
	for _, result := range results {
		homework := models.Homework{
			ID:        helper.ConvertUUIDToString(result.ID),
			TeacherID: helper.ConvertUUIDToString(result.TeacherID),
			LessonID:  helper.ConvertUUIDToString(result.LessonID),
			ClassID:   helper.ConvertUUIDToString(result.ClassID),
			Title:     result.Title,
			Content:   result.Content.String,
			DueDate:   result.DueDate.Time,
		}
		homeworks = append(homeworks, homework)
	}

	return homeworks, nil
}

func (hr HomeworkRepository) GetHomeworksByLessonID(lessonID string) ([]models.Homework, error) {
	ctx := context.Background()

	lessonUUID, err := helper.ConvertStringToUUID(lessonID)
	if err != nil {
		return nil, fmt.Errorf("invalid lesson ID: %w", err)
	}

	results, err := hr.queries.GetHomeworksByLessonID(ctx, lessonUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get homeworks by lesson ID: %w", err)
	}

	var homeworks []models.Homework
	for _, result := range results {
		homework := models.Homework{
			ID:        helper.ConvertUUIDToString(result.ID),
			TeacherID: helper.ConvertUUIDToString(result.TeacherID),
			LessonID:  helper.ConvertUUIDToString(result.LessonID),
			ClassID:   helper.ConvertUUIDToString(result.ClassID),
			Title:     result.Title,
			Content:   result.Content.String,
			DueDate:   result.DueDate.Time,
		}
		homeworks = append(homeworks, homework)
	}

	return homeworks, nil
}

func (hr HomeworkRepository) GetHomeworksByClassID(classID string) ([]models.Homework, error) {
	ctx := context.Background()

	classUUID, err := helper.ConvertStringToUUID(classID)
	if err != nil {
		return nil, fmt.Errorf("invalid class ID: %w", err)
	}

	results, err := hr.queries.GetHomeworksByClassID(ctx, classUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get homeworks by class ID: %w", err)
	}

	var homeworks []models.Homework
	for _, result := range results {
		homework := models.Homework{
			ID:        helper.ConvertUUIDToString(result.ID),
			TeacherID: helper.ConvertUUIDToString(result.TeacherID),
			LessonID:  helper.ConvertUUIDToString(result.LessonID),
			ClassID:   helper.ConvertUUIDToString(result.ClassID),
			Title:     result.Title,
			Content:   result.Content.String,
			DueDate:   result.DueDate.Time,
		}
		homeworks = append(homeworks, homework)
	}

	return homeworks, nil
}


