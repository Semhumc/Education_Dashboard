package repo

import (
	"Education_Dashboard/internal/infrastructure/db/postgresql/sqlc/tutorial"
	"Education_Dashboard/internal/models"

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


func(lr *LessonRepository) CreateLesson(lesson *models.Lesson) error{
	return nil
}

func(lr *LessonRepository)GetLessonByID(id string) (*models.Lesson, error){
	return nil,nil
}

func(lr *LessonRepository)UpdateLesson(lesson *models.Lesson) error{
	return nil
}

func(lr *LessonRepository)DeleteLesson(id string) error{
	return nil
}


func(lr *LessonRepository) GetAllLessons() ([]models.Lesson, error){
	return nil,nil
}