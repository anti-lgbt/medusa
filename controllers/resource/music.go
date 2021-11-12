package resource

import (
	"fmt"
	"path/filepath"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/entities"
	"github.com/anti-lgbt/medusa/controllers/queries"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/models/datatypes"
	"github.com/anti-lgbt/medusa/services"
	"github.com/anti-lgbt/medusa/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	MissingMusic  = "resource.music.missing_music"
	MusicNotValid = "resource.music.music_not_valid"
	ImageNotValid = "resource.music.image_not_valid"
)

// GET /api/v2/resource/musics
func GetMusics(c *fiber.Ctx) error {
	type Payload struct {
		queries.Pagination
	}

	var musics []*models.Music
	var params = new(Payload)
	if c.QueryParser(params) != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	user := c.Locals("CurrentUser").(*models.User)

	config.Database.Find(&musics, "user_id = ?", user.ID).Offset(params.Page*params.Limit - params.Limit).Limit(params.Limit)

	music_entities := make([]*entities.Music, 0)

	for _, music := range musics {
		music_entities = append(music_entities, music.ToEntity())
	}

	return c.Status(200).JSON(music_entities)
}

// GET /api/v2/resource/musics/:id
func GetMusic(c *fiber.Ctx) error {
	var music *models.Music
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	user := c.Locals("CurrentUser").(*models.User)

	if result := config.Database.Find(&music, "id = ? AND user_id = ?", id, user.ID); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	return c.Status(200).JSON(music.ToEntity())
}

type MusicPayload struct {
	ID          int64                `json:"id" form:"id"`
	Name        string               `json:"name" form:"name"`
	Description datatypes.NullString `json:"description" form:"description"`
	State       types.MusicState     `json:"state" form:"state"`
}

// POST /api/v2/resource/musics
func CreateMusic(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	params := new(MusicPayload)
	if c.BodyParser(params) != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	music := &models.Music{
		UserID:      user.ID,
		Name:        params.Name,
		Description: params.Description,
		State:       params.State,
	}

	music_file_header, err := c.FormFile("music")
	if err != nil {
		return c.Status(422).JSON(types.Error{
			Error: MissingMusic,
		})
	}

	if !services.VerifyFileType(music_file_header, types.FileTypeAudio) {
		return c.Status(422).JSON(types.Error{
			Error: MusicNotValid,
		})
	}

	file_name := uuid.New().String() + filepath.Ext(music_file_header.Filename)
	file_path := fmt.Sprintf("./uploads/musics/%s", file_name)

	c.SaveFile(music_file_header, file_path)

	music.Path = file_path

	image_file_header, err := c.FormFile("image")
	if err == nil {
		if !services.VerifyFileType(image_file_header, types.FileTypeImage) {
			return c.Status(422).JSON(types.Error{
				Error: ImageNotValid,
			})
		}

		file_name := uuid.New().String() + filepath.Ext(image_file_header.Filename)
		file_path := fmt.Sprintf("./uploads/images/%s", file_name)

		c.SaveFile(image_file_header, file_path)

		music.Image = datatypes.NullString{
			String: file_path,
			Valid:  true,
		}
	}

	config.Database.Create(&music)

	return c.Status(200).JSON(music.ToEntity())
}

// PUT /api/v2/resource/musics
func UpdateMusic(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	params := new(MusicPayload)
	if c.BodyParser(params) != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	var music *models.Music
	if result := config.Database.First(&music, "id = ? AND user_id = ?", params.ID, user.ID); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	music.Name = params.Name
	music.Description = params.Description
	music.State = params.State

	music_file_header, err := c.FormFile("music")
	if err != nil {
		return c.Status(422).JSON(types.Error{
			Error: MissingMusic,
		})
	}

	if !services.VerifyFileType(music_file_header, types.FileTypeAudio) {
		return c.Status(422).JSON(types.Error{
			Error: MusicNotValid,
		})
	}

	file_name := uuid.New().String() + filepath.Ext(music_file_header.Filename)
	file_path := fmt.Sprintf("./uploads/musics/%s", file_name)

	c.SaveFile(music_file_header, file_path)

	music.Path = file_path

	image_file_header, err := c.FormFile("image")
	if err == nil {
		if !services.VerifyFileType(image_file_header, types.FileTypeImage) {
			return c.Status(422).JSON(types.Error{
				Error: ImageNotValid,
			})
		}

		file_name := uuid.New().String() + filepath.Ext(image_file_header.Filename)
		file_path := fmt.Sprintf("./uploads/images/%s", file_name)

		c.SaveFile(image_file_header, file_path)

		music.Image = datatypes.NullString{
			String: file_path,
			Valid:  true,
		}
	}

	config.Database.Save(&music)

	return c.Status(200).JSON(music.ToEntity())
}

// DELETE /api/v2/resource/musics/:id
func DeleteMusic(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var music *models.Music
	if result := config.Database.First(&music, "id = ? AND user_id = ?", id, user.ID); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	music.Delete()

	return c.Status(200).JSON(200)
}

// POST /api/v2/resource/musics/:id/like
func LikeMusic(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var music *models.Music
	if result := config.Database.First(&music, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	music.Like(user.ID)

	return c.Status(200).JSON(200)
}

// POST /api/v2/resource/musics/:id/unlike
func UnLikeMusic(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var music *models.Music
	if result := config.Database.First(&music, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	music.UnLike(user.ID)

	return c.Status(200).JSON(200)
}

// POST /api/v2/resource/musics/:id/comment
func CommentMusic(c *fiber.Ctx) error {
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

	var music *models.Music
	if result := config.Database.First(&music, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	comment := music.Comment(user.ID, params.Content)

	return c.Status(200).JSON(comment.ToEntity())
}
