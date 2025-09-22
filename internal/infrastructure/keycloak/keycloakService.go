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

var (
	ADMIN_USERNAME       = "admin"
	ADMIN_PASSWORD       = "admin"
	KEYCLOAK_ADMIN_REALM = "master"
)

func NewKeycloakAuthService(hostname, clientId, clientSecret, realm string) (models.KeycloakService, models.ClassService) {
	return &KeycloakAuthService{
		Gocloak:      gocloak.NewClient(hostname),
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Realm:        realm,
		Hostname:     hostname,
	}, nil
}

func (kc *KeycloakAuthService) Login(login models.Login) (*models.LoginResponse, error) {
	ctx := context.Background()
	token, err := kc.Gocloak.Login(ctx, kc.ClientId, kc.ClientSecret, kc.Realm, login.Username, login.Password)
	if err != nil {
		return nil, fmt.Errorf("login fail: %w", err)
	}

	// Get admin token to retrieve user roles
	adminToken, err := kc.Gocloak.LoginAdmin(ctx, ADMIN_USERNAME, ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)
	if err != nil {
		return nil, fmt.Errorf("admin login fail to get user roles: %w", err)
	}

	// Get user ID from the obtained token (assuming username is unique and can be used to get user ID)
	users, err := kc.Gocloak.GetUsers(ctx, adminToken.AccessToken, kc.Realm, gocloak.GetUsersParams{Username: &login.Username})
	if err != nil || len(users) == 0 || users[0].ID == nil {
		return nil, fmt.Errorf("failed to get user details for role extraction: %w", err)
	}
	userID := *users[0].ID

	// Get user details by ID to extract role attribute
	userDetails, err := kc.Gocloak.GetUserByID(ctx, adminToken.AccessToken, kc.Realm, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID for role extraction: %w", err)
	}

	userRole := ""
	if userDetails.Attributes != nil && (*userDetails.Attributes)["role"] != nil && len((*userDetails.Attributes)["role"]) > 0 {
		userRole = (*userDetails.Attributes)["role"][0]
	}

	resp := &models.LoginResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Role:         userRole,
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

func (kc *KeycloakAuthService) CreateClass(class *models.Class) (string, error) {
	ctx := context.Background()

	adminToken, err := kc.Gocloak.LoginAdmin(ctx, ADMIN_USERNAME, ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)

	if err != nil {
		return "", fmt.Errorf("admin login fail: %w", err)
	}

	group := gocloak.Group{
		Name: gocloak.StringP(class.ClassName),
	}

	groupID, err := kc.Gocloak.CreateGroup(ctx, adminToken.AccessToken, kc.Realm, group)
	if err != nil {
		return "", fmt.Errorf("create group fail: %w", err)
	}

	return groupID, nil
}

func (kc *KeycloakAuthService) DeleteClass(classID string) error {
	ctx := context.Background()
	adminToken, err := kc.Gocloak.LoginAdmin(ctx, ADMIN_USERNAME, ADMIN_PASSWORD, kc.Realm)
	if err != nil {
		return fmt.Errorf("admin login fail:%w", err)
	}

	err = kc.Gocloak.DeleteGroup(ctx, adminToken.AccessToken, kc.Realm, classID)
	if err != nil {
		return fmt.Errorf("delete group fail:%w", err)
	}

	return nil
}

func (kc *KeycloakAuthService) GetAllClasses() ([]models.Class, error) {
	ctx := context.Background()
	adminToken, err := kc.Gocloak.LoginAdmin(ctx, ADMIN_USERNAME, ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)
	if err != nil {
		return nil, fmt.Errorf("admin login fail:%w", err)
	}

	groupParams := gocloak.GetGroupsParams{}

	groups, err := kc.Gocloak.GetGroups(ctx, adminToken.AccessToken, kc.Realm, groupParams)
	if err != nil {
		return nil, fmt.Errorf("get groups fail:%w", err)
	}

	var classes []models.Class

	for _, group := range groups {

		if group == nil || group.ID == nil || group.Name == nil {
			continue
		}

		teacherID := ""

		users, err := kc.Gocloak.GetGroupMembers(ctx, adminToken.AccessToken, kc.Realm, *group.ID, groupParams)

		if err != nil {
			return nil, fmt.Errorf("get users fail:%w", err)

		}

		for _, user := range users {
			if user.ID == nil {
				return nil, fmt.Errorf("user id null:%w", err)
			}

			userRoles, err := kc.Gocloak.GetRealmRolesByUserID(ctx, adminToken.AccessToken, kc.Realm, *user.ID)
			if err != nil {
				return nil, fmt.Errorf("get roles fail:%w", err)
			}

			for _, userRole := range userRoles {
				if *userRole.Name == "Teacher" {
					teacherID = *user.ID
					break
				}
			}
			if teacherID != "" {
				break
			}
		}

		class := models.Class{
			ID:        *group.ID,
			ClassName: *group.Name,
			TeacherID: teacherID,
		}

		classes = append(classes, class)

	}

	return classes, nil
}

func (kc *KeycloakAuthService) UpdateClass(class *models.Class) error {
	ctx := context.Background()

	adminToken, err := kc.Gocloak.LoginAdmin(ctx, ADMIN_USERNAME, ADMIN_PASSWORD, kc.Realm)
	if err != nil {
		return fmt.Errorf("admin login failed: %w", err)
	}

	updatedGroup := gocloak.Group{
		Name: gocloak.StringP("class_name"),
	}

	err = kc.Gocloak.UpdateGroup(ctx, adminToken.AccessToken, kc.Realm, updatedGroup)
	if err != nil {
		return fmt.Errorf("updated group fail:%w", err)

	}

	return nil
}

func (kc *KeycloakAuthService) GetClassesByTeacherID(teacherID string) ([]models.Class, error) {
	ctx := context.Background()
	adminToken, err := kc.Gocloak.LoginAdmin(ctx, ADMIN_USERNAME, ADMIN_PASSWORD, kc.Realm)
	if err != nil {
		return nil, fmt.Errorf("admin token fail:%w", err)
	}

	groupParams := gocloak.GetGroupsParams{}

	groupById, err := kc.Gocloak.GetUserGroups(ctx, adminToken.AccessToken, kc.Realm, teacherID, groupParams)
	if err != nil {
		return nil, fmt.Errorf("get groups by id fail:%w", err)

	}

	var classes []models.Class
	for _, kg := range groupById {
		class := models.Class{
			ID:        *kg.ID,   // Keycloak grubunun ID'si sınıf ID'si olarak kullanıldı
			ClassName: *kg.Name, // Keycloak grubunun adı sınıf adı olarak kullanıldı
		}
		classes = append(classes, class)
	}

	return classes, nil
}

func (kc *KeycloakAuthService) GetStudentsByClassID(classID string) ([]models.User, error) {
	ctx := context.Background()

	adminToken, err := kc.Gocloak.LoginAdmin(ctx, ADMIN_USERNAME, ADMIN_PASSWORD, KEYCLOAK_ADMIN_REALM)
	if err != nil {
		return nil, fmt.Errorf("admin login fail: %w", err)
	}

	groupMembers, err := kc.Gocloak.GetGroupMembers(ctx, adminToken.AccessToken, kc.Realm, classID, gocloak.GetGroupsParams{})
	if err != nil {
		return nil, fmt.Errorf("get group members fail: %w", err)
	}

	var students []models.User
	for _, member := range groupMembers {
		// Assuming a user with the role 'student' is a student
		if member.Attributes != nil && len((*member.Attributes)["role"]) > 0 && (*member.Attributes)["role"][0] == "student" {
			students = append(students, models.User{
				ID:          *member.ID,
				Username:    *member.Username,
				Email:       *member.Email,
				FirstName:   *member.FirstName,
				LastName:    *member.LastName,
				Phone:       (*member.Attributes)["phone"][0],
				Role:        (*member.Attributes)["role"][0],
				FamilyPhone: (*member.Attributes)["family_phone"][0],
			})
		}
	}
	return students, nil
}
