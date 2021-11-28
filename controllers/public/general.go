package public

import (
	"database/sql"
	"time"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/entities"
	"github.com/anti-lgbt/medusa/controllers/helpers"
	"github.com/anti-lgbt/medusa/controllers/queries"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/routes/middlewares"
	"github.com/anti-lgbt/medusa/types"
	"github.com/creasty/defaults"
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

	var avatar_url string
	if user.Avatar.Valid {
		avatar_url = user.Avatar.String
	} else {
		avatar_url = "./config/avatar.png"
	}

	return c.Status(200).SendFile(avatar_url, false)
}

// GET /api/v2/public/musics
func GetMusics(c *fiber.Ctx) error {
	type Payload struct {
		queries.Pagination
		queries.Period
		queries.Order
	}

	var params = new(Payload)
	if c.QueryParser(params) != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	defaults.Set(params)

	if err := helpers.Vaildate(params, "resource.music"); err != nil {
		return c.Status(422).JSON(types.Error{
			Error: err.Error(),
		})
	}

	tx := config.Database
	tx = queries.QueryPagination(tx, params.Limit, params.Page)
	tx = queries.QueryOrder(tx, params.OrderBy, params.Ordering)
	tx = queries.QueryPeriod(tx, params.TimeFrom, params.TimeTo)

	var musics []*models.Music
	tx.Find(&musics)

	user, _ := middlewares.IsAuth(c)
	music_entities := make([]*entities.Music, 0)

	for _, music := range musics {
		var entity *entities.Music

		if user != nil {
			entity = music.ToEntity(sql.NullInt64{
				Int64: user.ID,
				Valid: true,
			})
		} else {
			entity = music.ToEntity(sql.NullInt64{})
		}

		music_entities = append(music_entities, entity)
	}

	return c.Status(200).JSON(music_entities)
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

	user, _ := middlewares.IsAuth(c)
	var entity *entities.Music

	if user != nil {
		entity = music.ToEntity(sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		})
	} else {
		entity = music.ToEntity(sql.NullInt64{})
	}

	if user, err := middlewares.IsAuth(c); err == nil {
		entity.Liked = music.Liked(user.ID)
	}

	return c.Status(200).JSON(entity)
}

// GET /api/v2/public/musics/:id/comments
func GetMusicComments(c *fiber.Ctx) error {
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

	var comments []*models.Comment

	config.Database.Order("id ASC").Find(&comments, "music_id = ?", music.ID)

	user, _ := middlewares.IsAuth(c)

	comment_entities := make([]*entities.Comment, 0)
	for _, comment := range comments {
		var entity *entities.Comment

		if user != nil {
			entity = comment.ToEntity(sql.NullInt64{
				Int64: user.ID,
				Valid: true,
			})
		} else {
			entity = comment.ToEntity(sql.NullInt64{})
		}

		comment_entities = append(comment_entities, entity)
	}

	return c.Status(200).JSON(comment_entities)
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

	music.ViewCount++
	config.Database.Save(&music)

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

	user, _ := middlewares.IsAuth(c)

	var entity *entities.Album
	if user != nil {
		entity = album.ToEntity(sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		})
	} else {
		entity = album.ToEntity(sql.NullInt64{})
	}

	if user, err := middlewares.IsAuth(c); err == nil {
		entity.Liked = album.Liked(user.ID)
	}

	return c.Status(200).JSON(entity)
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

	return c.Status(200).SendFile(album.Image.String, false)
}

// GET /api/v2/public/albums/:id/comments
func GetAlbumComments(c *fiber.Ctx) error {
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

	var comments []*models.Comment

	config.Database.Order("id ASC").Find(&comments, "album_id = ?", album.ID)

	user, _ := middlewares.IsAuth(c)

	comment_entities := make([]*entities.Comment, 0)
	for _, comment := range comments {
		var entity *entities.Comment

		if user != nil {
			entity = comment.ToEntity(sql.NullInt64{
				Int64: user.ID,
				Valid: true,
			})
		} else {
			entity = comment.ToEntity(sql.NullInt64{})
		}

		comment_entities = append(comment_entities, entity)
	}

	return c.Status(200).JSON(comment_entities)
}

// GET /api/v2/public/time
func GetTime(c *fiber.Ctx) error {
	return c.Status(200).JSON(map[string]int64{
		"time": time.Now().Unix(),
	})
}
