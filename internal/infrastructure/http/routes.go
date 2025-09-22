package http

import (
	"Education_Dashboard/internal/infrastructure/http/handler"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, kh handler.KeycloakHandler) {
	api := app.Group("/v1/api") // Activate this line

	api.Post("/create", kh.RegisterHandler)
	api.Post("/login", kh.LoginHandler)
	api.Post("/logout", kh.LogoutHandler)

	user := api.Group("/user") // Change 'app.Group' to 'api.Group'
	user.Put("/update/:id", kh.UpdateUserHandler)
	user.Delete("/delete/:id", kh.DeleteUserHandler)
	user.Get("/currentuser/:id", kh.GetUserByIDHandler)
	user.Get("/allusers", kh.GetAllUsersHandler)
}
