package resource

import (
	"fmt"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/queries"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/models/datatypes"
	"github.com/anti-lgbt/medusa/services"
	"github.com/anti-lgbt/medusa/types"
)

type AlbumPayload struct {
	ID          int64                `json:"id" form:"id"`
	Name        string               `json:"name" form:"name"`
	Description datatypes.NullString `json:"description" form:"description"`
	Private     bool                 `json:"private" form:"private"`
	Musics      []int64              `json:"music" form:"music"`
}

// GET /api/v2/resource/albums
func GetAlbums(c *fiber.Ctx) error {
	type Payload struct {
		queries.Pagination
	}

	var albums []*models.Album
	var params = new(Payload)
	if c.QueryParser(params) != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	user := c.Locals("CurrentUser").(*models.User)

	config.Database.Find(&albums, "user_id = ?", user.ID).Offset(params.Page*params.Limit - params.Limit).Limit(params.Limit)

	return c.Status(200).JSON(albums)
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

	return c.Status(200).JSON(album)
}

// POST /api/v2/resource/albums
func CreateAlbum(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	var params *AlbumPayload
	if err := c.BodyParser(&params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	music_albums := make([]*models.MusicAlbum, 0)

	for _, music_id := range params.Musics {
		if !params.Private {
			var music *models.Music

			if result := config.Database.First(&music, "id = ? AND private = ?", music_id, true); result.Error != nil {
				return c.Status(404).JSON(types.Error{
					Error: types.RecordNotFound,
				})
			}
		}

		music_albums = append(music_albums, &models.MusicAlbum{
			MusicID: music_id,
		})
	}

	album := &models.Album{
		UserID:      user.ID,
		Name:        params.Name,
		Description: params.Description,
		Private:     params.Private,
		ViewCount:   0,
		MusicAlbums: music_albums,
	}

	file_header, err := c.FormFile("image")
	if err == nil {
		services.VerifyFileType(file_header, types.FileTypeAudio)

		file_name := uuid.New().String() + filepath.Ext(file_header.Filename)
		file_path := fmt.Sprintf("./uploads/%s", file_name)

		c.SaveFile(file_header, file_path)

		album.Image = datatypes.NullString{
			String: file_name,
			Valid:  true,
		}
	}

	config.Database.Create(&album)

	return c.Status(201).JSON(album)
}

// PUT /api/v2/resource/albums
func UpdateAlbum(c *fiber.Ctx) error {
	var params *AlbumPayload
	if err := c.BodyParser(&params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	var album *models.Album
	if result := config.Database.First(&album, params.ID); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	album.Name = params.Name
	album.Description = params.Description
	album.Private = params.Private

	file_header, err := c.FormFile("image")
	if err == nil {
		services.VerifyFileType(file_header, types.FileTypeAudio)

		file_name := uuid.New().String() + filepath.Ext(file_header.Filename)
		file_path := fmt.Sprintf("./uploads/%s", file_name)

		c.SaveFile(file_header, file_path)

		album.Image = datatypes.NullString{
			String: file_name,
			Valid:  true,
		}
	}

	config.Database.Delete(models.MusicAlbum{}, "album_id = ? AND music_id NOT IN (?)", params.ID, params.Musics)

	for _, music_id := range params.Musics {
		if !params.Private {
			var music *models.Music

			if result := config.Database.First(&music, "id = ? AND private = ?", music_id, true); result.Error != nil {
				return c.Status(404).JSON(types.Error{
					Error: types.RecordNotFound,
				})
			}
		}

		var music_album *models.MusicAlbum

		if result := config.Database.First(&music_album, "album_id = ? AND music_id = ?", params.ID, music_id); result.Error != nil {
			music_album = &models.MusicAlbum{
				MusicID: music_id,
				AlbumID: params.ID,
			}

			config.Database.Create(&music_album)
		}
	}

	config.Database.Save(&album)

	return c.Status(201).JSON(album)
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
		Private:     false,
		ViewCount:   0,
		MusicAlbums: music_albums,
		Image:       source_album.Image,
	}

	config.Database.Create(&album)

	return c.Status(201).JSON(album)
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

	var params *Payload
	if err := c.BodyParser(&params); err != nil {
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

	return c.Status(200).JSON(comment)
}
