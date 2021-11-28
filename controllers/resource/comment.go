package resource

import (
	"database/sql"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/controllers/helpers"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/types"
	"github.com/gofiber/fiber/v2"
)

// POST /api/v2/resource/comments/:id/like
func LikeComment(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var comment *models.Comment
	if result := config.Database.First(&comment, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	comment.Like(user.ID)

	return c.Status(200).JSON(200)
}

// POST /api/v2/resource/comments/:id/unlike
func UnLikeComment(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var comment *models.Comment
	if result := config.Database.First(&comment, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	comment.UnLike(user.ID)

	return c.Status(200).JSON(200)
}

// DELETE /api/v2/resource/comments/:id
func DeleteComment(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var comment *models.Comment
	if result := config.Database.First(&comment, "id = ? AND user_id = ?", id, user.ID); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	comment.Delete()

	return c.Status(200).JSON(200)
}

// POST /api/v2/resource/comments/reply
func CreateReply(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	type Payload struct {
		CommentID int64  `json:"comment_id"`
		Content   string `json:"content"`
	}

	params := new(Payload)
	if err := c.BodyParser(params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	if err := helpers.Vaildate(params, "resource.comment"); err != nil {
		return c.Status(422).JSON(types.Error{
			Error: err.Error(),
		})
	}

	reply := &models.Reply{
		UserID:    user.ID,
		CommentID: params.CommentID,
		Content:   params.Content,
	}

	config.Database.Create(&reply)

	return c.Status(201).JSON(reply.ToEntity(sql.NullInt64{
		Int64: user.ID,
		Valid: true,
	}))
}

// DELETE /api/v2/resource/comments/reply/:id
func DeleteReply(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var reply *models.Reply
	if result := config.Database.First(&reply, "id = ? AND user_id = ?", id, user.ID); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	reply.Delete()

	return c.Status(200).JSON(200)
}

// POST /api/v2/resource/comments/reply/:id/like
func LikeReply(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var reply *models.Reply
	if result := config.Database.First(&reply, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	reply.Like(user.ID)

	return c.Status(200).JSON(200)
}

// POST /api/v2/resource/comments/reply/:id/unlike
func UnLikeReply(c *fiber.Ctx) error {
	user := c.Locals("CurrentUser").(*models.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidQuery,
		})
	}

	var reply *models.Reply
	if result := config.Database.First(&reply, id); result.Error != nil {
		return c.Status(404).JSON(types.Error{
			Error: types.RecordNotFound,
		})
	}

	reply.UnLike(user.ID)

	return c.Status(200).JSON(200)
}
