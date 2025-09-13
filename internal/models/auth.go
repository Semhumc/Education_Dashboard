package models

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Register struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email" binding:"required"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Role        string `json:"role"`
	FamilyPhone string `json:"family_phone"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Role        string `json:"role"`
	FamilyPhone string `json:"family_phone"`
}

type KeycloakService interface{
	Register(register Register) error
	Login(login Login) (*LoginResponse, error)
	UpdateUser(id string, register Register) error
	DeleteUser(id string) error
	GetUserByID(id string) (User, error)
	GetAllUsers() ([]User, error)
	Logout(loginResponse LoginResponse) error
	//RefreshToken(refreshToken string) (*LoginResponse, error)
}