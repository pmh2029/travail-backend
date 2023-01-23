package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	sharedLogger "travail/pkg/shared/logger"
)

// DBConfig config for DB
type DBConfig struct {
	Host    string
	Name    string
	User    string
	Pass    string
	Port    string
	Charset string
	Type    DBType
}

// NewDB initialize database
func NewDB(config DBConfig, logger *logrus.Logger) (*gorm.DB, error) {
	var dbConn *gorm.DB
	var err error
	switch config.Type {
	case MySQL:
		// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", config.User, config.Pass, config.Host, config.Port, config.Name, config.Charset)

		dbConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: sharedLogger.NewGormLogger(logger),
		})
		if err != nil {
			return nil, err
		}
	case PostgreSQL:
		// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable"
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.Host, config.User, config.Pass, config.Name, config.Port)

		dbConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: sharedLogger.NewGormLogger(logger),
		})
		if err != nil {
			return nil, err
		}
	default:
		panic("Unknown type of database!")
	}
	err = Ping(dbConn)
	return dbConn, err
}

func CloseDB(db *gorm.DB, logger *logrus.Logger) {
	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorf("Error while returning *sql.DB: %v", err)
	}

	sqlDB.Close()
	logger.Info("Closing the DB connection pool")
}

func Ping(db *gorm.DB) error {
	myDB, err := db.DB()
	if err != nil {
		return err
	}

	return myDB.Ping()
}
