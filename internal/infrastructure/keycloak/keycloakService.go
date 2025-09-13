package keycloak

import (
	"Education_Dashboard/internal/models"

	"context"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
)

type KeycloakAuthService struct {
	Gocloak      *gocloak.GoCloak
	ClientId     string
	ClientSecret string
	Realm        string
	Hostname     string
}

const (
	ADMIN_USERNAME       = "admin"
	ADMIN_PASSWORD       = "admin"
	KEYCLOAK_ADMIN_REALM = "master"
)

func NewKeycloakAuthService(hostname, clientId, clientSecret, realm string) models.KeycloakService {
	return &KeycloakAuthService{
		Gocloak:      gocloak.NewClient(hostname),
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Realm:        realm,
		Hostname:     hostname,
	}
}

func (kc *KeycloakAuthService) Login(login models.Login) (*models.LoginResponse, error) {
	ctx := context.Background()
	token, err := kc.Gocloak.Login(ctx, kc.ClientId, kc.ClientSecret, kc.Realm, login.Username, login.Password)
	if err != nil {
		return nil, fmt.Errorf("login fail: %w", err)
	}

	resp := &models.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	return resp, nil
}

func (kc *KeycloakAuthService) Register(register models.Register) error {
	ctx := context.Background()

	adminToken, err := kc.Gocloak.LoginAdmin(ctx, ADMIN_USERNAME, ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)
	if err != nil {
		return fmt.Errorf("admin login fail: %w", err)
	}

	user := gocloak.User{
		Username:      gocloak.StringP(register.Username),
		Email:         gocloak.StringP(register.Email),
		FirstName:     gocloak.StringP(register.FirstName),
		LastName:      gocloak.StringP(register.LastName),
		Attributes:    &map[string][]string{"role": {register.Role}, "family_phone": {register.FamilyPhone}, "phone": {register.Phone}},
		EmailVerified: gocloak.BoolP(true),
		Enabled:       gocloak.BoolP(true),
	}

	userID, err := kc.Gocloak.CreateUser(ctx, adminToken.AccessToken, kc.Realm, user)
	if err != nil {
		return fmt.Errorf("create user fail: %w", err)
	}

	cred := gocloak.CredentialRepresentation{
		Type:      gocloak.StringP("password"),
		Value:     gocloak.StringP(register.Password),
		Temporary: gocloak.BoolP(false),
	}

	err = kc.Gocloak.SetPassword(ctx, adminToken.AccessToken, userID, KEYCLOAK_ADMIN_REALM, *cred.Value, false)
	if err != nil {
		return fmt.Errorf("set password fail: %w", err)
	}

	err = kc.Gocloak.SendVerifyEmail(ctx, adminToken.AccessToken, userID, kc.Realm, gocloak.SendVerificationMailParams{
		ClientID:    gocloak.StringP(kc.ClientId),
		RedirectURI: gocloak.StringP("http://localhost:3000/"),
	})
	if err != nil {
		return fmt.Errorf("send verify email failed: %w", err)
	}

	return nil
}

func (kc *KeycloakAuthService) UpdateUser(userID string, register models.Register) error {
	ctx := context.Background()

	adminToken, err := kc.Gocloak.LoginAdmin(ctx, ADMIN_USERNAME, ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)

	if err != nil {
		return fmt.Errorf("admin login fail: %w", err)
	}
	user := gocloak.User{
		Username:   gocloak.StringP(register.Username),
		Email:      gocloak.StringP(register.Email),
		FirstName:  gocloak.StringP(register.FirstName),
		LastName:   gocloak.StringP(register.LastName),
		Attributes: &map[string][]string{"role": {register.Role}, "family_phone": {register.FamilyPhone}, "phone": {register.Phone}},
	}

	err = kc.Gocloak.UpdateUser(ctx, adminToken.AccessToken, kc.Realm, user)

	if err != nil {
		return fmt.Errorf("update user fail: %w", err)
	}
	return nil
}

func (kc *KeycloakAuthService) DeleteUser(userID string) error {
	ctx := context.Background()
	adminToken, err := kc.Gocloak.LoginAdmin(ctx, ADMIN_USERNAME, ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)
	if err != nil {
		return fmt.Errorf("admin login failed: %w", err)
	}

	err = kc.Gocloak.DeleteUser(ctx, adminToken.AccessToken, kc.Realm, userID)
	if err != nil {
		return fmt.Errorf("delete user failed: %w", err)
	}
	return nil
}

func (kc *KeycloakAuthService) GetUserByID(userID string) (models.User, error) {
	ctx := context.Background()
	adminToken, err := kc.Gocloak.LoginAdmin(ctx, ADMIN_USERNAME, ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)
	if err != nil {
		return models.User{}, fmt.Errorf("admin login fail: %w", err)
	}

	user, err := kc.Gocloak.GetUserByID(ctx, adminToken.AccessToken, userID, kc.Realm)
	if err != nil {
		return models.User{}, fmt.Errorf("get user by id fail: %w", err)
	}

	User := models.User{
		ID:          *user.ID,
		Username:    *user.Username,
		Email:       *user.Email,	
		FirstName:   *user.FirstName,
		LastName:    *user.LastName,
		Phone:       (*user.Attributes)["phone"][0],
		Role:        (*user.Attributes)["role"][0],
		FamilyPhone: (*user.Attributes)["family_phone"][0],
	}

	return User, nil

}

func (kc *KeycloakAuthService) GetAllUsers() ([]models.User, error) {
	ctx := context.Background()
	adminToken, err := kc.Gocloak.LoginAdmin(ctx, ADMIN_USERNAME, ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)
	if err != nil {
		return nil, fmt.Errorf("admin login fail: %w", err)
	}
	users, err := kc.Gocloak.GetUsers(ctx, adminToken.AccessToken, kc.Realm, gocloak.GetUsersParams{})
	if err != nil {
		return nil, fmt.Errorf("get all users fail: %w", err)
	}

	var result []models.User
	for _, user := range users {
		result = append(result, models.User{
			ID:          *user.ID,
			Username:    *user.Username,
			Email:       *user.Email,
			FirstName:   *user.FirstName,
			LastName:    *user.LastName,
			Phone:       (*user.Attributes)["phone"][0],
			Role:        (*user.Attributes)["role"][0],
			FamilyPhone: (*user.Attributes)["family_phone"][0],
		})
	}

	return result, nil
}

func (kc *KeycloakAuthService) Logout(loginResponse models.LoginResponse) error {
	ctx := context.Background()
	err := kc.Gocloak.Logout(ctx, kc.ClientId, kc.ClientSecret, kc.Realm, loginResponse.RefreshToken)
	if err != nil {	
		return fmt.Errorf("logout fail: %w", err)
	}	
	return nil
}