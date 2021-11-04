package config

import "gorm.io/gorm"

var Database *gorm.DB

func InitializeConfig() error {
	NewLoggerService()

	if err := ConnectDatabase(); err != nil {
		return err
	}
	if err := InitSessionStore(); err != nil {
		return err
	}

	return nil
}
