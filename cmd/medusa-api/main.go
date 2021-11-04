package main

import (
	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/routes"
)

func main() {
	if err := config.InitializeConfig(); err != nil {
		config.Logger.Error(err.Error())
		return
	}

	config.Database.AutoMigrate(
		&models.User{},
	)

	r := routes.SetupRouter()

	r.Listen(":3000")
}
