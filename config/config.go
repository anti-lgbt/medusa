package config

import "gorm.io/gorm"

var Database *gorm.DB

func InitializeConfig() {
	NewLoggerService()
	ConnectDatabase()
	InitSessionStore()
}
