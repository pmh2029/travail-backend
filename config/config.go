package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"travail/internal/pkg/migrations"
	"travail/pkg/shared/database"
	sharedLogger "travail/pkg/shared/logger"
)

var (
	DBConn *gorm.DB
)

func LoadConfig() {
	logger := sharedLogger.NewLogger()
	LoadEnv(logger)
	LoadDB(logger)
}

func LoadEnv(logger *logrus.Logger) {
	err := godotenv.Load(filepath.Join(".env"))
	if err != nil {
		logger.Fatalln("Fail to load .env")
	}
}

func LoadDB(logger *logrus.Logger) {
	dbConfig := database.DBConfig{
		Host:    os.Getenv("DB_HOST"),
		Name:    os.Getenv("DB_NAME"),
		User:    os.Getenv("DB_USER"),
		Pass:    os.Getenv("DB_PASS"),
		Port:    os.Getenv("DB_PORT"),
		Type:    database.PostgreSQL,
		Charset: "utf8mb4",
	}

	logger.Info("Init Database")
	DBConn, err := database.NewDB(dbConfig, logger)
	if err != nil {
		logger.Fatalln("Fail to connect to database")
		panic(err)
	}
	logger.Info("Init Database Success")

	defer database.CloseDB(DBConn, logger)

	logger.Info("Migrate Database Start")
	err = migrations.Migrate(DBConn)
	if err != nil {
		logger.Fatalln("Fail to migrate database")
		panic(err)
	}
	logger.Info("Migrate Database Success")
}
