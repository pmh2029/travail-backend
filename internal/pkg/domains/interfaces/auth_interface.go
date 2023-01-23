package interfaces

import (
	"travail/internal/pkg/domains/models/dtos/req"
	"travail/internal/pkg/domains/models/entities"
)

// AuthRepository interface
type AuthRepository interface {
	SignUp(user entities.User) (entities.User, error)
	TakeByConditions(conditions map[string]interface{}) (entities.User, error)
}

// AuthUsecase interface
type AuthUsecase interface {
	SignUp(req req.UserSignUpRequest) (entities.User, error)
	TakeByConditions(conditions map[string]interface{}) (entities.User, error)
	SignIn(req req.UserSignInRequest) (entities.User, string, error)
}
