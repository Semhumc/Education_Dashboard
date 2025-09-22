package handler

import (
	"Education_Dashboard/internal/models"

	"github.com/gofiber/fiber/v2"
)

type KeycloakHandler struct {
	keycloakService models.KeycloakService
}

func NewKeycloakHandler(kc models.KeycloakService) KeycloakHandler {
	return KeycloakHandler{
		keycloakService: kc,
	}
}

func (kch KeycloakHandler) LoginHandler(c *fiber.Ctx) error {

	var loginparams models.Login
	if err := c.BodyParser(&loginparams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	if loginparams.Username == "" || loginparams.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "username and password are required",
		})
	}
	
	resp, err := kch.keycloakService.Login(loginparams)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (kch KeycloakHandler) RegisterHandler(c *fiber.Ctx) error {

	var registerparams models.Register
	if err := c.BodyParser(&registerparams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{	
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}
	if registerparams.Username == "" || registerparams.Password == "" || registerparams.Email == "" || registerparams.FirstName == "" 	|| registerparams.LastName == "" || registerparams.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "username, password, email, first name, last name and phone are required",
		})
	}
	err := kch.keycloakService.Register(registerparams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{	
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

func (kch KeycloakHandler) UpdateUserHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "user id is required",
		})
	}
	var registerparams models.Register
	if err := c.BodyParser(&registerparams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}
	if registerparams.Username == "" || registerparams.Password == "" || registerparams.Email == "" || registerparams.FirstName == "" || registerparams.LastName == "" || registerparams.Phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "username, password, email, first name, last name and phone are required",
		})
	}
	err := kch.keycloakService.UpdateUser(userID, registerparams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

func (kch KeycloakHandler) DeleteUserHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {	
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "user id is required",
		})
	}	
	err := kch.keycloakService.DeleteUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func (kch KeycloakHandler) GetUserByIDHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "user id is required",
		})
	}
	user, err := kch.keycloakService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(user)
}

func (kch KeycloakHandler) GetAllUsersHandler(c *fiber.Ctx) error {
	users, err := kch.keycloakService.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

func (kch KeycloakHandler) LogoutHandler(c *fiber.Ctx) error {
	var loginResponse models.LoginResponse
	if err := c.BodyParser(&loginResponse); err != nil {	
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}	
		
	if loginResponse.AccessToken == "" || loginResponse.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "access token and refresh token are required",
		})
	}

	err := kch.keycloakService.Logout(loginResponse)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (kch KeycloakHandler) CreateUserHandler(c *fiber.Ctx) error {
	var registerparams models.Register
	if err := c.BodyParser(&registerparams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": err.Error(),
		})
	}

	if registerparams.Username == "" || registerparams.Password == "" || registerparams.Email == "" || registerparams.FirstName == "" || registerparams.LastName == "" || registerparams.Phone == "" || registerparams.Role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "username, password, email, first name, last name, phone, and role are required",
		})
	}

	err := kch.keycloakService.Register(registerparams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

