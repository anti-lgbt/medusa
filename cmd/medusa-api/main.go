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
		&models.Activity{},
		&models.Album{},
		&models.Music{},
		&models.Code{},
		&models.Comment{},
		&models.Reply{},
		&models.MusicAlbum{},
		&models.Like{},
		&models.Label{},
		&models.Popular{},
		&models.User{},
	)

	r := routes.SetupRouter()

	r.Listen(":3000")
}
