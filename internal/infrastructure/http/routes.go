package http

import (
	"github.com/gofiber/fiber/v2"
	"Education_Dashboard/internal/infrastructure/http/handler"
	
	
)

func AuthRoutes(app *fiber.App, kh handler.KeycloakHandler) {
	api := app.Group("/v1/api")

	api.Post("/create", kh.RegisterHandler)
	api.Post("/login", kh.LoginHandler)
	api.Post("/logout", kh.LogoutHandler)

	user := api.Group("/user")
	user.Put("/update/:id", kh.UpdateUserHandler)
	user.Delete("/delete/:id", kh.DeleteUserHandler)
	user.Get("/currentuser/:id", kh.GetUserByIDHandler)
	user.Get("/allusers", kh.GetAllUsersHandler)
}
