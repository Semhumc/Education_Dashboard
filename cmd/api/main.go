package main

import (
	"Education_Dashboard/internal/infrastructure/http"
	"Education_Dashboard/internal/infrastructure/http/handler"
	"Education_Dashboard/internal/infrastructure/keycloak"
	"os"

	"github.com/gofiber/fiber/v2"
)

var (
	keycloak_base_url      = os.Getenv("KEYCLOAK_BASE_URL")
	keycloak_realm         = os.Getenv("KEYCLOAK_REALM")
	keycloak_client_id     = os.Getenv("KEYCLOAK_CLIENT_ID")
	keycloak_client_secret = os.Getenv("KEYCLOAK_CLIENT_SECRET")
	port                   = "5000"
)

func main() {

	app := fiber.New()

	keycloakAuthService := keycloak.NewKeycloakAuthService(
		keycloak_client_id,
		keycloak_client_secret,
		keycloak_realm,
		keycloak_base_url,
	)

	authHandler := handler.NewKeycloakHandler(keycloakAuthService)

	http.AuthRoutes(app, authHandler)
	app.Listen(":" + port)
}
