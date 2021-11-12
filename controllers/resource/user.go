package resource

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/models"
	"github.com/anti-lgbt/medusa/services"
	"github.com/anti-lgbt/medusa/types"
	"github.com/gofiber/fiber/v2"
)

const (
	UserAvatarNotValid          = "resource.user.avatar_not_valid"
	PasswordDoesntNotMatch      = "resource.password.not_match"
	PasswordPrevPassNotConrrect = "resource.password.prev_pass.not_conrrect"
	PasswordNoChangeProvided    = "resource.password.no_change_provided"
)

// PUT /api/v2/resource/users
func UpdateUser(c *fiber.Ctx) error {
	type Payload struct {
		FirstName string         `json:"first_name" form:"first_name"`
		LastName  string         `json:"last_name" form:"last_name"`
		Bio       sql.NullString `json:"bio" form:"bio"`
	}

	user := c.Locals("CurrentUser").(*models.User)
	var params *Payload
	if err := c.BodyParser(&params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	user.FirstName = params.FirstName
	user.LastName = params.LastName
	user.Bio = params.Bio

	file_header, err := c.FormFile("avatar")
	if err != nil {
		if !services.VerifyFileType(file_header, types.FileTypeImage) {
			return c.Status(422).JSON(types.Error{
				Error: UserAvatarNotValid,
			})
		}

		file_name := uuid.New().String() + filepath.Ext(file_header.Filename)
		file_path := fmt.Sprintf("./uploads/%s", file_name)

		c.SaveFile(file_header, file_path)

		user.Avatar = sql.NullString{
			Valid:  true,
			String: file_path,
		}
	}

	config.Database.Save(&user)

	return c.Status(201).JSON(201)
}

// PUT /api/v2/resouce/users/password
func UpdateUserPassword(c *fiber.Ctx) error {
	type Payload struct {
		OldPassword     string `json:"old_password" form:"old_password"`
		NewPassword     string `json:"new_password" form:"new_password"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
	}

	user := c.Locals("CurrentUser").(*models.User)

	var params *Payload
	if err := c.BodyParser(&params); err != nil {
		return c.Status(500).JSON(types.Error{
			Error: types.ServerInvalidBody,
		})
	}

	if params.NewPassword != params.ConfirmPassword {
		return c.Status(422).JSON(types.Error{
			Error: PasswordDoesntNotMatch,
		})
	}

	if services.DecryptPassword(user.Password) != params.OldPassword {
		return c.Status(422).JSON(types.Error{
			Error: PasswordPrevPassNotConrrect,
		})
	}

	if params.OldPassword == params.NewPassword {
		return c.Status(422).JSON(types.Error{
			Error: PasswordNoChangeProvided,
		})
	}

	user.UpdatePassword(params.NewPassword)

	services.SendEmail("password_changed", user.Email, user.Language(), nil)

	return c.Status(201).JSON(201)
}
