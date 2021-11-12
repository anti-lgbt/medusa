package public

import (
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/types"
	"github.com/gofiber/fiber/v2"
)

// GET /api/v2/public/user/avatar/:id
func GetUserAvatar(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.ServerInvalidQuery)
	}

	var user *models.User
	if result := config.Database.First(&user, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	return c.Status(200).SendFile(user.Avatar.String, false)
}

// GET /api/v2/public/musics/:id
func GetMusic(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.ServerInvalidQuery)
	}

	var music *models.Music
	if result := config.Database.First(&music, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	return c.Status(200).JSON(music.ToEntity())
}

// GET /api/v2/public/musics/:id/audio
func GetMusicAudio(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.ServerInvalidQuery)
	}

	var music *models.Music
	if result := config.Database.First(&music, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	return c.Status(200).SendFile(music.Path, false)
}

// GET /api/v2/public/musics/:id/image
func GetMusicImage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.ServerInvalidQuery)
	}

	var music *models.Music
	if result := config.Database.First(&music, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	return c.Status(200).SendFile(music.Image.String, false)
}

// GET /api/v2/public/albums/:id
func GetAlbum(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.ServerInvalidQuery)
	}

	var album *models.Album
	if result := config.Database.First(&album, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	if album.Private {
		user := c.Locals("CurrentUser").(*models.User)

		if user == nil {
			return c.Status(404).JSON(types.Error{
				Error: types.RecordNotFound,
			})
		}

		if user.ID == album.UserID {
			return c.Status(200).JSON(album.ToEntity())
		}
	}

	return c.Status(200).JSON(album.ToEntity())
}

// GET /api/v2/public/albums/:id/image
func GetAlbumImage(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.ServerInvalidQuery)
	}

	var album *models.Album
	if result := config.Database.First(&album, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	if album.Private {
		user := c.Locals("CurrentUser").(*models.User)

		if user == nil {
			return c.Status(404).JSON(types.Error{
				Error: types.RecordNotFound,
			})
		}

		if user.ID == album.UserID {
			return c.Status(200).SendFile(album.Image.String, false)
		}
	}

	return c.Status(200).SendFile(album.Image.String, false)
}

// GET /api/v2/public/time
func GetTime(c *fiber.Ctx) error {
	return c.Status(200).JSON(map[string]int64{
		"time": time.Now().Unix(),
	})
}
