package repositories

import (
	"travail/internal/pkg/domains/interfaces"
	"travail/internal/pkg/domains/models/entities"

	"gorm.io/gorm"
)

type AuthRepository struct {
	DBConn *gorm.DB
}

func NewAuthRepository(dbConn *gorm.DB) interfaces.AuthRepository {
	return &AuthRepository{
		DBConn: dbConn,
	}
}

func (authRepo *AuthRepository) SignUp(user entities.User) (entities.User, error) {
	result := authRepo.DBConn.Create(&user)

	return user, result.Error
}

func (authRepo *AuthRepository) TakeByConditions(conditions map[string]interface{}) (entities.User, error) {
	user := entities.User{}
	result := authRepo.DBConn.Where(conditions).Take(&user)

	return user, result.Error
}
