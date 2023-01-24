package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"travail/internal/pkg/domains/models/dtos/res"
	"travail/internal/pkg/handlers"
)

// Router is application struct
type Router struct {
	Engine *gin.Engine
	DBConn *gorm.DB
	Logger *logrus.Logger
}

// InitializeRouter initializes Engine and middleware
func (r *Router) InitializeRouter(logger *logrus.Logger) {
	r.Engine.Use(gin.Logger())
	r.Engine.Use(gin.Recovery())
	r.Engine.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}))
	r.Logger = logger
}

func (r *Router) SetupHandler() {
	authHandler := handlers.NewUserHandler(r.DBConn)

	// ping
	r.Engine.GET("/ping", func(c *gin.Context) {
		data := res.BaseResponse{
			Status: "success",
			Data:   gin.H{"message": "Pong!"},
			Error:  nil,
		}
		c.JSON(http.StatusOK, data)
	})

	// router api
	publicApi := r.Engine.Group("/api")
	{
		// auth
		authAPI := publicApi.Group("/auth")
		{
			authAPI.POST("/signup", authHandler.SignUp)
			authAPI.POST("/signin", authHandler.SignIn)
			authAPI.GET("/google/signin", authHandler.SignInWithGoogle)
			authAPI.GET("/google/redirect", authHandler.Redirect)
			authAPI.POST("/forgot_password", authHandler.ForgotPassword)
			authAPI.PATCH("/reset_password", authHandler.ResetPassword)
		}
	}
}
