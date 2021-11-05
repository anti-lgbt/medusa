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
