package http

import (
	"Education_Dashboard/internal/application/handlers"
	"Education_Dashboard/internal/infrastructure/http/handler"
	"Education_Dashboard/internal/infrastructure/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, kh handler.KeycloakHandler, kch handler.KeycloakClassHandler, sh *handlers.ScheduleHandler, ah *handlers.AttendanceHandler, lh *handlers.LessonHandler, hwh *handlers.HomeworkHandler, authMiddleware middleware.AuthMiddleware) {
	api := app.Group("/v1/api")

	// Auth routes
	api.Post("/register", kh.RegisterHandler)
	api.Post("/login", kh.LoginHandler)
	api.Post("/logout", kh.LogoutHandler)

	// User routes
	user := api.Group("/user")
	user.Use(authMiddleware.AuthMiddleware())
	user.Put("/update/:id", authMiddleware.HasRole("admin", "teacher", "student"), kh.UpdateUserHandler)
	user.Delete("/delete/:id", authMiddleware.HasRole("admin"), kh.DeleteUserHandler)
	user.Get("/currentuser/:id", authMiddleware.HasRole("admin", "teacher", "student"), kh.GetUserByIDHandler)
	user.Get("/allusers", authMiddleware.HasRole("admin"), kh.GetAllUsersHandler)
	user.Post("/create", authMiddleware.HasRole("admin"), kh.CreateUserHandler)

	// Class routes
	class := api.Group("/class")
	class.Use(authMiddleware.AuthMiddleware())
	class.Post("/create", authMiddleware.HasRole("admin"), kch.CreateClassHandler)
	class.Delete("/delete/:id", authMiddleware.HasRole("admin"), kch.DeleteClassHandler)
	class.Put("/update", authMiddleware.HasRole("admin"), kch.UpdateClassHandler)
	class.Get("/all", authMiddleware.HasRole("admin", "teacher", "student"), kch.GetAllClassesHandler)
	class.Get("/teacher/:teacherID", authMiddleware.HasRole("teacher"), kch.GetClassesByTeacherIDHandler)
	class.Get("/students/:classID", authMiddleware.HasRole("teacher", "admin"), kch.GetStudentsByClassIDHandler)

	// Schedule routes
	schedule := api.Group("/schedule")
	schedule.Use(authMiddleware.AuthMiddleware())
	schedule.Post("/create", authMiddleware.HasRole("admin", "teacher"), sh.CreateScheduleHandler)
	schedule.Get("/all", authMiddleware.HasRole("admin", "teacher", "student"), sh.GetAllSchedulesHandler)
	schedule.Get("/:id", authMiddleware.HasRole("admin", "teacher", "student"), sh.GetScheduleByIDHandler)
	schedule.Put("/update/:id", authMiddleware.HasRole("admin", "teacher"), sh.UpdateScheduleHandler)
	schedule.Delete("/delete/:id", authMiddleware.HasRole("admin", "teacher"), sh.DeleteScheduleHandler)
	schedule.Post("/check-conflicts", authMiddleware.HasRole("admin", "teacher"), sh.CheckConflictsHandler)
	schedule.Put("/reschedule/:id", authMiddleware.HasRole("admin", "teacher"), sh.RescheduleHandler)
	schedule.Get("/week", authMiddleware.HasRole("admin", "teacher", "student"), sh.GetWeekSchedulesHandler)
	schedule.Get("/upcoming/:teacherId", authMiddleware.HasRole("admin", "teacher", "student"), sh.GetUpcomingSchedulesHandler)

	// Attendance routes
	attendance := api.Group("/attendance")
	attendance.Use(authMiddleware.AuthMiddleware())
	attendance.Post("/create", authMiddleware.HasRole("teacher"), ah.CreateAttendanceHandler)
	attendance.Post("/mark", authMiddleware.HasRole("teacher"), ah.MarkAttendanceHandler)

	attendance.Get("/:id", authMiddleware.HasRole("admin", "teacher", "student"), ah.GetAttendanceByIDHandler)
	attendance.Put("/update/:id", authMiddleware.HasRole("teacher"), ah.UpdateAttendanceHandler)
	attendance.Delete("/delete/:id", authMiddleware.HasRole("teacher"), ah.DeleteAttendanceHandler)
	attendance.Get("/student/:studentID", authMiddleware.HasRole("student", "teacher", "admin"), ah.GetAttendanceByStudentIDHandler)

	// Lesson routes
	lesson := api.Group("/lesson")
	lesson.Use(authMiddleware.AuthMiddleware())
	lesson.Post("/create", authMiddleware.HasRole("admin", "teacher"), lh.CreateLessonHandler)
	lesson.Get("/all", authMiddleware.HasRole("admin", "teacher", "student"), lh.GetAllLessonsHandler)
	lesson.Get("/:id", authMiddleware.HasRole("admin", "teacher", "student"), lh.GetLessonByIDHandler)
	lesson.Put("/update/:id", authMiddleware.HasRole("admin", "teacher"), lh.UpdateLessonHandler)
	lesson.Delete("/delete/:id", authMiddleware.HasRole("admin", "teacher"), lh.DeleteLessonHandler)

	// Homework routes
	homework := api.Group("/homework")
	homework.Use(authMiddleware.AuthMiddleware())
	homework.Post("/create", authMiddleware.HasRole("teacher"), hwh.CreateHomeworkHandler)
	homework.Get("/all", authMiddleware.HasRole("admin", "teacher", "student"), hwh.GetAllHomeworksHandler)
	homework.Get("/:id", authMiddleware.HasRole("admin", "teacher", "student"), hwh.GetHomeworkByIDHandler)
	homework.Put("/update/:id", authMiddleware.HasRole("teacher"), hwh.UpdateHomeworkHandler)
	homework.Delete("/delete/:id", authMiddleware.HasRole("teacher"), hwh.DeleteHomeworkHandler)
	homework.Get("/class/:classID", authMiddleware.HasRole("admin", "teacher", "student"), hwh.GetHomeworksByClassIDHandler)
	homework.Get("/teacher/:teacherID", authMiddleware.HasRole("admin", "teacher"), hwh.GetHomeworksByTeacherIDHandler)
	homework.Get("/lesson/:lessonID", authMiddleware.HasRole("admin", "teacher", "student"), hwh.GetHomeworksByLessonIDHandler)
	homework.Get("/active", authMiddleware.HasRole("admin", "teacher", "student"), hwh.GetActiveHomeworksHandler)
	homework.Get("/overdue", authMiddleware.HasRole("admin", "teacher", "student"), hwh.GetOverdueHomeworksHandler)
	homework.Get("/due-soon", authMiddleware.HasRole("admin", "teacher", "student"), hwh.GetHomeworksDueSoonHandler)
	homework.Put("/extend/:id", authMiddleware.HasRole("teacher"), hwh.ExtendDueDateHandler)
}
