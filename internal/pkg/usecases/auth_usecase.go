package usecases

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"os"
	"text/template"
	"time"

	"gopkg.in/gomail.v2"

	"travail/internal/pkg/domains/interfaces"
	"travail/internal/pkg/domains/models/dtos/req"
	"travail/internal/pkg/domains/models/entities"
	"travail/pkg/shared/auth"
	"travail/pkg/shared/utils"
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
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	return user, token, err
}

func (authUsecase *AuthUsecase) SendMailForgotPassword(req req.ForgotPasswordRequest) error {
	var body bytes.Buffer
	t, err := template.ParseFiles("template/forgot_password_template.html")
	if err != nil {
		return err
	}

	user, err := authUsecase.AuthRepo.TakeByConditions(map[string]interface{}{
		"email": req.Email,
	})
	if err != nil {
		return err
	}

	encodedUsername := base64.StdEncoding.EncodeToString([]byte(user.Username))

	token, err := auth.GenerateHS256JWT(map[string]interface{}{
		"email": user.Email,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
	})
	if err != nil {
		return err
	}

	resetPasswordUrl := "http://localhost:3000/api/auth/reset_password/" + encodedUsername + "/" + token

	t.Execute(&body, struct {
		Name string
		Url  string
	}{Name: user.Username, Url: resetPasswordUrl})

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", os.Getenv("TRAVAIL_EMAIL"))

	// Set E-Mail receivers
	m.SetHeader("To", req.Email)

	// Set E-Mail subject
	m.SetHeader("Subject", "Forgot Password")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", body.String())

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("TRAVAIL_EMAIL"), os.Getenv("TRAVAIL_APP_PASS"))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	err = d.DialAndSend(m)

	return err
}

func (authUsecase *AuthUsecase) ResetPassword(req req.ResetPasswordRequest) error {
	user, err := authUsecase.AuthRepo.TakeByConditions(map[string]interface{}{
		"email": req.Email,
	})
	if err != nil {
		return err
	}

	hashPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashPassword
	err = authUsecase.AuthRepo.ResetPassword(int(user.ID), user)

	return err
}
