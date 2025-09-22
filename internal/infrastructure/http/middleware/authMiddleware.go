package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"Education_Dashboard/internal/models"
)

type AuthMiddleware struct {
	keycloakService models.KeycloakService
}

func NewAuthMiddleware(kc models.KeycloakService) AuthMiddleware {
	return AuthMiddleware{
		keycloakService: kc,
	}
}

func (am AuthMiddleware) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Missing or malformed JWT",
			})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Missing or malformed JWT",
			})
		}

		// Validate the token with Keycloak using introspection
		res, err := am.keycloakService.RetrospectToken(tokenString)
		if err != nil || !res {
			fmt.Println("Token introspection failed or token is inactive:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Invalid or inactive token",
			})
		}

		claimsMap, err := am.keycloakService.DecodeToken(tokenString)
		if err != nil {
			fmt.Println("Error decoding token:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Invalid token claims",
			})
		}

		claims := claimsMap

		// Extract user roles from Keycloak claims
		var roles []string
		if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
			if realmRoles, ok := realmAccess["roles"].([]interface{}); ok {
				for _, role := range realmRoles {
					if r, isString := role.(string); isString {
						roles = append(roles, r)
					}
				}
			}
		}

		userID, ok := claims["sub"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "User ID not found in token",
			})
		}

		c.Locals("userID", userID)
		c.Locals("userRoles", roles)
		c.Locals("isAuthenticated", true)

		return c.Next()
	}
}

func (am AuthMiddleware) HasRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoles, ok := c.Locals("userRoles").([]string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Forbidden",
				"message": "User roles not found in context",
			})
		}

		for _, requiredRole := range roles {
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					return c.Next()
				}
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   "Forbidden",
			"message": "User does not have the required role",
		})
	}
}

func (am AuthMiddleware) GetUserID(c *fiber.Ctx) (string, error) {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}
