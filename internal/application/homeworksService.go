package application

import (
	"Education_Dashboard/internal/models"
	"fmt"
	"time"
)

type HomeworkService struct {
	homeworkRepo models.HomeworkRepository
	lessonRepo   models.LessonRepository
}

func NewHomeworkService(homeworkRepo models.HomeworkRepository, lessonRepo models.LessonRepository) models.HomeworkService {
	return &HomeworkService{
		homeworkRepo: homeworkRepo,
		lessonRepo:   lessonRepo,
	}
}

func (hs *HomeworkService) CreateHomework(homework *models.Homework) error {
	// Validate required fields
	if homework.TeacherID == "" {
		return fmt.Errorf("teacher ID is required")
	}

	if homework.LessonID == "" {
		return fmt.Errorf("lesson ID is required")
	}

	if homework.ClassID == "" {
		return fmt.Errorf("class ID is required")
	}

	if homework.Title == "" {
		return fmt.Errorf("homework title is required")
	}

	// Validate lesson exists
	_, err := hs.lessonRepo.GetLessonByID(homework.LessonID)
	if err != nil {
		return fmt.Errorf("lesson not found: %w", err)
	}

	// Validate due date is in the future
	if homework.DueDate.Before(time.Now()) {
		return fmt.Errorf("due date must be in the future")
	}

	return hs.homeworkRepo.CreateHomework(homework)
}

func (hs *HomeworkService) GetHomeworkByID(id string) (*models.Homework, error) {
	if id == "" {
		return nil, fmt.Errorf("homework ID is required")
	}

	return hs.homeworkRepo.GetHomeworkByID(id)
}

func (hs *HomeworkService) UpdateHomework(homework *models.Homework) error {
	// Validate homework exists
	existing, err := hs.homeworkRepo.GetHomeworkByID(homework.ID)
	if err != nil {
		return fmt.Errorf("homework not found: %w", err)
	}

	// Validate required fields
	if homework.TeacherID == "" {
		return fmt.Errorf("teacher ID is required")
	}

	if homework.LessonID == "" {
		return fmt.Errorf("lesson ID is required")
	}

	if homework.ClassID == "" {
		return fmt.Errorf("class ID is required")
	}

	if homework.Title == "" {
		return fmt.Errorf("homework title is required")
	}

	// Validate lesson exists if changed
	if existing.LessonID != homework.LessonID {
		_, err := hs.lessonRepo.GetLessonByID(homework.LessonID)
		if err != nil {
			return fmt.Errorf("lesson not found: %w", err)
		}
	}

	// Check if due date is being changed to past (only if not already past)
	if !existing.DueDate.Before(time.Now()) && homework.DueDate.Before(time.Now()) {
		return fmt.Errorf("cannot set due date to past for active homework")
	}

	return hs.homeworkRepo.UpdateHomework(homework)
}

func (hs *HomeworkService) DeleteHomework(id string) error {
	if id == "" {
		return fmt.Errorf("homework ID is required")
	}

	// Validate homework exists
	homework, err := hs.homeworkRepo.GetHomeworkByID(id)
	if err != nil {
		return fmt.Errorf("homework not found: %w", err)
	}

	// Business rule: Cannot delete homework that is past due date
	if homework.DueDate.Before(time.Now()) {
		return fmt.Errorf("cannot delete homework that is past due date")
	}

	return hs.homeworkRepo.DeleteHomework(id)
}

func (hs *HomeworkService) GetAllHomeworks() ([]models.Homework, error) {
	return hs.homeworkRepo.GetAllHomeworks()
}

func (hs *HomeworkService) GetHomeworksByTeacherID(teacherID string) ([]models.Homework, error) {
	if teacherID == "" {
		return nil, fmt.Errorf("teacher ID is required")
	}

	return hs.homeworkRepo.GetHomeworksByTeacherID(teacherID)
}

func (hs *HomeworkService) GetHomeworksByLessonID(lessonID string) ([]models.Homework, error) {
	if lessonID == "" {
		return nil, fmt.Errorf("lesson ID is required")
	}

	// Validate lesson exists
	_, err := hs.lessonRepo.GetLessonByID(lessonID)
	if err != nil {
		return nil, fmt.Errorf("lesson not found: %w", err)
	}

	return hs.homeworkRepo.GetHomeworksByLessonID(lessonID)
}

func (hs *HomeworkService) GetHomeworksByClassID(classID string) ([]models.Homework, error) {
	if classID == "" {
		return nil, fmt.Errorf("class ID is required")
	}

	return hs.homeworkRepo.GetHomeworksByClassID(classID)
}

// Additional business methods

func (hs *HomeworkService) GetActiveHomeworks() ([]models.Homework, error) {
	allHomeworks, err := hs.homeworkRepo.GetAllHomeworks()
	if err != nil {
		return nil, fmt.Errorf("failed to get all homeworks: %w", err)
	}

	var activeHomeworks []models.Homework
	now := time.Now()

	for _, homework := range allHomeworks {
		if homework.DueDate.After(now) {
			activeHomeworks = append(activeHomeworks, homework)
		}
	}

	return activeHomeworks, nil
}

func (hs *HomeworkService) GetOverdueHomeworks() ([]models.Homework, error) {
	allHomeworks, err := hs.homeworkRepo.GetAllHomeworks()
	if err != nil {
		return nil, fmt.Errorf("failed to get all homeworks: %w", err)
	}

	var overdueHomeworks []models.Homework
	now := time.Now()

	for _, homework := range allHomeworks {
		if homework.DueDate.Before(now) {
			overdueHomeworks = append(overdueHomeworks, homework)
		}
	}

	return overdueHomeworks, nil
}

func (hs *HomeworkService) GetHomeworksDueSoon(hours int) ([]models.Homework, error) {
	if hours <= 0 {
		return nil, fmt.Errorf("hours must be positive")
	}

	allHomeworks, err := hs.homeworkRepo.GetAllHomeworks()
	if err != nil {
		return nil, fmt.Errorf("failed to get all homeworks: %w", err)
	}

	var dueSoonHomeworks []models.Homework
	now := time.Now()
	threshold := now.Add(time.Duration(hours) * time.Hour)

	for _, homework := range allHomeworks {
		if homework.DueDate.After(now) && homework.DueDate.Before(threshold) {
			dueSoonHomeworks = append(dueSoonHomeworks, homework)
		}
	}

	return dueSoonHomeworks, nil
}

func (hs *HomeworkService) ExtendDueDate(homeworkID string, newDueDate time.Time) error {
	if homeworkID == "" {
		return fmt.Errorf("homework ID is required")
	}

	if newDueDate.Before(time.Now()) {
		return fmt.Errorf("new due date must be in the future")
	}

	homework, err := hs.homeworkRepo.GetHomeworkByID(homeworkID)
	if err != nil {
		return fmt.Errorf("homework not found: %w", err)
	}

	if newDueDate.Before(homework.DueDate) {
		return fmt.Errorf("new due date cannot be earlier than current due date")
	}

	homework.DueDate = newDueDate
	return hs.homeworkRepo.UpdateHomework(homework)
}