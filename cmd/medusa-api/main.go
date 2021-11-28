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
		&models.User{},
		&models.Activity{},
		&models.Label{},
		&models.Popular{},
		&models.Album{},
		&models.Music{},
		&models.MusicAlbum{},
		&models.Comment{},
		&models.Reply{},
		&models.Like{},
		&models.Code{},
		&models.TrendingMusic{},
	)

	r := routes.SetupRouter()

	r.Listen(":3000")
}
