package keycloak

import (
	"Education_Dashboard/internal/models"
)


type KeycloakClassService struct {
	Kcauth models.KeycloakService
}

func NewKeycloakClassService(kcauth models.KeycloakService) models.ClassService {
	return &KeycloakClassService{
		Kcauth: kcauth,
	}
}

func (kc *KeycloakClassService) CreateClass(class *models.Class) error {
	
	return nil
}

func (kc *KeycloakClassService) GetClassByID(id string) (*models.Class, error) {
	
	return nil, nil
}

func (kc *KeycloakClassService) UpdateClass(class *models.Class) error {
	
	return nil
}

func (kc *KeycloakClassService) DeleteClass(id string) error {

	return nil
}

func (kc *KeycloakClassService) GetAllClasses() ([]models.Class, error) {
	return nil, nil
}


func (kc *KeycloakClassService) GetClassesByTeacherID(teacherID string) ([]models.Class, error) {
	return nil, nil
}