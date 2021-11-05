package main

import (
	"log"

	"github.com/anti-lgbt/medusa/config"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	config.ConnectDatabase()

	log.Println(config.Database)

	config.Database.AutoMigrate(
		&Product{},
	)

	// r := routes.SetupRouter()

	// r.Listen(":3000")
}
