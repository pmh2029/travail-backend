package migrations

import (
	"gorm.io/gorm"
)

func Migrate(dbConn *gorm.DB) error {
	err := dbConn.AutoMigrate()

	return err
}
