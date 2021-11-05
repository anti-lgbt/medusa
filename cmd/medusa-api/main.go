package main

import (
	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/models"
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
		&models.Activity{},
		&models.Code{},
		&models.Like{},
		&models.Reply{},
		&models.Comment{},
		&models.MusicAlbum{},
		&models.Music{},
		&models.Album{},
		&models.Popular{},
		&models.Label{},
		&models.User{},
	)

	r := routes.SetupRouter()

	r.Listen(":3000")
}
