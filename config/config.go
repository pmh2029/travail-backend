package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"travail/internal/pkg/migrations"
	"travail/pkg/shared/database"
)

func LoadConfig(logger *logrus.Logger) {
	LoadEnv(logger)
	LoadDB(logger)
}

func LoadEnv(logger *logrus.Logger) {
	err := godotenv.Load(filepath.Join(".env"))
	if err != nil {
		logger.Fatalln("Fail to load .env")
	}
}

func LoadDB(logger *logrus.Logger) *gorm.DB {
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
	dbConn, err := database.NewDB(dbConfig, logger)
	if err != nil {
		logger.Fatalln("Fail to connect to database")
		panic(err)
	}
	logger.Info("Init Database Success")

	logger.Info("Migrate Database Start")
	err = migrations.Migrate(dbConn)
	if err != nil {
		logger.Fatalln("Fail to migrate database")
		panic(err)
	}
	logger.Info("Migrate Database Success")
	
	return dbConn
}
