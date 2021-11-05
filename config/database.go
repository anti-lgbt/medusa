package config

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectDatabase() {
	dsn := "host=" + os.Getenv("DATABASE_HOST") +
		" port=" + os.Getenv("DATABASE_PORT") +
		" user=" + os.Getenv("DATABASE_USER") +
		" password=" + os.Getenv("DATABASE_PASS") +
		" dbname=" + os.Getenv("DATABASE_NAME") +
		" sslmode=disable"

	var err error

	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
}
