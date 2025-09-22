package main

import (
	"Education_Dashboard/internal/application"
	"Education_Dashboard/internal/application/handlers"
	"Education_Dashboard/internal/infrastructure/db/postgresql/repo"

	"Education_Dashboard/internal/infrastructure/http"
	"Education_Dashboard/internal/infrastructure/http/handler"
	"Education_Dashboard/internal/infrastructure/http/middleware"
	"Education_Dashboard/internal/infrastructure/keycloak"
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
	keycloak_base_url      string
	keycloak_realm         string
	keycloak_client_id     string
	keycloak_client_secret string
	port                   string

	// Database connection variables
	db_host          string
	db_port          string
	db_name          string
	db_user          string
	db_password      string
	app_frontend_url string
)

func init() {
	// Set environment variables
	keycloak_base_url = os.Getenv("KEYCLOAK_BASE_URL")
	keycloak_realm = os.Getenv("KEYCLOAK_REALM")
	keycloak_client_id = os.Getenv("KEYCLOAK_CLIENT_ID")
	keycloak_client_secret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	port = os.Getenv("APP_PORT")
	if port == "" {
		port = "5000" // Default to 5000 if APP_PORT is not set
	}

	app_frontend_url = os.Getenv("APP_FRONTEND_URL")

	db_host = os.Getenv("DB_POSTGRES_HOST")
	db_port = os.Getenv("DB_POSTGRES_PORT")
	db_name = os.Getenv("DB_POSTGRES_NAME")
	db_user = os.Getenv("DB_POSTGRES_USER")
	db_password = os.Getenv("DB_POSTGRES_PASSWORD")
}

func main() {
	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	// Add middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     app_frontend_url, // Allows requests from this origin
		AllowMethods:     "*",              // Allow all methods
		AllowHeaders:     "*",              // Allow all headers (e.g., Authorization, Content-Type)
		AllowCredentials: true,             // Allow cookies to be sent
	}))

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

	// Initialize Auth Middleware
	authMiddleware := middleware.NewAuthMiddleware(keycloakAuthService)

	// Setup routes
	http.SetupRoutes(app, authHandler, classHandler, scheduleHandler, attendanceHandler, lessonHandler, homeworkHandler, authMiddleware)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
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
