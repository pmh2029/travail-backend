package config

import "gorm.io/gorm"

var (
	DBConn     *gorm.DB
	AppConfig  OAuthConfig
	GoogleUser GoogleUserInfo
)
