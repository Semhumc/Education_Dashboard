package main

import (
	"Education_Dashboard/internal/application"
	"Education_Dashboard/internal/infrastructure/http"
	"Education_Dashboard/internal/infrastructure/http/handler"
	"Education_Dashboard/internal/application/handlers"
	"Education_Dashboard/internal/infrastructure/keycloak"
	"Education_Dashboard/internal/infrastructure/db/postgresql/repo"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	keycloak_base_url      = os.Getenv("KEYCLOAK_BASE_URL")
	keycloak_realm         = os.Getenv("KEYCLOAK_REALM")
	keycloak_client_id     = os.Getenv("KEYCLOAK_CLIENT_ID")
	keycloak_client_secret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	port                   = "5000"
	
	// Database connection variables
	db_host     = os.Getenv("DB_POSTGRES_HOST")
	db_port     = os.Getenv("DB_POSTGRES_PORT")
	db_name     = os.Getenv("DB_POSTGRES_NAME")
	db_user     = os.Getenv("DB_POSTGRES_USER")
	db_password = os.Getenv("DB_POSTGRES_PASSWORD")
)

func main() {
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": true,
				"message": err.Error(),
			})
		},
	})

	// Add middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize database connection
	dbPool, err := initializeDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer dbPool.Close()

	// Initialize repositories
	attendanceRepo := repo.NewAttendanceRepository(dbPool)
	homeworkRepo := repo.NewHomeworkRepository(dbPool)
	lessonRepo := repo.NewLessonRepository(dbPool)
	scheduleRepo := repo.NewSchuedleRepository(dbPool)

	// Initialize application services
	attendanceService := application.NewAttendanceService(attendanceRepo, scheduleRepo)
	homeworkService := application.NewHomeworkService(homeworkRepo, lessonRepo)
	lessonService := application.NewLessonService(lessonRepo, homeworkRepo, scheduleRepo)
	scheduleService := application.NewScheduleService(scheduleRepo, lessonRepo, attendanceRepo)

	// Initialize Keycloak service
	keycloakAuthService, keycloakClassService := keycloak.NewKeycloakAuthService(
		keycloak_base_url,
		keycloak_client_id,
		keycloak_client_secret,
		keycloak_realm,
	)

	// Initialize handlers
	authHandler := handler.NewKeycloakHandler(keycloakAuthService)
	classHandler := handler.NewKeycloakClassHandler(keycloakClassService)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceService)
	homeworkHandler := handlers.NewHomeworkHandler(homeworkService)
	lessonHandler := handlers.NewLessonHandler(lessonService)
	scheduleHandler := handlers.NewScheduleHandler(scheduleService)

	// Setup routes
	setupRoutes(app, authHandler, classHandler, attendanceHandler, *homeworkHandler, *lessonHandler, *scheduleHandler)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Education Dashboard API is running",
		})
	})

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initializeDatabase() (*pgxpool.Pool, error) {
	// Build connection string
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		db_user, db_password, db_host, db_port, db_name,
	)

	// Create connection pool
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set connection pool settings
	config.MaxConns = 10
	config.MinConns = 2

	// Connect to database
	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established successfully")
	return pool, nil
}

func setupRoutes(app *fiber.App, authHandler handler.KeycloakHandler, classHandler handler.KeycloakClassHandler, attendanceHandler handlers.AttendanceHandler, homeworkHandler handlers.HomeworkHandler, lessonHandler handlers.LessonHandler, scheduleHandler handlers.ScheduleHandler) {
	// Authentication routes
	http.AuthRoutes(app, authHandler)
	
	// API group
	api := app.Group("/v1/api")
	
	// Class routes
	classGroup := api.Group("/class")
	classGroup.Post("/create", classHandler.CreateClassHandler)
	classGroup.Get("/all", classHandler.GetAllClassesHandler)
	classGroup.Get("/teacher/:teacherID", classHandler.GetClassesByTeacherIDHandler)
	classGroup.Put("/update/:id", classHandler.UpdateClassHandler)
	classGroup.Delete("/delete/:id", classHandler.DeleteClassHandler)
	
	// Attendance routes
	attendanceGroup := api.Group("/attendance")
	attendanceGroup.Post("/create", attendanceHandler.CreateAttendanceHandler)
	attendanceGroup.Get("/:id", attendanceHandler.GetAttendanceByIDHandler)
	attendanceGroup.Put("/update/:id", attendanceHandler.UpdateAttendanceHandler)
	attendanceGroup.Delete("/delete/:id", attendanceHandler.DeleteAttendanceHandler)
	attendanceGroup.Get("/student/:studentID", attendanceHandler.GetAttendanceByStudentIDHandler)
	attendanceGroup.Get("/schedule/:scheduleID", attendanceHandler.GetAttendanceByScheduleIDHandler)
	attendanceGroup.Post("/mark", attendanceHandler.MarkAttendanceHandler)
	//attendanceGroup.Get("/rate/:studentID", attendanceHandler.GetAttendanceRateHandler)
	
	// Homework routes
	homeworkGroup := api.Group("/homework")
	homeworkGroup.Post("/create", homeworkHandler.CreateHomeworkHandler)
	homeworkGroup.Get("/:id", homeworkHandler.GetHomeworkByIDHandler)
	homeworkGroup.Put("/update/:id", homeworkHandler.UpdateHomeworkHandler)
	homeworkGroup.Delete("/delete/:id", homeworkHandler.DeleteHomeworkHandler)
	homeworkGroup.Get("/all", homeworkHandler.GetAllHomeworksHandler)
	homeworkGroup.Get("/teacher/:teacherID", homeworkHandler.GetHomeworksByTeacherIDHandler)
	homeworkGroup.Get("/lesson/:lessonID", homeworkHandler.GetHomeworksByLessonIDHandler)
	homeworkGroup.Get("/class/:classID", homeworkHandler.GetHomeworksByClassIDHandler)
	homeworkGroup.Get("/active", homeworkHandler.GetActiveHomeworksHandler)
	homeworkGroup.Get("/overdue", homeworkHandler.GetOverdueHomeworksHandler)
	homeworkGroup.Get("/due-soon", homeworkHandler.GetHomeworksDueSoonHandler)
	homeworkGroup.Put("/extend/:id", homeworkHandler.ExtendDueDateHandler)
	
	// Lesson routes
	lessonGroup := api.Group("/lesson")
	lessonGroup.Post("/create", lessonHandler.CreateLessonHandler)
	lessonGroup.Get("/:id", lessonHandler.GetLessonByIDHandler)
	lessonGroup.Put("/update/:id", lessonHandler.UpdateLessonHandler)
	lessonGroup.Delete("/delete/:id", lessonHandler.DeleteLessonHandler)
	lessonGroup.Get("/all", lessonHandler.GetAllLessonsHandler)
	//lessonGroup.Get("/search", lessonHandler.SearchLessonsHandler)
	//lessonGroup.Get("/stats/:id", lessonHandler.GetLessonStatsHandler)
	//lessonGroup.Get("/homework-count", lessonHandler.GetLessonsWithHomeworkCountHandler)
	
	// Schedule routes
	scheduleGroup := api.Group("/schedule")
	scheduleGroup.Post("/create", scheduleHandler.CreateScheduleHandler)
	scheduleGroup.Get("/:id", scheduleHandler.GetScheduleByIDHandler)
	scheduleGroup.Put("/update/:id", scheduleHandler.UpdateScheduleHandler)
	scheduleGroup.Delete("/delete/:id", scheduleHandler.DeleteScheduleHandler)
	scheduleGroup.Get("/all", scheduleHandler.GetAllSchedulesHandler)
	scheduleGroup.Get("/teacher/:teacherID", scheduleHandler.GetSchedulesByTeacherIDHandler)
	scheduleGroup.Get("/class/:classID", scheduleHandler.GetSchedulesByClassIDHandler)
	scheduleGroup.Get("/today", scheduleHandler.GetTodaySchedulesHandler)
	scheduleGroup.Get("/week", scheduleHandler.GetWeekSchedulesHandler)
	scheduleGroup.Get("/upcoming/:teacherID", scheduleHandler.GetUpcomingSchedulesHandler)
	scheduleGroup.Post("/check-conflicts", scheduleHandler.CheckConflictsHandler)
	scheduleGroup.Put("/reschedule/:id", scheduleHandler.RescheduleHandler)
}