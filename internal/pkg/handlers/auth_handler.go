package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"travail/config"
	"travail/internal/pkg/domains/interfaces"
	"travail/internal/pkg/domains/models/dtos/req"
	"travail/internal/pkg/domains/models/dtos/res"
	"travail/internal/pkg/repositories"
	"travail/internal/pkg/usecases"
	"travail/pkg/shared/auth"
	"travail/pkg/shared/constants"
	"travail/pkg/shared/utils"
)

type AuthHandler struct {
	AuthUsecase interfaces.AuthUsecase
}

func NewUserHandler(dbConn *gorm.DB) *AuthHandler {
	authRepo := repositories.NewAuthRepository(dbConn)
	authUsecase := usecases.NewAuthUsecase(authRepo)
	return &AuthHandler{
		AuthUsecase: authUsecase,
	}
}

func (authHandler *AuthHandler) SignUp(c *gin.Context) {
	req := req.UserSignUpRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status: "failed",
			Error: &res.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	user, err := authHandler.AuthUsecase.SignUp(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.BaseResponse{
			Status: "failed",
			Error: &res.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, res.BaseResponse{
		Status: "success",
		Data:   gin.H{"user_info": utils.ConvertUserEntityToUserResponse(user)},
	})
}

func (authHandler *AuthHandler) SignIn(c *gin.Context) {
	req := req.UserSignInRequest{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, res.BaseResponse{
			Status: "failed",
			Error: &res.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	user, token, err := authHandler.AuthUsecase.SignIn(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.BaseResponse{
			Status: "failed",
			Error: &res.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, res.BaseResponse{
		Status: "success",
		Data:   gin.H{"access_token": token, "user_info": utils.ConvertUserEntityToUserResponse(user)},
	})
}

func (authHandler *AuthHandler) SignInWithGoogle(c *gin.Context) {
	if c.Request.Method != "GET" {
		c.JSON(http.StatusMethodNotAllowed, res.BaseResponse{
			Status: "failed",
			Error: &res.ErrorResponse{
				ErrorMessage: "method not allowed",
			},
		})
		return
	}

	// Create oauthState cookie
	oauthState := utils.GenerateStateOauthCookie(c)
	/*
		AuthCodeURL receive state that is a token to protect the user
		from CSRF attacks. You must always provide a non-empty string
		and validate that it matches the the state query parameter
		on your redirect callback.
	*/
	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (authHandler *AuthHandler) Redirect(c *gin.Context) {
	// check is method is correct
	if c.Request.Method != "GET" {
		c.JSON(http.StatusMethodNotAllowed, res.BaseResponse{
			Status: "failed",
			Error: &res.ErrorResponse{
				ErrorMessage: "method not allowed",
			},
		})
		return
	}

	// get oauth state from cookie for this user
	oauthState, _ := c.Request.Cookie("oauthstate")
	state := c.Request.FormValue("state")
	code := c.Request.FormValue("code")
	c.Header("content-type", "application/json")

	// ERROR : Invalid OAuth State
	if state != oauthState.Value {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		c.JSON(http.StatusInternalServerError, res.BaseResponse{
			Status: "failed",
			Error: &res.ErrorResponse{
				ErrorMessage: "invalid oauth state",
			},
		})
		return
	}

	// Exchange Auth Code for Tokens
	token, err := config.AppConfig.GoogleLoginConfig.Exchange(context.Background(), code)

	// ERROR : Auth Code Exchange Failed
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.BaseResponse{
			Status: "failed",
			Error: &res.ErrorResponse{
				ErrorMessage: "falied code exchange: " + err.Error(),
			},
		})
		return
	}

	// Fetch User Data from google server
	response, err := http.Get(constants.OauthGoogleUrlAPI + token.AccessToken)
	// ERROR : Unable to get user data from google
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.BaseResponse{
			Status: "failed",
			Error: &res.ErrorResponse{
				ErrorMessage: "failed getting user info: " + err.Error(),
			},
		})
		return
	}

	// Parse user data JSON Object
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.BaseResponse{
			Status: "failed",
			Error: &res.ErrorResponse{
				ErrorMessage: "failed read response: " + err.Error(),
			},
		})
		return
	}

	err = json.Unmarshal(contents, &config.GoogleUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, res.BaseResponse{
			Status: "failed",
			Error: &res.ErrorResponse{
				ErrorMessage: "failed unmarshal response: " + err.Error(),
			},
		})
		return
	}

	user, err := authHandler.AuthUsecase.TakeByConditions(map[string]interface{}{
		"email": config.GoogleUser.Email,
	})

	if err != nil && err == gorm.ErrRecordNotFound {
		req := req.UserSignUpRequest{
			Username: config.GoogleUser.Name + "@" + strings.Split(config.GoogleUser.Email, "@")[0],
			Email:    config.GoogleUser.Email,
			Password: config.GoogleUser.ID,
		}

		user, err = authHandler.AuthUsecase.SignUp(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, res.BaseResponse{
				Status: "failed",
				Error: &res.ErrorResponse{
					ErrorMessage: err.Error(),
				},
			})
			return
		}

		jwtToken, _ := auth.GenerateHS256JWT(map[string]interface{}{
			"email":    config.GoogleUser.Email,
			"username": config.GoogleUser.Name,
		}, time.Now().Add(time.Hour*72))

		c.JSON(http.StatusOK, res.BaseResponse{
			Status: "success",
			Data:   gin.H{"access_token": jwtToken, "user_info": utils.ConvertUserEntityToUserResponse(user)},
		})
	} else if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, res.BaseResponse{
			Status: "failed",
			Error: &res.ErrorResponse{
				ErrorMessage: err.Error(),
			},
		})
	} else {
		jwtToken, _ := auth.GenerateHS256JWT(map[string]interface{}{
			"email":    config.GoogleUser.Email,
			"username": config.GoogleUser.Name,
		}, time.Now().Add(time.Hour*72))

		c.JSON(http.StatusOK, res.BaseResponse{
			Status: "success",
			Data:   gin.H{"access_token": jwtToken, "user_info": utils.ConvertUserEntityToUserResponse(user)},
		})
	}
}
