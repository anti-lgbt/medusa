package resource

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/creasty/defaults"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/volatiletech/null"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/entities"
	"github.com/anti-lgbt/medusa/controllers/helpers"
	"github.com/anti-lgbt/medusa/controllers/queries"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/services"
	"github.com/anti-lgbt/medusa/types"
)

type AlbumPayload struct {
	Name        string      `json:"name" form:"name" validate:"required"`
	Description null.String `json:"description" form:"description"`
	Musics      []int64     `json:"music" form:"music" validate:"required" default:"[]"`
}

// GET /api/v2/resource/albums
func GetAlbums(c *fiber.Ctx) error {
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

	if err := helpers.Vaildate(params, "resource.album"); err != nil {
		return c.Status(422).JSON(types.Error{
			Error: err.Error(),
		})
	}

	user := c.Locals("CurrentUser").(*models.User)

	tx := config.Database
	tx = queries.QueryPagination(tx, params.Limit, params.Page)
	tx = queries.QueryOrder(tx, params.OrderBy, params.Ordering)
	tx = queries.QueryPeriod(tx, params.TimeFrom, params.TimeTo)

	var albums []*models.Album
	tx.Find(&albums, "user_id = ?", user.ID)

	album_entities := make([]*entities.Album, 0)

	for _, album := range albums {
		album_entities = append(album_entities, album.ToEntity(sql.NullInt64{
			Int64: user.ID,
			Valid: true,
		}))
	}

	return c.Status(200).JSON(album_entities)
}

// GET /api/v2/resource/albums/:id
func GetAlbum(c *fiber.Ctx) error {
	var album *models.Music
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	user := c.Locals("CurrentUser").(*models.User)

	if result := config.Database.Find(&album, "id = ? AND user_id = ?", id, user.ID); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	return c.Status(200).JSON(album.ToEntity(sql.NullInt64{
		Int64: user.ID,
		Valid: true,
	}))
}

// POST /api/v2/resource/albums
func CreateAlbum(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	params := new(AlbumPayload)
	if err := c.BodyParser(params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	defaults.Set(params)

	if err := helpers.Vaildate(params, "resource.album"); err != nil {
		return c.Status(422).JSON(types.Error{
			Error: err.Error(),
		})
	}

	music_albums := make([]*models.MusicAlbum, 0)

	for _, music_id := range params.Musics {
		music_albums = append(music_albums, &models.MusicAlbum{
			MusicID: music_id,
		})
	}

	album := &models.Album{
		UserID: user.ID,
		Name:   params.Name,
		Description: sql.NullString{
			String: params.Description.String,
			Valid:  params.Description.Valid,
		},
		ViewCount:   0,
		MusicAlbums: music_albums,
	}

	file_header, err := c.FormFile("image")
	if err == nil {
		services.VerifyFileType(file_header, types.FileTypeAudio)

		file_name := uuid.New().String() + filepath.Ext(file_header.Filename)
		file_path := fmt.Sprintf("./uploads/%s", file_name)

		c.SaveFile(file_header, file_path)

		album.Image = sql.NullString{
			String: file_name,
			Valid:  true,
		}
	}

	config.Database.Create(&album)

	return c.Status(201).JSON(album.ToEntity(sql.NullInt64{
		Int64: user.ID,
		Valid: true,
	}))
}

// PUT /api/v2/resource/albums/:id
func UpdateAlbum(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	params := new(AlbumPayload)
	if err := c.BodyParser(params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	defaults.Set(params)

	if err := helpers.Vaildate(params, "resource.album"); err != nil {
		return c.Status(422).JSON(types.Error{
			Error: err.Error(),
		})
	}

	var album *models.Album
	if result := config.Database.First(&album, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	album.Name = params.Name
	album.Description = sql.NullString{
		String: params.Description.String,
		Valid:  params.Description.Valid,
	}

	file_header, err := c.FormFile("image")
	if err == nil {
		services.VerifyFileType(file_header, types.FileTypeAudio)

		file_name := uuid.New().String() + filepath.Ext(file_header.Filename)
		file_path := fmt.Sprintf("./uploads/%s", file_name)

		c.SaveFile(file_header, file_path)

		album.Image = sql.NullString{
			String: file_name,
			Valid:  true,
		}
	}

	config.Database.Delete(models.MusicAlbum{}, "album_id = ? AND music_id NOT IN (?)", id, params.Musics)

	for _, music_id := range params.Musics {
		var music_album *models.MusicAlbum

		if result := config.Database.First(&music_album, "album_id = ? AND music_id = ?", id, music_id); result.Error != nil {
			music_album = &models.MusicAlbum{
				MusicID: music_id,
				AlbumID: int64(id),
			}

			config.Database.Create(&music_album)
		}
	}

	config.Database.Save(&album)

	return c.Status(201).JSON(album.ToEntity(sql.NullInt64{
		Int64: user.ID,
		Valid: true,
	}))
}

// POST /api/v2/resource/albums/:id/fork
func ForkAlbum(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var source_album *models.Album
	if result := config.Database.First(&source_album, "id = ? AND private = ?", id, false); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	music_albums := make([]*models.MusicAlbum, 0)

	for _, music_album := range source_album.MusicAlbums {
		music_albums = append(music_albums, &models.MusicAlbum{
			MusicID: music_album.MusicID,
		})
	}

	album := &models.Album{
		UserID:      user.ID,
		Name:        source_album.Name,
		Description: source_album.Description,
		ViewCount:   0,
		MusicAlbums: music_albums,
		Image:       source_album.Image,
	}

	config.Database.Create(&album)

	return c.Status(201).JSON(album.ToEntity(sql.NullInt64{
		Int64: user.ID,
		Valid: true,
	}))
}

// DELETE /api/v2/resource/albums/:id
func DeleteAlbum(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var album *models.Album
	if result := config.Database.First(&album, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	album.Delete()

	return c.Status(200).JSON(200)
}

// POST /api/v2/resource/albums/:id/like
func LikeAlbum(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var album *models.Album
	if result := config.Database.First(&album, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	album.Like(user.ID)

	return c.Status(200).JSON(200)
}

// POST /api/v2/resource/albums/:id/unlike
func UnLikeAlbum(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var album *models.Album
	if result := config.Database.First(&album, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	album.UnLike(user.ID)

	return c.Status(200).JSON(200)
}

// POST /api/v2/resource/albums/:id/comment
func CommentAlbum(c *fiber.Ctx) error {
	type Payload struct {
		Content string `json:"content" form:"content"`
	}

	user := c.Locals("CurrentUser").(*models.User)

	params := new(Payload)
	if err := c.BodyParser(params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var album *models.Album
	if result := config.Database.First(&album, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	comment := album.Comment(user.ID, params.Content)

	return c.Status(200).JSON(comment.ToEntity(sql.NullInt64{
		Int64: user.ID,
		Valid: true,
	}))
}
