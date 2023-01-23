package usecases

import (
	"errors"
	"travail/internal/pkg/domains/interfaces"
	"travail/internal/pkg/domains/models/dtos/req"
	"travail/internal/pkg/domains/models/entities"
	"travail/pkg/shared/auth"
	"travail/pkg/shared/utils"
	"time"
)

type AuthUsecase struct {
	AuthRepo interfaces.AuthRepository
}

func NewAuthUsecase(authRepo interfaces.AuthRepository) interfaces.AuthUsecase {
	return &AuthUsecase{
		AuthRepo: authRepo,
	}
}

func (authUsecase *AuthUsecase) SignUp(req req.UserSignUpRequest) (entities.User, error) {
	user := entities.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return entities.User{}, err
	}

	user.Password = hashPassword
	user, err = authUsecase.AuthRepo.SignUp(user)
	return user, err
}

func (authUsecase *AuthUsecase) TakeByConditions(conditions map[string]interface{}) (entities.User, error) {
	user, err := authUsecase.AuthRepo.TakeByConditions(conditions)

	return user, err
}

func (authUsecase *AuthUsecase) SignIn(req req.UserSignInRequest) (user entities.User, token string, err error) {
	conditions := map[string]interface{}{}
	if req.Username == "" && req.Email == "" {
		return entities.User{}, "", errors.New("username or email is required")
	}

	if req.Username != "" {
		conditions["username"] = req.Username
	}

	if req.Email != "" {
		conditions["email"] = req.Email
	}

	user, err = authUsecase.AuthRepo.TakeByConditions(conditions)
	if err != nil {
		return entities.User{}, "", err
	}

	result := utils.CheckHashPassword(req.Password, user.Password)
	if !result {
		return entities.User{}, "", errors.New("password is incorrect")
	}

	token, err = auth.GenerateHS256JWT(map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
	}, time.Now().Add(time.Hour*72))

	return user, token, err
}
