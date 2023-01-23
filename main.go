package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"travail/config"
	"travail/internal/app/myapi/router"
	sharedLogger "travail/pkg/shared/logger"
)

func main() {
	logger := sharedLogger.NewLogger()
	gin.SetMode(gin.DebugMode)

	config.LoadConfig()

	engine := gin.New()
	router := &router.Router{
		Engine: engine,
		DBConn: config.DBConn,
	}

	router.InitializeRouter(logger)
	router.SetupHandler()

	engine.Run(":" + os.Getenv("API_PORT"))
}
