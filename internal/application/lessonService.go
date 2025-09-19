package application

import (
	"Education_Dashboard/internal/models"
	"fmt"
	"strings"
)

type LessonService struct {
	lessonRepo   models.LessonRepository
	homeworkRepo models.HomeworkRepository
	scheduleRepo models.ScheduleRepository
}

func NewLessonService(lessonRepo models.LessonRepository, homeworkRepo models.HomeworkRepository, scheduleRepo models.ScheduleRepository) models.LessonService {
	return &LessonService{
		lessonRepo:   lessonRepo,
		homeworkRepo: homeworkRepo,
		scheduleRepo: scheduleRepo,
	}
}

func (ls *LessonService) CreateLesson(lesson *models.Lesson) error {
	// Validate required fields
	if lesson.LessonName == "" {
		return fmt.Errorf("lesson name is required")
	}

	// Normalize lesson name
	lesson.LessonName = strings.TrimSpace(lesson.LessonName)
	
	if len(lesson.LessonName) < 2 {
		return fmt.Errorf("lesson name must be at least 2 characters long")
	}

	if len(lesson.LessonName) > 255 {
		return fmt.Errorf("lesson name cannot exceed 255 characters")
	}

	// Check for duplicate lesson names
	existingLessons, err := ls.lessonRepo.GetAllLessons()
	if err != nil {
		return fmt.Errorf("failed to check existing lessons: %w", err)
	}

	for _, existing := range existingLessons {
		if strings.EqualFold(existing.LessonName, lesson.LessonName) {
			return fmt.Errorf("lesson with name '%s' already exists", lesson.LessonName)
		}
	}

	return ls.lessonRepo.CreateLesson(lesson)
}

func (ls *LessonService) GetLessonByID(id string) (*models.Lesson, error) {
	if id == "" {
		return nil, fmt.Errorf("lesson ID is required")
	}

	return ls.lessonRepo.GetLessonByID(id)
}

func (ls *LessonService) UpdateLesson(lesson *models.Lesson) error {
	// Validate lesson exists
	existing, err := ls.lessonRepo.GetLessonByID(lesson.ID)
	if err != nil {
		return fmt.Errorf("lesson not found: %w", err)
	}

	// Validate required fields
	if lesson.LessonName == "" {
		return fmt.Errorf("lesson name is required")
	}

	// Normalize lesson name
	lesson.LessonName = strings.TrimSpace(lesson.LessonName)
	
	if len(lesson.LessonName) < 2 {
		return fmt.Errorf("lesson name must be at least 2 characters long")
	}

	if len(lesson.LessonName) > 255 {
		return fmt.Errorf("lesson name cannot exceed 255 characters")
	}

	// Check for duplicate lesson names (excluding current lesson)
	if !strings.EqualFold(existing.LessonName, lesson.LessonName) {
		allLessons, err := ls.lessonRepo.GetAllLessons()
		if err != nil {
			return fmt.Errorf("failed to check existing lessons: %w", err)
		}

		for _, otherLesson := range allLessons {
			if otherLesson.ID != lesson.ID && strings.EqualFold(otherLesson.LessonName, lesson.LessonName) {
				return fmt.Errorf("lesson with name '%s' already exists", lesson.LessonName)
			}
		}
	}

	return ls.lessonRepo.UpdateLesson(lesson)
}

func (ls *LessonService) DeleteLesson(id string) error {
	if id == "" {
		return fmt.Errorf("lesson ID is required")
	}

	// Validate lesson exists
	_, err := ls.lessonRepo.GetLessonByID(id)
	if err != nil {
		return fmt.Errorf("lesson not found: %w", err)
	}

	// Check if lesson has associated homeworks
	homeworks, err := ls.homeworkRepo.GetHomeworksByLessonID(id)
	if err != nil {
		return fmt.Errorf("failed to check associated homeworks: %w", err)
	}

	if len(homeworks) > 0 {
		return fmt.Errorf("cannot delete lesson with associated homeworks (%d found). Please delete or reassign homeworks first", len(homeworks))
	}

	// Check if lesson has associated schedules
	schedules, err := ls.scheduleRepo.GetAllSchedules()
	if err != nil {
		return fmt.Errorf("failed to check associated schedules: %w", err)
	}

	var associatedSchedules []models.Schedule
	for _, schedule := range schedules {
		if schedule.LessonID == id {
			associatedSchedules = append(associatedSchedules, schedule)
		}
	}

	if len(associatedSchedules) > 0 {
		return fmt.Errorf("cannot delete lesson with associated schedules (%d found). Please delete or reassign schedules first", len(associatedSchedules))
	}

	return ls.lessonRepo.DeleteLesson(id)
}

func (ls *LessonService) GetAllLessons() ([]models.Lesson, error) {
	return ls.lessonRepo.GetAllLessons()
}

// Additional business methods

func (ls *LessonService) SearchLessonsByName(searchTerm string) ([]models.Lesson, error) {
	if searchTerm == "" {
		return nil, fmt.Errorf("search term is required")
	}

	searchTerm = strings.TrimSpace(strings.ToLower(searchTerm))
	if len(searchTerm) < 2 {
		return nil, fmt.Errorf("search term must be at least 2 characters long")
	}

	allLessons, err := ls.lessonRepo.GetAllLessons()
	if err != nil {
		return nil, fmt.Errorf("failed to get all lessons: %w", err)
	}

	var matchingLessons []models.Lesson
	for _, lesson := range allLessons {
		if strings.Contains(strings.ToLower(lesson.LessonName), searchTerm) {
			matchingLessons = append(matchingLessons, lesson)
		}
	}

	return matchingLessons, nil
}

func (ls *LessonService) GetLessonStats(lessonID string) (*LessonStats, error) {
	if lessonID == "" {
		return nil, fmt.Errorf("lesson ID is required")
	}

	// Validate lesson exists
	lesson, err := ls.lessonRepo.GetLessonByID(lessonID)
	if err != nil {
		return nil, fmt.Errorf("lesson not found: %w", err)
	}

	// Get associated homeworks
	homeworks, err := ls.homeworkRepo.GetHomeworksByLessonID(lessonID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lesson homeworks: %w", err)
	}

	// Get associated schedules
	allSchedules, err := ls.scheduleRepo.GetAllSchedules()
	if err != nil {
		return nil, fmt.Errorf("failed to get schedules: %w", err)
	}

	var scheduleCount int
	for _, schedule := range allSchedules {
		if schedule.LessonID == lessonID {
			scheduleCount++
		}
	}

	stats := &LessonStats{
		LessonID:      lesson.ID,
		LessonName:    lesson.LessonName,
		HomeworkCount: len(homeworks),
		ScheduleCount: scheduleCount,
	}

	return stats, nil
}

func (ls *LessonService) GetLessonsWithHomeworkCount() ([]LessonWithHomeworkCount, error) {
	allLessons, err := ls.lessonRepo.GetAllLessons()
	if err != nil {
		return nil, fmt.Errorf("failed to get all lessons: %w", err)
	}

	var result []LessonWithHomeworkCount
	for _, lesson := range allLessons {
		homeworks, err := ls.homeworkRepo.GetHomeworksByLessonID(lesson.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get homeworks for lesson %s: %w", lesson.ID, err)
		}

		result = append(result, LessonWithHomeworkCount{
			Lesson:        lesson,
			HomeworkCount: len(homeworks),
		})
	}

	return result, nil
}

// Helper structs for additional business methods
type LessonStats struct {
	LessonID      string `json:"lesson_id"`
	LessonName    string `json:"lesson_name"`
	HomeworkCount int    `json:"homework_count"`
	ScheduleCount int    `json:"schedule_count"`
}

type LessonWithHomeworkCount struct {
	Lesson        models.Lesson `json:"lesson"`
	HomeworkCount int           `json:"homework_count"`
}