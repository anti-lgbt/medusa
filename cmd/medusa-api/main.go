package main

import (
	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/routes"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	config.InitializeConfig()

	config.Database.AutoMigrate(
		&Product{},
	)

	r := routes.SetupRouter()

	r.Listen(":3000")
}
